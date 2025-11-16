package handler

import (
	"errors"
	"saas/internal/comment/domain"
	"saas/internal/common/reqkit/bind"
	"saas/internal/common/reskit/codes"
	"saas/internal/common/reskit/response"
	"saas/internal/common/server"

	"github.com/gin-gonic/gin"
)

type HttpHandler struct {
	service domain.CommentService
}

func NewHttpHandler(service domain.CommentService) *HttpHandler {
	return &HttpHandler{
		service: service,
	}
}

// Create godoc
// @Summary      创建/回复评论
// @Tags         comment
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        tenant_id   	path string true "租户id"
// @Param        belong_key   path string true "板块key"
// @Param        request body handler.CreateRequest true "请求参数"
// @Success      200  {object}  response.successResponse "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/comment/{tenant_id}/{belong_key} [post]
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

	// 验证root_id和parent_id的组合是否合法
	if err := req.Validate(); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if err = h.service.Create(&domain.Comment{
		TenantID: req.TenantID,
		RootID:   req.RootID,
		ParentID: req.ParentID,
		Content:  req.Content,
		UserID:   domain.UserID(userID),
	}, req.BelongKey); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx)
}

// Delete godoc
// @Summary      删除评论
// @Tags         comment
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        tenant_id   	path string true "租户id"
// @Param        id   				path string true "评论id"
// @Success      200  {object}  response.successResponse "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/comment/{tenant_id}/{id} [delete]
func (h *HttpHandler) Delete(ctx *gin.Context) {
	userID, err := server.GetUserID(ctx)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	req := new(DeleteRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	if err := h.service.Delete(req.TenantID, domain.UserID(userID), req.ID); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx)
}

// ListRoots godoc
// @Summary      获取根级评论列表
// @Tags         comment
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        tenant_id   path 	 string 	  true 	 "租户id"
// @Param        belong_key  path 	 string 		true   "评论板块"
// @Param        last_id     query   string     false  "上页最后一条记录id"
// @Param        page_size   query   int    		false  "页码"
// @Success      200  {object}  response.successResponse{data=[]handler.CommentRootResponse} "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/comment/{tenant_id}/{belong_key}/roots [get]
func (h *HttpHandler) ListRoots(ctx *gin.Context) {
	userID, _ := server.GetUserID(ctx)

	req := new(ListRootsRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	data, err := h.service.ListRoots(req.BelongKey, domain.UserID(userID), &domain.CommentRootsQuery{
		TenantID: req.TenantID,
		PageSize: req.PageSize,
		LastID:   req.LastID,
	})

	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainCommentRootsToResponse(data))
}

// ListReplies godoc
// @Summary      获取根树下的回复评论列表
// @Tags         comment
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        tenant_id   path 	 string 	  true 	 "租户id"
// @Param        belong_key  path 	 string true   "评论板块"
// @Param        root_id   	 path 	 string 	  true   "根评论id"
// @Param        last_id     query   string    false  "上页最后一条记录id"
// @Param        page_size   query   int    false  "页码"
// @Success      200  {object}  response.successResponse{data=[]handler.CommentReplyResponse} "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/comment/{tenant_id}/{belong_key}/{root_id}/replies [get]
func (h *HttpHandler) ListReplies(ctx *gin.Context) {
	userID, _ := server.GetUserID(ctx)

	req := new(ListRepliesRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	data, err := h.service.ListReplies(req.BelongKey, domain.UserID(userID), &domain.CommentRepliesQuery{
		TenantID: req.TenantID,
		RootID:   req.RootID,
		LastID:   req.LastID,
		PageSize: req.PageSize,
	})

	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainCommentRepliesToResponse(data))
}

// ListNoAudits godoc
// @Summary      获取未审核的评论
// @Tags         comment
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        tenant_id   path 	 string 	true 	 "租户id"
// @Param        belong_key  query 	 string 	true   "评论板块"
// @Param        page_size   query   int    	false  "页码"
// @Success      200  {object}  response.successResponse{data=[]handler.CommentNoAuditResponse} "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/comment/{tenant_id}/audit [get]
func (h *HttpHandler) ListNoAudit(ctx *gin.Context) {
	req := new(ListNoAuditRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	data, err := h.service.ListNoAudits(req.BelongKey, &domain.CommentNoAuditQuery{
		TenantID: req.TenantID,
		PageSize: req.PageSize,
		Keyword:  req.Keyword,
	})

	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainCommentNoAuditsToResponse(data))
}

