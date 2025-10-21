package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"saas/internal/common/reqkit/bind"
	"saas/internal/common/reskit/response"
	"saas/internal/common/server"
	"saas/internal/tenant/domain"
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
	req := new(CreateRequest)

	userID, err := server.GetUserID(ctx)
	if err != nil {
		response.InvalidParams(ctx, err)
		return
	}

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
// @Summary      获取列表
// @Tags         tenant
// @Accept       json
// @Produce      json
// @Param        keyword    query     string  false  "关键词"
// @Param        page       query     int     false  "页号"
// @Param        page_size  query     int     false  "页码"
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
// @Success      200  {object}  response.successResponse "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/tenant/upgrade/{id} [put]
func (h *HttpHandler) Upgrade(ctx *gin.Context) {
	response.Error(ctx, errors.New("暂未实现"))
	// return
}
