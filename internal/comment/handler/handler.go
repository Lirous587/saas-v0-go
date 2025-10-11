package handler

import (
	"saas/internal/comment/domain"
	"saas/internal/common/reqkit/bind"
	"saas/internal/common/reskit/response"
	"saas/internal/common/server"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type HttpHandler struct {
	service domain.CommentService
}

func NewHttpHandler(service domain.CommentService) *HttpHandler {
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
// @Tags         comment
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body handler.CreateRequest true "请求参数"
// @Success      200  {object}  response.successResponse{data=handler.CommentResponse} "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/comment [post]
func (h *HttpHandler) Create(ctx *gin.Context) {
	userID, err := server.GetUserID(ctx)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	req := new(CreateRequest)

	if err = bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	data, err := h.service.Create(&domain.Comment{
		Plate: &domain.PlateBelong{
			BelongKey: req.BelongKey,
		},
		TenantID: req.TenantID,
		ParentID: req.ParentID,
		Content:  req.Content,
		User: &domain.UserInfo{
			ID: userID,
		},
	})

	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainCommentToResponse(data))
}

// Delete godoc
// @Summary      删除
// @Tags         comment
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path int true "评论id"
// @Success      200  {object}  response.successResponse "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/comment/{id} [delete]
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
// @Tags         comment
// @Accept       json
// @Produce      json
// @Param        keyword    query     string  false  "关键词"
// @Param        page       query     int     false  "页码"
// @Param        page_size  query     int     false  "每页数量"
// @Success      200  {object}  response.successResponse{data=handler.CommentListResponse} "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/comment [get]
func (h *HttpHandler) List(ctx *gin.Context) {
	req := new(ListRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	data, err := h.service.List(&domain.CommentQuery{
		// Keyword:  req.KeyWord,
		Page:     req.Page,
		PageSize: req.PageSize,
	})

	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainCommentListToResponse(data))
}

// CreatePlate godoc
// @Summary      新增评论板块
// @Tags         comment
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        tenant_id   path int true "租户id"
// @Param        request body handler.CreatePlateRequest true "请求参数"
// @Success      200  {object}  response.successResponse "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/comment/{tenant_id}/plate [post]
func (h *HttpHandler) CreatePlate(ctx *gin.Context) {
	req := new(CreatePlateRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	if err := h.service.CreatePlate(&domain.Plate{
		TenantID:  req.TenantID,
		BelongKey: req.BelongKey,
		Summary:   req.Summary,
	}); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx)
}

// DeletePlate godoc
// @Summary      删除评论板块
// @Tags         comment
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        tenant_id   path int true "租户id"
// @Param        id   			 path int true "板块id"
// @Success      200  {object}  response.successResponse "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/comment/{tenant_id}/plate/{id} [delete]
func (h *HttpHandler) DeletePlate(ctx *gin.Context) {
	req := new(DeletePlateRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	if err := h.service.DeletePlate(req.TenantID, req.ID); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx)
}

// ListPlate godoc
// @Summary      评论板块列表
// @Tags         comment
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        tenant_id   path int true "租户id"
// @Param        keyword    query     string  false  "关键词"
// @Param        page       query     int     false  "页码"
// @Param        page_size  query     int     false  "每页数量"
// @Success      200  {object}  response.successResponse "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/comment/{tenant_id}/plate [get]
func (h *HttpHandler) ListPlate(ctx *gin.Context) {
	req := new(PlateListRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	data, err := h.service.ListPlate(&domain.PlateQuery{
		TenantID: req.TenantID,
		Keyword:  req.Keyword,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainPlateListToResponse(data))
}

// SetTenantConfig godoc
// @Summary      设置租户级别的评论系统配置
// @Tags         comment
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        tenant_id   path int true "租户id"
// @Param        request body handler.SetTenantConfigRequest true "请求参数"
// @Success      200  {object}  response.successResponse "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/comment/{tenant_id}/config [put]
func (h *HttpHandler) SetTenantConfig(ctx *gin.Context) {
	req := new(SetTenantConfigRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	if err := h.service.SetTenantConfig(&domain.TenantConfig{
		TenantID: req.TenantID,
		IfAudit:  *req.IfAudit,
	}); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx)
}

// GetTenantConfig godoc
// @Summary      获取租户级别的评论系统配置
// @Tags         comment
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        tenant_id   path int true "租户id"
// @Success      200  {object}  response.successResponse{data=handler.TenantConfigResponse} "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/comment/{tenant_id}/config [get]
func (h *HttpHandler) GetTenantConfig(ctx *gin.Context) {
	req := new(GetTenantConfigRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	res, err := h.service.GetTenantConfig(req.TenantID)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainTenantConfigToResponse(res))
}

// SetPlateConfig godoc
// @Summary      设置板块级别的评论系统配置
// @Tags         comment
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        tenant_id   	path int true "租户id"
// @Param        belong_key   path string true "板块key"
// @Param        request body handler.SetPlateConfigRequest true "请求参数"
// @Success      200  {object}  response.successResponse "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/comment/{tenant_id}/{belong_key}/config [put]
func (h *HttpHandler) SetPlateConfig(ctx *gin.Context) {
	req := new(SetPlateConfigRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	if err := h.service.SetPlateConfig(&domain.PlateConfig{
		TenantID: req.TenantID,
		Plate: &domain.PlateBelong{
			BelongKey: req.BelongKey,
		},
		IfAudit: *req.IfAudit,
	}); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx)
}

// GetPlateConfig godoc
// @Summary      获取板块级别的评论系统配置
// @Tags         comment
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        tenant_id   path int true 		"租户id"
// @Param        belong_key  path string true "板块key"
// @Success      200  {object}  response.successResponse{data=handler.PlateConfigResponse} "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/comment/{tenant_id}/{belong_key}/config [get]
func (h *HttpHandler) GetPlateConfig(ctx *gin.Context) {
	req := new(GetPlateConfigRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	res, err := h.service.GetPlateConfig(req.TenantID, req.BelongKey)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainPlateConfigToResponse(res))
}