// Audit godoc
// @Summary      评论审计
// @Tags         comment
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        tenant_id   path string true "租户id"
// @Param        id   			 path string true "评论id"
// @Param        request body handler.AuditRequest true "请求参数"
// @Success      200  {object}  response.successResponse "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/comment/{tenant_id}/audit/{id} [put]
func (h *HttpHandler) Audit(ctx *gin.Context) {
	req := new(AuditRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	var status domain.CommentStatus
	if req.Action.isAccept() {
		status.SetApproved()
	} else {
		status.SetPending()
	}

	err := h.service.Audit(req.TenantID, req.ID, status)

	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx)
}

// ToggleLike godoc
// @Summary      点赞/取消点赞评论
// @Tags         comment
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        tenant_id   	path string true "租户id"
// @Param        id   				path string true "评论id"
// @Success      200  {object}  response.successResponse "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/comment/{tenant_id}/like/{id} [put]
func (h *HttpHandler) ToggleLike(ctx *gin.Context) {
	userID, err := server.GetUserID(ctx)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	req := new(ToggleLikeRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	if err := h.service.ToggleLike(req.TenantID, domain.UserID(userID), req.ID); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx)
}

// SetTenantConfig godoc
// @Summary      设置租户级别的评论系统配置
// @Tags         tenant
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        tenant_id   path string true "租户id"
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
// @Tags         tenant
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        tenant_id   path string true "租户id"
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
	if err != nil && !errors.Is(err, codes.ErrCommentTenantConfigNotFound) {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainTenantConfigToResponse(res))
}

// CreatePlate godoc
// @Summary      新增评论板块
// @Tags         comment
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        tenant_id   path string true "租户id"
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
		TenantID:   req.TenantID,
		BelongKey:  req.BelongKey,
		RelatedURL: req.RelatedURL,
		Summary:    req.Summary,
	}); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx)
}

// UpdatePlate godoc
// @Summary      修改评论板块
// @Tags         comment
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        tenant_id   path string true "租户id"
// @Param        id   			 path string true "板块id"
// @Param        request body handler.CreatePlateRequest true "请求参数"
// @Success      200  {object}  response.successResponse "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/comment/{tenant_id}/plate/:id [put]
func (h *HttpHandler) UpdatePlate(ctx *gin.Context) {
	req := new(UpdatePlateRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	if err := h.service.UpdatePlate(&domain.Plate{
		ID:         req.ID,
		TenantID:   req.TenantID,
		BelongKey:  req.BelongKey,
		RelatedURL: req.RelatedURL,
		Summary:    req.Summary,
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
// @Param        tenant_id   path string true "租户id"
// @Param        id   			 path string true "板块id"
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
// @Param        tenant_id   path string true "租户id"
// @Param        keyword    query     string  false  "关键词"
// @Param        page       query     int     false  "页号"
// @Param        page_size  query     int     false  "页码"
// @Success      200  {object}  response.successResponse{data=handler.PlateListResponse} "请求成功"
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

// SetPlateConfig godoc
// @Summary      设置板块级别的评论系统配置
// @Tags         comment
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        tenant_id   	path string true "租户id"
// @Param        belong_key   path string true "板块key"
// @Param        request body handler.SetPlateConfigRequest true "请求参数"
// @Success      200  {object}  response.successResponse "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/comment/{tenant_id}/plate/config/{belong_key} [put]
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
// @Param        tenant_id path string true 	"租户id"
// @Param        id  		   path string true 	"板块id"
// @Success      200  {object}  response.successResponse{data=handler.PlateConfigResponse} "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/comment/{tenant_id}/plate/config/{id} [get]
func (h *HttpHandler) GetPlateConfig(ctx *gin.Context) {
	req := new(GetPlateConfigRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	res, err := h.service.GetPlateConfig(req.TenantID, req.ID)
	if err != nil && !errors.Is(err, codes.ErrCommentPlateConfigNotFound) {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainPlateConfigToResponse(res))
}

// CheckPlateBelongKey godoc
// @Summary      检测是否有相同的板块名
// @Tags         comment
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        tenant_id   			path string true "租户id"
// @Param        belong_key   		query string true "租户id"
// @Success      200  {object}  response.successResponse "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/comment/{tenant_id}/plate/check_name [get]
func (h *HttpHandler) CheckPlateBelongKey(ctx *gin.Context) {
	req := new(PlateCheckBelongKeyRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	exist, err := h.service.CheckPlateBelongKey(req.TenantID, req.BelongKey)

	if err != nil {
		response.Error(ctx, err)
		return
	}

	if exist {
		response.Error(ctx, codes.ErrCommentPlateHasSameBelongKey)
		return
	}

	response.Success(ctx)
}
