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
// @Tags         role
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path int true "id"
// @Param        request body handler.CreateRequest true "请求参数"
// @Success      200  {object}  response.successResponse "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/role [post]
func (h *HttpHandler) Create(ctx *gin.Context) {
	req := new(CreateRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	if err := h.service.Create(&domain.Role{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
	}); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx)
}

// Update godoc
// @Summary      更新
// @Tags         role
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path int true "id"
// @Param        request body handler.UpdateRequest true "请求参数"
// @Success      200  {object}  response.successResponse "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/role/{id} [put]
func (h *HttpHandler) Update(ctx *gin.Context) {
	req := new(UpdateRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	if err := h.service.Update(&domain.Role{
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
// @Tags         role
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path int true "id"
// @Success      200  {object}  response.successResponse "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/role/{id} [delete]
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
// @Tags         role
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.successResponse{data=handler.RoleListResponse} "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/role [get]
func (h *HttpHandler) List(ctx *gin.Context) {
	data, err := h.service.List()

	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainRoleListToResponse(data))
}
