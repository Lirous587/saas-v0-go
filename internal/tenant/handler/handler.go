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
// @Summary      创建
// @Tags         tenant
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body handler.CreateRequest true "请求参数"
// @Success      200  {object}  response.successResponse{data=handler.TenantResponse} "请求成功"
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
// @Summary      读取单条数据
// @Tags         tenant
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path int true "id"
// @Success      200  {object}  response.successResponse{data=handler.TenantResponse} "请求成功"
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
// @Summary      更新
// @Tags         tenant
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path int true "id"
// @Param        request body handler.UpdateRequest true "请求参数"
// @Success      200  {object}  response.successResponse{data=handler.TenantResponse} "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/tenant/{id} [put]
func (h *HttpHandler) Update(ctx *gin.Context) {
	req := new(UpdateRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	data, err := h.service.Update(&domain.Tenant{
		ID:          req.ID,
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
// @Summary      删除
// @Tags         tenant
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path int true "id"
// @Success      200  {object}  response.successResponse "请求成功"
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
// @Summary      获取列表
// @Tags         tenant
// @Accept       json
// @Produce      json
// @Param        keyword    query     string  false  "关键词"
// @Param        page       query     int     false  "页码"
// @Param        page_size  query     int     false  "每页数量"
// @Success      200  {object}  response.successResponse{data=handler.TenantListResponse} "请求成功"
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
// @Tags         tenant
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path int true "id"
// @Param        request body handler.UpgradeRequest true "请求参数"
// @Success      200  {object}  response.successResponse{data=handler.TenantResponse} "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/tenant/upgrade/{id} [put]
func (h *HttpHandler) Upgrade(ctx *gin.Context) {
	panic("暂未实现")
}

// GenInviteToken godoc
// @Summary      生成邀请令牌
// @Tags         tenant
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path int true "id"
// @Param        request body handler.GenInviteTokenRequest true "请求参数"
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
// @Tags         tenant
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path int true "id"
// @Param        request body handler.InviteRequest true  "请求参数"
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
// @Tags         tenant
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path int true "id"
// @Param        request body handler.EntryRequest true "请求参数"
// @Success      200  {object}  response.successResponse "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/tenant/entry/{id} [get]
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

// GetUsers godoc
// @Summary      获取租户下的用户
// @Tags         tenant
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   			path  		int  		true 	 "租户id"
// @Param        nickname   query     string  false  "用户名"
// @Param        page       query     int     false  "页号"
// @Param        page_size  query     int     false  "页码"
// @Success      200  {object}  response.successResponse{data=handler.UserListResponse} "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/tenant/:id/users [get]
func (h *HttpHandler) GetUsers(ctx *gin.Context) {
	req := new(ListUserRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	data, err := h.service.ListUsers(&domain.UserQuery{
		TenantID: req.TenantID,
		Nickname: req.Nickname,
		Page:     req.Page,
		PageSize: req.PageSize,
	})

	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainUserListToResponse(data))
}
