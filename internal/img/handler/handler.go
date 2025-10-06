package handler

import (
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"saas/internal/common/reqkit/bind"
	"saas/internal/common/reskit/response"
	"saas/internal/img/domain"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type HttpHandler struct {
	service domain.ImgService
}

func NewHttpHandler(service domain.ImgService) *HttpHandler {
	return &HttpHandler{
		service: service,
	}
}

func isImage(file multipart.File) (bool, string, error) {
	buf := make([]byte, 512)
	n, _ := file.Read(buf)

	// 复位文件指针
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return false, "", fmt.Errorf("file.Seek 复位文件指针失败,reason:%v", err)
	}

	contentType := http.DetectContentType(buf[:n])
	switch contentType {
	case "image/jpeg", "image/png", "image/gif", "image/webp", "image/avif", "image/bmp", "image/svg+xml":
		return true, contentType, nil
	default:
		return false, contentType, nil
	}
}

func generateImgPath(ext string) string {
	now := time.Now().Format("2006_01_02_150405.000")
	random := rand.Intn(1000000)
	return fmt.Sprintf("%s_%d%s", now, random, ext)
}

func getExtByContentType(realType string) (ext string) {
	switch realType {
	case "image/jpeg":
		ext = ".jpg"
	case "image/png":
		ext = ".png"
	case "image/gif":
		ext = ".gif"
	case "image/webp":
		ext = ".webp"
	case "image/avif":
		ext = ".avif"
	case "image/bmp":
		ext = ".bmp"
	case "image/svg+xml":
		ext = ".svg"
	default:
		ext = ""
	}
	return ext
}

// Upload godoc
// @Summary      上传图片
// @Description  上传单张图片（支持 jpeg/png/gif/webp/avif/bmp/svg）
// @Tags         img
// @Accept       multipart/form-data
// @Produce      json
// @Param        object      formData  file  true  "图片文件"
// @Param        path        formData  string false "自定义图片路径（可选）"
// @Param        description formData  string false "图片描述（可选，最长60字）"
// @Param        category_id formData  int64  false "分类ID（可选）"
// @Success      200 {object} response.successResponse{data=handler.ImgResponse} "上传成功"
// @Failure      400 {object} response.invalidParamsResponse "参数错误"
// @Failure      500 {object} response.errorResponse "服务器错误"
// @Security     BearerAuth
// @Router       /v1/img/{tenant_id}/upload [post]
func (h *HttpHandler) Upload(ctx *gin.Context) {
	fileHeader, _ := ctx.FormFile("object")
	if fileHeader == nil {
		response.InvalidParams(ctx, errors.New("未携带对象上传"))
		return
	}

	// 将 *multipart.FileHeader 转为 io.Reader
	file, _ := fileHeader.Open()
	defer file.Close()

	ok, realType, err := isImage(file)
	if err != nil {
		response.Error(ctx, errors.Errorf("isImage执行失败: %s", err))
		return
	}
	if !ok {
		response.InvalidParams(ctx, errors.Errorf("仅支持基本图片类型上传，实际类型: %s", realType))
		return
	}

	req := new(UploadRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	var imgPath string

	fmt.Println(req.Path)

	fmt.Println(imgPath)

	// 无path则生成path
	if strings.TrimSpace(req.Path) != "" {
		ext := filepath.Ext(req.Path)
		if ext == "" {
			imgPath = req.Path + getExtByContentType(realType)
		} else {
			imgPath = req.Path
		}
	} else {
		imgPath = generateImgPath(getExtByContentType(realType))
	}

	// res, err := h.service.Upload(
	// 	file,
	// 	&domain.Img{
	// 		TenantID:    req.TenantID,
	// 		Path:        imgPath,
	// 		Description: req.Description,
	// 	},
	// 	req.CategoryID,
	// )

	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx)
	// response.Success(ctx, domainImgToResponse(res))
}

