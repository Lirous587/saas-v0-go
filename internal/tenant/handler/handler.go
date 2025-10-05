package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"saas/internal/common/reqkit/bind"
	"saas/internal/common/reskit/response"
	"saas/internal/common/server"
	"saas/internal/tenant/domain"
	"strconv"
)

type HttpHandler struct {
	service domain.TenantService
}

func NewHttpHandler(service domain.TenantService) *HttpHandler {
	return &HttpHandler{
		service: service,
	}
}

func (h *HttpHandler) getID(ctx *gin.Context) (int64, error) {
	idStr := ctx.Param("id")
	if idStr == "" {
		return 0, errors.New("请传递id参数")
	}
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	if idInt == 0 {
		return 0, errors.WithStack(errors.New("无效的id"))
	}
	return int64(idInt), nil
}

// Create godoc
// @Summary      新建租户
// @Description  创建新的 Tenant
// @Tags         tenant
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body handler.CreateRequest true "创建 Tenant 请求"
// @Success      200  {object}  response.successResponse{data=handler.TenantResponse} "成功创建 Tenant"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/tenant [post]
func (h *HttpHandler) Create(ctx *gin.Context) {
	req := new(CreateRequest)

	userID, err := server.GetUserID(ctx)
	if err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	data, err := h.service.Create(&domain.Tenant{
		Name:        req.Name,
		Description: req.Description,
	},
		req.PlanID,
		userID,
	)

	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainTenantToResponse(data))
}

// Read godoc
// @Summary      查询租户基础信息
// @Description  查询租户基础信息
// @Tags         tenant
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path int true "Tenant ID"
// @Success      200  {object}  response.successResponse{data=handler.TenantResponse} "成功查询 Tenant"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/tenant/{id} [get]
func (h *HttpHandler) Read(ctx *gin.Context) {
	id, err := h.getID(ctx)
	if err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	data, err := h.service.Read(id)

	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainTenantToResponse(data))
}

// Update godoc
// @Summary      更新 Tenant
// @Description  根据ID更新 Tenant 信息
// @Tags         tenant
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path int true "Tenant ID"
// @Param        request body handler.UpdateRequest true "更新 Tenant 请求"
// @Success      200  {object}  response.successResponse{data=handler.TenantResponse} "成功更新 Tenant"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/tenant/{id} [put]
func (h *HttpHandler) Update(ctx *gin.Context) {
	id, err := h.getID(ctx)
	if err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	req := new(UpdateRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	data, err := h.service.Update(&domain.Tenant{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
	})

	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainTenantToResponse(data))
}

// Delete godoc
// @Summary      删除 Tenant
// @Description  根据ID删除 Tenant
// @Tags         tenant
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path int true "Tenant ID"
// @Success      200  {object}  response.successResponse "成功删除 Tenant"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/tenant/{id} [delete]
func (h *HttpHandler) Delete(ctx *gin.Context) {
	id, err := h.getID(ctx)
	if err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if err := h.service.Delete(id); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx)
}

// List godoc
// @Summary      获取 Tenant 列表
// @Description  根据查询参数获取Tenant列表，返回当前页数据和total数量
// @Tags         tenant
// @Accept       json
// @Produce      json
// @Param        keyword    query     string  false  "关键词搜索"
// @Param        page       query     int     false  "页码" default(1)
// @Param        page_size  query     int     false  "每页数量" default(10)
// @Success      200  {object}  response.successResponse{data=handler.TenantListResponse} "Tenant列表"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/tenant [get]
func (h *HttpHandler) List(ctx *gin.Context) {
	req := new(ListRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	data, err := h.service.List(&domain.TenantQuery{
		Keyword:  req.KeyWord,
		Page:     req.Page,
		PageSize: req.PageSize,
	})

	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainTenantListToResponse(data))
}

// Upgrade godoc
// @Summary      升级租户
// @Description  升级租户计划
// @Tags         tenant
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path int true "Tenant ID"
// @Param        request body handler.UpgradeRequestBody true "升级 Tenant 请求"
// @Success      200  {object}  response.successResponse{data=handler.TenantResponse} "成功升级 Tenant"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/tenant/upgrade/{id} [put]
func (h *HttpHandler) Upgrade(ctx *gin.Context) {

}

// GenInviteToken godoc
// @Summary      生成邀请令牌
// @Description  生成邀请令牌
// @Tags         tenant
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path int true "Tenant ID"
// @Param        request body handler.GenInviteTokenRequestBody true "生成邀请令牌 Tenant 请求"
// @Success      200  {object}  response.successResponse "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/tenant/{id}/gen_invite_token [post]
func (h *HttpHandler) GenInviteToken(ctx *gin.Context) {
	req := new(GenInviteTokenRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	token, err := h.service.GenInviteToken(&domain.GenInviteTokenPayload{
		TenantID:     req.TenantID,
		ExpireSecond: req.ExpireSecond,
	})

	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, token)
}

// Invite godoc
// @Summary      邀请指定人员,通过邮箱通知
// @Description  邀请指定人员,通过邮箱通知
// @Tags         tenant
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path int true "Tenant ID"
// @Param        request body handler.InviteRequestBody true  "邀请参数"
// @Success      200  {object}  response.successResponse "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/tenant/{id}/invite [post]
func (h *HttpHandler) Invite(ctx *gin.Context) {
	req := new(InviteRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	err := h.service.Invite(&domain.InvitePayload{
		TenantID:     req.TenantID,
		ExpireSecond: req.ExpireSecond,
		Emails:       req.Emails,
	})

	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx)
}

// Enter godoc
// @Summary      加入租户
// @Description  加入指定租户 并分配指定角色
// @Tags         tenant
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path int true "Tenant ID"
// @Param        request body handler.EntryRequestBody true "加入租户请求参数"
// @Success      200  {object}  response.successResponse "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/tenant/entry [get]
func (h *HttpHandler) Enter(ctx *gin.Context) {
	req := new(EntryRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	err := h.service.Enter(&domain.EnterPayload{
		TenantID:  req.TenantID,
		TokenKind: req.TokenKind,
		Token:     req.Token,
		Email:     req.Email,
	})

	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx)
}

// GetUserWithRole godoc
// @Summary      获取租户下的用户
// @Description  获取租户下的用户的非敏感信息以及角色信息
// @Tags         tenant
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path int true "Tenant ID"
// @Param        request body handler.ListUserWithRoleQueryRequestBody true "请求参数"
// @Success      200  {object}  response.successResponse "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/tenant/:id/users [get]
func (h *HttpHandler) GetUserWithRole(ctx *gin.Context) {
	req := new(ListUserWithRoleQueryRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	// 验证role_id是否有效
	if err := h.service.CheckRoleValidity(req.RoleID); err != nil {
		// 这里归为Error更加合适
		response.Error(ctx, err)
		return
	}

	data, err := h.service.ListUsersWithRole(&domain.UserWithRoleQuery{
		TenantID: req.TenantID,
		Nickname: req.Nickname,
		RoleID:   req.RoleID,
		Page:     req.Page,
		PageSize: req.PageSize,
	})

	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainUserWithRoleListToResponse(data))
}
