package handler

import (
	"saas/internal/comment/domain"
	"saas/internal/common/reqkit/bind"
	"saas/internal/common/reskit/response"
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
// @Summary      创建评论
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
	req := new(CreateRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	data, err := h.service.Create(&domain.Comment{
		// Title:       req.Title,
		// Description: req.Description,
	})

	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainCommentToResponse(data))
}

// Delete godoc
// @Summary      删除 Comment
// @Description  根据ID删除 Comment
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
// @Summary      获取评论列表
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

// SetCommentTenantConfig godoc
// @Summary      设置租户级别的评论系统配置
// @Description  设置租户级别的评论系统配置
// @Tags         comment
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        tenant_id   path int true "租户id"
// @Param        request body handler.SetCommentTenantConfigRequest true "请求参数"
// @Success      200  {object}  response.successResponse "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/comment/{tenant_id}/config [put]
func (h *HttpHandler) SetCommentTenantConfig(ctx *gin.Context) {
	req := new(SetCommentTenantConfigRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	if err := h.service.SetCommentTenantConfig(&domain.CommentTenantConfig{
		TenantID: req.TenantID,
		IfAudit:  *req.IfAudit,
	}); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx)
}

// GetCommentTenantConfig godoc
// @Summary      获取租户级别的评论系统配置
// @Description  获取租户级别的评论系统配置
// @Tags         comment
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        tenant_id   path int true "租户id"
// @Success      200  {object}  response.successResponse{data=handler.CommentTenantConfigResponse} "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/comment/{tenant_id}/config [get]
func (h *HttpHandler) GetCommentTenantConfig(ctx *gin.Context) {
	req := new(GetCommentTenantConfigRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	res, err := h.service.GetCommentTenantConfig(req.TenantID)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainCommentTenantConfigToResponse(res))
}

// SetCommentConfig godoc
// @Summary      设置板块级别的评论系统配置
// @Tags         comment
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        tenant_id   	path int true "租户id"
// @Param        belong_key   path string true "板块唯一键"
// @Param        request body handler.SetCommentConfigRequest true "请求参数"
// @Success      200  {object}  response.successResponse "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/comment/{tenant_id}/{belong_key}/config [put]
func (h *HttpHandler) SetCommentConfig(ctx *gin.Context) {
	req := new(SetCommentConfigRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	if err := h.service.SetCommentConfig(&domain.CommentConfig{
		TenantID:  req.TenantID,
		BelongKey: req.BelongKey,
		IfAudit:   *req.IfAudit,
	}); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx)
}

// GetCommentConfig godoc
// @Summary      获取板块级别的评论系统配置
// @Tags         comment
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        tenant_id   path int true "租户id"
// @Success      200  {object}  response.successResponse{data=handler.CommentConfigResponse} "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/comment/{tenant_id}/{belong_key}/config [get]
func (h *HttpHandler) GetCommentConfig(ctx *gin.Context) {
	req := new(GetCommentConfigRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	res, err := h.service.GetCommentConfig(req.TenantID, req.BelongKey)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainCommentConfigToResponse(res))
}
