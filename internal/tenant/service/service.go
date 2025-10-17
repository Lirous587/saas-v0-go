package service

import (
	"fmt"
	"os"
	"saas/internal/common/email"
	"saas/internal/common/reskit/codes"
	"saas/internal/common/utils"
	planDomain "saas/internal/plan/domain"
	roleDomain "saas/internal/role/domain"
	userDomain "saas/internal/user/domain"

	"saas/internal/tenant/domain"
	"saas/internal/tenant/templates"
	"time"

	"github.com/friendsofgo/errors"
	"go.uber.org/zap"
)

type service struct {
	myDomain    string
	repo        domain.TenantRepository
	cache       domain.TenantCache
	mailer      email.Mailer
	planService planDomain.PlanService
	roleService roleDomain.RoleService
	userRepo    userDomain.UserRepository
}

func NewTenantService(repo domain.TenantRepository, cache domain.TenantCache, mailer email.Mailer, planService planDomain.PlanService, roleService roleDomain.RoleService, userRepo userDomain.UserRepository) domain.TenantService {
	myDomain := os.Getenv("DOMAIN")
	if myDomain == "" {
		panic("环境变量 DOMAIN 加载失败")
	}

	return &service{
		myDomain:    myDomain,
		repo:        repo,
		cache:       cache,
		mailer:      mailer,
		planService: planService,
		roleService: roleService,
		userRepo:    userRepo,
	}
}

func (s *service) Create(tenant *domain.Tenant, planID int64, userID int64) error {
	// 1.创建事务
	tx, err := s.repo.BeginTx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1.向tenants插入数据
	res, err := s.repo.InsertTx(tx, tenant)
	if err != nil {
		return errors.WithMessage(err, "向tenants插入数据失败")
	}

	tenantID := res.ID

	// 2.向tenant_plan插入数据
	if err = s.planService.AttchToTenantTx(tx, planID, tenantID); err != nil {
		return errors.WithMessage(err, "向tenant_plan插入数据失败")
	}

	// 3.为此租户设置tenantadmin
	superadmin := s.roleService.NewRole().GetTenantadmin()

	if err = s.repo.AssignTenantUserRoleTx(tx, tenantID, userID, superadmin.ID); err != nil {
		return errors.WithMessage(err, "为此租户设置superadmin失败")
	}

	return tx.Commit()
}

func (s *service) Read(id int64) (*domain.Tenant, error) {
	return s.repo.FindByID(id)
}

func (s *service) Update(tenant *domain.Tenant) error {
	return s.repo.Update(tenant)
}

func (s *service) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *service) List(query *domain.TenantQuery) (*domain.TenantList, error) {
	return s.repo.List(query)
}

func (s *service) GenInviteToken(payload *domain.GenInviteTokenPayload) (string, error) {
	// 生成公共签名
	token, err := s.cache.GenPublicInviteToken(payload.TenantID, time.Duration(payload.ExpireSecond))
	if err != nil {
		return "", errors.WithStack(err)
	}
	return token, nil
}

func (s *service) Invite(payload *domain.InvitePayload) error {
	// 查询到租户的基础信息
	tenant, err := s.Read(payload.TenantID)
	if err != nil {
		return errors.WithStack(err)
	}

	uniqueNumber := utils.UniqueStrings(payload.Emails)

	for _, number := range uniqueNumber {
		// 生成私有签名
		token, err := s.cache.GenSecretInviteToken(payload.TenantID, time.Duration(payload.ExpireSecond), number)
		if err != nil {
			return errors.WithStack(err)
		}

		// 构建邀请链接
		inviteLink := fmt.Sprintf("%s/v1/tenant/entry?tenant_id=%d&number=%s&token_kind=secret&token=%s", s.myDomain, payload.TenantID, number, token)

		// 计算过期时间
		expireTime := time.Now().Add(time.Duration(payload.ExpireSecond)).Format(time.DateTime)

		data := map[string]interface{}{
			"TenantName":        tenant.Name,
			"TenantDescription": tenant.Description,
			"InviteLink":        inviteLink,
			"ExpireTime":        expireTime,
		}

		if err := s.mailer.SendWithTemplate(number, "新的租户邀请", templates.TemplateInvite, data); err != nil {
			return errors.WithStack(err)
		}

	}

	return nil
}

func (s *service) Enter(paylod *domain.EnterPayload) error {
	// 使用public token认证
	if paylod.TokenKind.IsPublic() {
		if err := s.cache.ValidatePublicInviteToken(paylod.TenantID, paylod.Token); err != nil {
			return errors.WithStack(err)
		}

		if err := s.enterHelp(paylod.TenantID, paylod.Email); err != nil {
			return errors.WithStack(err)
		}

		return nil
	}

	// 使用secret token认证
	// 后续可考虑在这里分配指定角色 相较于使用public token认证 可以认为使用secret token认证的场景有更高的可信度
	err := s.cache.ValidateSecretInviteToken(paylod.TenantID, paylod.Email, paylod.Token)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := s.enterHelp(paylod.TenantID, paylod.Email); err != nil {
		return errors.WithStack(err)
	}

	if err := s.cache.DeleteSecretInviteToken(paylod.TenantID, paylod.Email); err != nil {
		zap.L().Error("删除secret租户邀请令牌失败",
			zap.Int64("tenant_id", paylod.TenantID),
			zap.String("email", paylod.Email),
			zap.Error(err),
		)
	}

	return nil
}

func (s *service) enterHelp(tenantID int64, email string) error {
	// 1.查询该邮箱是否注册
	_, err := s.userRepo.FindByEmail(email)
	if !errors.Is(codes.ErrUserNotFound, err) {
		return errors.WithStack(err)
	}

	user := new(userDomain.User)

	user.Email = email
	user.Nickname = email
	// 2.没注册就去初始化
	resUser, err := s.userRepo.Create(user)
	if err != nil {
		return errors.WithStack(err)
	}

	// 3.将用户加入到租户
	viewer := s.roleService.NewRole().GetViewer()
	if err := s.repo.AssignTenantUserRole(tenantID, resUser.ID, viewer.ID); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *service) ListUsers(query *domain.UserQuery) (*domain.UserList, error) {
	return s.repo.ListUsers(query)
}

func (s *service) CheckRoleValidity(roleID int64) error {
	return s.roleService.NewRole().CheckRoleID(roleID)
}
