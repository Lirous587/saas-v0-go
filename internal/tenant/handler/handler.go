package handler

import (
	"saas/internal/common/reqkit/bind"
	"saas/internal/common/reskit/codes"
	"saas/internal/common/reskit/response"
	"saas/internal/common/server"
	"saas/internal/tenant/domain"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type HttpHandler struct {
	service domain.TenantService
}

func NewHttpHandler(service domain.TenantService) *HttpHandler {
	return &HttpHandler{
		service: service,
	}
}

// Create godoc
// @Summary      创建
// @Tags         tenant
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body handler.CreateRequest true "请求参数"
// @Success      200  {object}  response.successResponse "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/tenant [post]
func (h *HttpHandler) Create(ctx *gin.Context) {
	userID, err := server.GetUserID(ctx)
	if err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	req := new(CreateRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	if err := h.service.Create(&domain.Tenant{
		Name:        req.Name,
		Description: req.Description,
		CreatorID:   userID,
	},
		req.PlanID,
	); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx)
}

// Update godoc
// @Summary      更新
// @Tags         tenant
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path int true "租户id"
// @Param        request body handler.UpdateRequest true "请求参数"
// @Success      200  {object}  response.successResponse "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/tenant/{id} [put]
func (h *HttpHandler) Update(ctx *gin.Context) {
	req := new(UpdateRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	if err := h.service.Update(&domain.Tenant{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
	}); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx)
}

// Delete godoc
// @Summary      删除
// @Tags         tenant
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path int true "租户id"
// @Success      200  {object}  response.successResponse "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/tenant/{id} [delete]
func (h *HttpHandler) Delete(ctx *gin.Context) {
	req := new(DeleteRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	if err := h.service.Delete(req.ID); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx)
}

// List godoc
// @Summary      获取用户的租户列表
// @Tags         tenant
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        keyword    query     string  false  "关键词"
// @Param        last_id    query     int     false  "上页最后一条记录id"
// @Param        page_size  query     int     false  "页码"
// @Success      200  {object}  response.successResponse{data=[]handler.TenantResponse} "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/tenant [get]
func (h *HttpHandler) List(ctx *gin.Context) {
	userID, err := server.GetUserID(ctx)
	if err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	req := new(ListRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	data, err := h.service.List(&domain.TenantQuery{
		CreatorID: userID,
		Keyword:   req.KeyWord,
		LastID:    req.LastID,
		PageSize:  req.PageSize,
	})

	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainTenantsToResponse(data))
}

// CheckName godoc
// @Summary      检测是否有相同的租户名
// @Tags         tenant
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        name    query     string  false  "租户名称"
// @Success      200  {object}  response.successResponse "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/tenant/check_name [get]
func (h *HttpHandler) CheckName(ctx *gin.Context) {
	userID, err := server.GetUserID(ctx)
	if err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	req := new(CheckNameRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	exist, err := h.service.CheckName(userID, req.Name)

	if err != nil {
		response.Error(ctx, err)
		return
	}

	if exist {
		response.Error(ctx, codes.ErrTenantHasSameName)
		return
	}

	response.Success(ctx)
}

// Upgrade godoc
// @Summary      升级租户
// @Tags         tenant
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path int true "id"
// @Param        request body handler.UpgradeRequest true "请求参数"
// @Success      200  {object}  response.successResponse "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/tenant/upgrade/{id} [put]
func (h *HttpHandler) Upgrade(ctx *gin.Context) {
	response.Error(ctx, errors.New("暂未实现"))
	// return
}
