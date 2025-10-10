﻿package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"saas/internal/common/reqkit/bind"
	"saas/internal/common/reskit/response"
	"saas/internal/plan/domain"
	"strconv"
)

type HttpHandler struct {
	service domain.PlanService
}

func NewHttpHandler(service domain.PlanService) *HttpHandler {
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
// @Tags         plan
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body handler.CreateRequest true "请求参数"
// @Success      200  {object}  response.successResponse{data=handler.PlanResponse} "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/plan [post]
func (h *HttpHandler) Create(ctx *gin.Context) {
	req := new(CreateRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	data, err := h.service.Create(&domain.Plan{
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
	})

	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainPlanToResponse(data))
}

// Read godoc
// @Summary      读取单条数据
// @Tags         plan
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path int true "id"
// @Success      200  {object}  response.successResponse{data=handler.PlanResponse} "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/plan/{id} [get]
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

	response.Success(ctx, domainPlanToResponse(data))
}

// Update godoc
// @Summary      更新
// @Tags         plan
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path int true "id"
// @Param        request body handler.UpdateRequest true "请求参数"
// @Success      200  {object}  response.successResponse{data=handler.PlanResponse} "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/plan/{id} [put]
func (h *HttpHandler) Update(ctx *gin.Context) {
	req := new(UpdateRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	data, err := h.service.Update(&domain.Plan{
		ID:          req.ID,
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
	})

	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainPlanToResponse(data))
}

// Delete godoc
// @Summary      删除
// @Tags         plan
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path int true "id"
// @Success      200  {object}  response.successResponse "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/plan/{id} [delete]
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
// @Tags         plan
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.successResponse{data=handler.PlanListResponse} "请求成功"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /v1/plan [get]
func (h *HttpHandler) List(ctx *gin.Context) {
	data, err := h.service.List()

	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainPlanListToResponse(data))
}