// Delete godoc
// @Summary      删除图片
// @Description  删除图片（软删除或硬删除）
// @Tags         img
// @Accept       json
// @Produce      json
// @Param        hard  query  bool   false "是否硬删除（默认false）"
// @Success      200 {object} response.successResponse "删除成功"
// @Failure      400 {object} response.invalidParamsResponse "参数错误"
// @Failure      500 {object} response.errorResponse "服务器错误"
// @Security     BearerAuth
// @Router       /v1/img/{tenant_id}/{id} [delete]
func (h *HttpHandler) Delete(ctx *gin.Context) {
	req := new(DeleteRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	if req.Hard {
		if err := h.service.Delete(req.TenantID, req.ID, true); err != nil {
			response.Error(ctx, err)
			return
		}
	} else {
		if err := h.service.Delete(req.TenantID, req.ID, false); err != nil {
			response.Error(ctx, err)
			return
		}
	}

	response.Success(ctx)
}

// List godoc
// @Summary      图片列表
// @Description  分页获取图片列表
// @Tags         img
// @Accept       json
// @Produce      json
// @Param        page        query int    false "页码（默认1）"
// @Param        page_size   query int    false "每页数量（默认5，最大50）"
// @Param        keyword     query string false "关键词（可选，最长20字）"
// @Param        deleted     query bool   false "是否查询回收站图片（默认false）"
// @Param        category_id query int64  false "分类ID（可选）"
// @Success      200 {object} response.successResponse{data=handler.ImgListResponse} "查询成功"
// @Failure      400 {object} response.invalidParamsResponse "参数错误"
// @Failure      500 {object} response.errorResponse "服务器错误"
// @Security     BearerAuth
// @Router       /v1/img/{tenant_id} [get]
func (h *HttpHandler) List(ctx *gin.Context) {
	req := new(ListRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	list, err := h.service.List(&domain.ImgQuery{
		TenantID:   req.TenantID,
		Keyword:    req.KeyWord,
		Page:       req.Page,
		PageSize:   req.PageSize,
		Deleted:    req.Deleted,
		CategoryID: req.CategoryID,
	})

	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainImgListToResponse(list))
}

// ClearRecycleBin godoc
// @Summary      清空回收站图片
// @Description  彻底删除回收站中的图片
// @Tags         img
// @Accept       json
// @Produce      json
// @Param        id path int64 true "图片ID"
// @Success      200 {object} response.successResponse "清空成功"
// @Failure      400 {object} response.invalidParamsResponse "参数错误"
// @Failure      500 {object} response.errorResponse "服务器错误"
// @Security     BearerAuth
// @Router       /v1/img/{tenant_id}/recycle/{id} [delete]
func (h *HttpHandler) ClearRecycleBin(ctx *gin.Context) {
	req := new(ClearRecycleBinRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	if err := h.service.ClearRecycleBin(req.TenantID, req.ID); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx)
}

// RestoreFromRecycleBin godoc
// @Summary      恢复回收站图片
// @Description  从回收站恢复图片
// @Tags         img
// @Accept       json
// @Produce      json
// @Param        id path int64 true "图片ID"
// @Param        id path int64 true "图片ID"
// @Success      200 {object} response.successResponse{data=handler.ImgResponse} "恢复成功"
// @Failure      400 {object} response.invalidParamsResponse "参数错误"
// @Failure      500 {object} response.errorResponse "服务器错误"
// @Security     BearerAuth
// @Router       /v1/img/{tenant_id}/recycle/{id} [put]
func (h *HttpHandler) RestoreFromRecycleBin(ctx *gin.Context) {
	req := new(RestoreFromRecycleBinRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	res, err := h.service.RestoreFromRecycleBin(req.TenantID, req.ID)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, res)
}

func (h *HttpHandler) ListenDeleteQueue() {
	h.service.ListenDeleteQueue()
}

// --- 分类管理 ---

// CreateCategory godoc
// @Summary      创建图片分类
// @Description  新建图片分类
// @Tags         img-category
// @Accept       json
// @Produce      json
// @Param        request body handler.CreateCategoryRequest true "创建分类请求"
// @Success      200 {object} response.successResponse{data=handler.CategoryResponse} "创建成功"
// @Failure      400 {object} response.invalidParamsResponse "参数错误"
// @Failure      500 {object} response.errorResponse "服务器错误"
// @Security     BearerAuth
// @Router       /v1/img/{tenant_id}/category [post]
func (h *HttpHandler) CreateCategory(ctx *gin.Context) {
	req := new(CreateCategoryRequest)
	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	res, err := h.service.CreateCategory(&domain.Category{
		TenantID: req.TenantID,
		Title:    req.Title,
		Prefix:   req.Prefix,
	})

	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainCategoryToResponse(res))
}

// UpdateCategory godoc
// @Summary      更新图片分类
// @Description  修改图片分类信息
// @Tags         img-category
// @Accept       json
// @Produce      json
// @Param        id      path   int64  true  "分类ID"
// @Param        request body   handler.UpdateCategoryRequest true "更新分类请求"
// @Success      200 {object} response.successResponse{data=handler.CategoryResponse} "更新成功"
// @Failure      400 {object} response.invalidParamsResponse "参数错误"
// @Failure      500 {object} response.errorResponse "服务器错误"
// @Security     BearerAuth
// @Router       /v1/img/{tenant_id}/category/{id} [put]
func (h *HttpHandler) UpdateCategory(ctx *gin.Context) {
	req := new(UpdateCategoryRequest)
	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	res, err := h.service.UpdateCategory(&domain.Category{
		ID:       req.ID,
		TenantID: req.TenantID,
		Title:    req.Title,
		Prefix:   req.Prefix,
	})

	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainCategoryToResponse(res))
}

// DeleteCategory godoc
// @Summary      删除图片分类
// @Description  删除指定图片分类
// @Tags         img-category
// @Accept       json
// @Produce      json
// @Param        id path int64 true "分类ID"
// @Success      200 {object} response.successResponse "删除成功"
// @Failure      400 {object} response.invalidParamsResponse "参数错误"
// @Failure      500 {object} response.errorResponse "服务器错误"
// @Security     BearerAuth
// @Router       /v1/img/{tenant_id}/category/{id} [delete]
func (h *HttpHandler) DeleteCategory(ctx *gin.Context) {
	req := new(UpdateCategoryRequest)
	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	if err := h.service.DeleteCategory(req.TenantID, req.ID); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx)
}

// ListCategories godoc
// @Summary      分类列表
// @Description  获取所有图片分类
// @Tags         img-category
// @Accept       json
// @Produce      json
// @Success      200 {object} response.successResponse{data=[]handler.CategoryResponse} "查询成功"
// @Failure      500 {object} response.errorResponse "服务器错误"
// @Security     BearerAuth
// @Router       /v1/img/{tenant_id}/categories [get]
func (h *HttpHandler) ListCategories(ctx *gin.Context) {
	req := new(ListCategoryRequest)
	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	res, err := h.service.ListCategories(req.TenantID)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainCategoriesToResponse(res))
}
