package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"saas/internal/common/reqkit/bind"
	"saas/internal/common/reskit/response"
	"saas/internal/role/domain"
	"strconv"
)

type HttpHandler struct {
	service domain.RoleService
}

func NewHttpHandler(service domain.RoleService) *HttpHandler {
	return &HttpHandler{
		service: service,
	}
}

func (h *HttpHandler) getID(ctx *gin.Context) (int64, error) {
	idStr := ctx.Param("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		response.InvalidParams(ctx, err)
		return 0, err
	}
	if idInt == 0 {
		err := errors.New("无效的id")
		response.InvalidParams(ctx, err)
		return 0, err
	}
	return int64(idInt), err
}

// Create godoc
// @Summary      创建 Role
// @Description  创建新的 Role
// @Tags         role
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body handler.CreateRequest true "创建 Role 请求"
// @Success      200  {object}  response.successResponse{data=handler.RoleResponse} "成功创建 Role"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/role [post]
func (h *HttpHandler) Create(ctx *gin.Context) {
	req := new(CreateRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	data, err := h.service.Create(&domain.Role{
		TenantID:    req.TenantID,
		Description: req.Description,
	})

	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainRoleToResponse(data))
}

// Update godoc
// @Summary      更新 Role
// @Description  根据ID更新 Role 信息
// @Tags         role
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path int true "Role ID"
// @Param        request body handler.UpdateRequest true "更新 Role 请求"
// @Success      200  {object}  response.successResponse{data=handler.RoleResponse} "成功更新 Role"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/role/{id} [put]
func (h *HttpHandler) Update(ctx *gin.Context) {
	id, err := h.getID(ctx)
	if err != nil {
		return
	}

	req := new(UpdateRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	data, err := h.service.Update(&domain.Role{
		ID: id,
		// Title:       req.Title,
		Description: req.Description,
	})

	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainRoleToResponse(data))
}

// Delete godoc
// @Summary      删除 Role
// @Description  根据ID删除 Role
// @Tags         role
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path int true "Role ID"
// @Success      200  {object}  response.successResponse "成功删除 Role"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/role/{id} [delete]
func (h *HttpHandler) Delete(ctx *gin.Context) {
	id, err := h.getID(ctx)
	if err != nil {
		return
	}

	if err := h.service.Delete(id); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx)
}

// List godoc
// @Summary      获取 Role 列表
// @Description  根据查询参数获取Role列表，返回当前页数据和total数量
// @Tags         role
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        keyword    query     string  false  "关键词搜索"
// @Param        page       query     int     false  "页码" default(1)
// @Param        page_size  query     int     false  "每页数量" default(10)
// @Success      200  {object}  response.successResponse{data=handler.RoleListResponse} "Role列表"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/role/:tenant_id [get]
func (h *HttpHandler) List(ctx *gin.Context) {
	req := new(ListRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	data, err := h.service.List(&domain.RoleQuery{
		Keyword:  req.KeyWord,
		Page:     req.Page,
		PageSize: req.PageSize,
	})

	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainRoleListToResponse(data))
}
