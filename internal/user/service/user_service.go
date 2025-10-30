package service

import (
	"saas/internal/common/reskit/codes"
	"saas/internal/common/utils"

	"go.uber.org/zap"

	"github.com/pkg/errors"

	"saas/internal/user/domain"
)

type userService struct {
	userRepo     domain.UserRepository
	tokenService domain.TokenService
}

var (
	githubClientID     string
	githubClientSecret string
)

func NewUserService(userRepo domain.UserRepository, tokenService domain.TokenService) domain.UserService {
	githubClientID = utils.GetEnv("GITHUB_CLIENT_ID")
	githubClientSecret = utils.GetEnv("GITHUB_CLIENT_SECRET")

	return &userService{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

func (s *userService) AuthenticateWithOAuth(provider string, userInfo *domain.OAuthUserInfo) (
	*domain.User2Token, error,
) {
	// 1. 查找或创建用户
	user, _, err := s.findOrCreateUserByOAuth(provider, userInfo)
	if err != nil {
		return nil, err
	}

	// 2. 更新最后登录时间
	if err := s.userRepo.UpdateLastLogin(user.ID); err != nil {
		// 这个错误不应该阻止登录流程，记录日志即可
		zap.L().Error("更新用户最后登录时间失败", zap.String("user_id", user.ID), zap.Error(err))
	}

	// 3. 生成 Token
	payload := &domain.JwtPayload{
		UserID: user.ID,
	}

	accessToken, err := s.tokenService.GenerateAccessToken(payload)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(payload)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &domain.User2Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *userService) RefreshUserToken(refreshToken string) (*domain.User2Token, error) {
	//1 . 生成新的 access token
	accessToken, err := s.tokenService.RefreshAccessToken(refreshToken)
	if err != nil {
		return nil, err
	}

	//2. 获取payload
	payload, err := s.tokenService.ParseAccessToken(accessToken)
	if err != nil {
		return nil, err
	}

	//3. 生成新的refresh token
	newRefreshToken, err := s.tokenService.GenerateRefreshToken(payload)
	if err != nil {
		return nil, err
	}

	//4. 移除旧的refresh token
	if err := s.tokenService.RemoveRefreshToken(refreshToken); err != nil {
		return nil, err
	}

	return &domain.User2Token{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

// 私有辅助方法
func (s *userService) findOrCreateUserByOAuth(provider string, userInfo *domain.OAuthUserInfo) (
	user *domain.User, isNew bool, err error,
) {
	// 1. 先通过 OAuth ID 查找
	user, err = s.userRepo.FindByOAuthID(provider, userInfo.ID)
	if err == nil {
		// 找到用户，更新信息
		return user, false, nil
	}

	if !errors.Is(err, codes.ErrUserNotFound) {
		return nil, false, errors.WithStack(err)
	}

	// 2. 通过邮箱查找现有用户
	if userInfo.Email != "" {
		user, err = s.userRepo.FindByEmail(userInfo.Email)
		if err == nil {
			// 绑定 OAuth 到现有用户
			user, err = s.bindOAuthToUser(user, provider, userInfo)
			return user, false, err
		}

		if !errors.Is(err, codes.ErrUserNotFound) {
			return nil, false, errors.WithStack(err)
		}
	}

	// 3. 创建新用户
	user, err = s.createUserFromOAuth(provider, userInfo)
	return user, true, err
}

func (s *userService) createUserFromOAuth(provider string, userInfo *domain.OAuthUserInfo) (*domain.User, error) {
	user := &domain.User{
		Email:    userInfo.Email,
		Nickname: userInfo.Nickname,
	}

	// 设置 OAuth ID
	switch provider {
	case "github":
		user.GithubID = userInfo.ID
	}

	return s.userRepo.Create(user)
}

func (s *userService) bindOAuthToUser(user *domain.User, provider string, userInfo *domain.OAuthUserInfo) (
	*domain.User, error,
) {
	// 设置 OAuth ID
	switch provider {
	case "github":
		user.GithubID = userInfo.ID
	}

	// 更新头像等信息
	if userInfo.AvatarURL != "" {
		user.AvatarURL = userInfo.AvatarURL
	}

	return s.userRepo.Update(user)
}

func (s *userService) GetUser(id string) (*domain.User, error) {
	if err := s.userRepo.UpdateLastLogin(id); err != nil {
		zap.L().Error("更新用户登录时间失败",
			zap.String("id", id),
			zap.Error(err))
	}
	return s.userRepo.FindByID(id)
}
