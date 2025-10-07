package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
    "strconv"
	"saas/internal/common/reqkit/bind"
	"saas/internal/common/reskit/response"
	"saas/internal/comment/domain"
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
        return 0,errors.New("请传递id参数")
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
// @Summary      创建 Comment
// @Description  创建新的 Comment
// @Tags         comment
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body handler.CreateRequest true "创建 Comment 请求"
// @Success      200  {object}  response.successResponse{data=handler.CommentResponse} "成功创建 Comment"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/comment [post]
func (h *HttpHandler) Create(ctx *gin.Context) {
    req := new(CreateRequest)

	if err := bind.BindingRegularAndResponse(ctx,req); err != nil {
		return
	}

    data, err := h.service.Create(&domain.Comment{
        Title:    req.Title,
        Description:  req.Description,
    })

    if err != nil {
        response.Error(ctx, err)
        return
    }

    response.Success(ctx, domainCommentToResponse(data))
}

// Read godoc
// @Summary      读取单条 Comment
// @Description  读取单条 Comment
// @Tags         comment
// @Accept       json
// @Produce      json
// @Param        id   path int true "Comment ID"
// @Success      200  {object}  response.successResponse{data=handler.CommentResponse} "成功查询 Comment"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/comment/{id} [get]
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

	response.Success(ctx, domainCommentToResponse(data))
}

// Update godoc
// @Summary      更新 Comment
// @Description  根据ID更新 Comment 信息
// @Tags         comment
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path int true "Comment ID"
// @Param        request body handler.UpdateRequest true "更新 Comment 请求"
// @Success      200  {object}  response.successResponse{data=handler.CommentResponse} "成功更新 Comment"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/comment/{id} [put]
func (h *HttpHandler) Update(ctx *gin.Context) {
    req := new(UpdateRequest)

	if err := bind.BindingRegularAndResponse(ctx,req); err != nil {
		return
	}

    data, err := h.service.Update(&domain.Comment{
        ID:           req.ID,
        Title:        req.Title,
        Description:  req.Description,
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
// @Param        id   path int true "Comment ID"
// @Success      200  {object}  response.successResponse "成功删除 Comment"
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
// @Summary      获取 Comment 列表
// @Description  根据查询参数获取Comment列表，返回当前页数据和total数量
// @Tags         comment
// @Accept       json
// @Produce      json
// @Param        keyword    query     string  false  "关键词搜索"
// @Param        page       query     int     false  "页码" default(1)
// @Param        page_size  query     int     false  "每页数量" default(10)
// @Success      200  {object}  response.successResponse{data=handler.CommentListResponse} "Comment列表"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/comment [get]
func (h *HttpHandler) List(ctx *gin.Context) {
    req := new(ListRequest)

	if err := bind.BindingRegularAndResponse(ctx,req); err != nil {
		return
	}

    data, err := h.service.List(&domain.CommentQuery{
        Keyword:  req.KeyWord,
        Page:     req.Page,
        PageSize: req.PageSize,
    })

    if err != nil {
        response.Error(ctx, err)
        return
    }

    response.Success(ctx, domainCommentListToResponse(data))
}