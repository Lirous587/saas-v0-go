package handler

import (
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"saas/internal/common/reqkit/bind"
	"saas/internal/common/reskit/codes"
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
// @Param        path        formData  string false "自定义图片路径"
// @Param        description formData  string false "图片描述"
// @Param        category_id formData  string  false "分类id"
// @Success      200 {object} response.successResponse "请求成功"
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

	if err := h.service.Upload(
		file,
		&domain.Img{
			TenantID:    req.TenantID,
			Path:        imgPath,
			Description: req.Description,
		},
		req.CategoryID,
	); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx)
}

// Delete godoc
// @Summary      删除图片
// @Tags         img
// @Accept       json
// @Produce      json
// @Param        hard  query  bool   false "是否硬删除"
// @Success      200 {object} response.successResponse "请求成功"
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
// @Tags         img
// @Accept       json
// @Produce      json
// @Param        page        query int    false "页号"
// @Param        page_size   query int    false "页码"
// @Param        keyword     query string false "关键词"
// @Param        deleted     query bool   false "是否查询回收站图片"
// @Param        category_id query string  false "分类id"
// @Success      200 {object} response.successResponse{data=handler.ImgListResponse} "请求成功"
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
// @Summary      移除回收站图片
// @Tags         img
// @Accept       json
// @Produce      json
// @Param        id path string true "图片id"
// @Success      200 {object} response.successResponse "请求成功"
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
// @Tags         img
// @Accept       json
// @Produce      json
// @Param        id path string true "图片id"
// @Param        id path string true "图片id"
// @Success      200 {object} response.successResponse "请求成功"
// @Failure      400 {object} response.invalidParamsResponse "参数错误"
// @Failure      500 {object} response.errorResponse "服务器错误"
// @Security     BearerAuth
// @Router       /v1/img/{tenant_id}/recycle/{id} [put]
func (h *HttpHandler) RestoreFromRecycleBin(ctx *gin.Context) {
	req := new(RestoreFromRecycleBinRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	if err := h.service.RestoreFromRecycleBin(req.TenantID, req.ID); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx)
}

func (h *HttpHandler) ListenDeleteQueue() {
	h.service.ListenDeleteQueue()
}

// --- 分类管理 ---

// CreateCategory godoc
// @Summary      创建图片分类
// @Tags         img-category
// @Accept       json
// @Produce      json
// @Param        tenant_id    path   string  true  "租户id"
// @Param        request body handler.CreateCategoryRequest true "请求参数"
// @Success      200 {object} response.successResponse "请求成功"
// @Failure      400 {object} response.invalidParamsResponse "参数错误"
// @Failure      500 {object} response.errorResponse "服务器错误"
// @Security     BearerAuth
// @Router       /v1/img/{tenant_id}/category [post]
func (h *HttpHandler) CreateCategory(ctx *gin.Context) {
	req := new(CreateCategoryRequest)
	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	if err := h.service.CreateCategory(&domain.Category{
		TenantID: req.TenantID,
		Title:    req.Title,
		Prefix:   req.Prefix,
	}); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx)
}

// UpdateCategory godoc
// @Summary      更新图片分类
// @Tags         img-category
// @Accept       json
// @Produce      json
// @Param        id      		path   string  true  "分类id"
// @Param        tenant_id  path   string  true  "租户id"
// @Param        request body   handler.UpdateCategoryRequest true "请求参数"
// @Success      200 {object} response.successResponse "请求成功"
// @Failure      400 {object} response.invalidParamsResponse "参数错误"
// @Failure      500 {object} response.errorResponse "服务器错误"
// @Security     BearerAuth
// @Router       /v1/img/{tenant_id}/category/{id} [put]
func (h *HttpHandler) UpdateCategory(ctx *gin.Context) {
	req := new(UpdateCategoryRequest)
	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	if err := h.service.UpdateCategory(&domain.Category{
		ID:       req.ID,
		TenantID: req.TenantID,
		Title:    req.Title,
		Prefix:   req.Prefix,
	}); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx)
}

// DeleteCategory godoc
// @Summary      删除图片分类
// @Tags         img-category
// @Accept       json
// @Produce      json
// @Param        id path string true "分类id"
// @Param        tenant_id path string true "租户id"
// @Success      200 {object} response.successResponse "删除成功"
// @Failure      400 {object} response.invalidParamsResponse "参数错误"
// @Failure      500 {object} response.errorResponse "服务器错误"
// @Security     BearerAuth
// @Router       /v1/img/{tenant_id}/category/{id} [delete]
func (h *HttpHandler) DeleteCategory(ctx *gin.Context) {
	req := new(DeleteCategoryRequest)
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
// @Tags         img-category
// @Accept       json
// @Produce      json
// @Success      200 {object} response.successResponse{data=[]handler.CategoryResponse} "请求成功"
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

// SetR2Config godoc
// @Summary      配置图库R2配置(不包含密钥)
// @Tags         tenant
// @Accept       json
// @Produce      json
// @Param        tenant_id      path   string  true  "租户id"
// @Param        request body   handler.SetR2ConfigRequest true "请求参数"
// @Success      200 {object} response.successResponse "请求成功"
// @Failure      500 {object} response.errorResponse "服务器错误"
// @Security     BearerAuth
// @Router       /v1/img/{tenant_id}/r2_config [put]
func (h *HttpHandler) SetR2Config(ctx *gin.Context) {
	req := new(SetR2ConfigRequest)
	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	config := &domain.R2Config{
		TenantID:        req.TenantID,
		AccountID:       req.AccountID,
		AccessKeyID:     req.AccessKeyID,
		PublicBucket:    req.PublicBucket,
		PublicURLPrefix: req.PublicURLPrefix,
		DeleteBucket:    req.DeleteBucket,
	}

	err := h.service.SetR2Config(config)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx)
}

// GetR2Config godoc
// @Summary      获取配置图库R2配置
// @Tags         tenant
// @Accept       json
// @Produce      json
// @Param        tenant_id      path   string  true  "租户id"
// @Success      200 {object} response.successResponse{data=handler.R2ConfigResponse} "请求成功"
// @Failure      500 {object} response.errorResponse "服务器错误"
// @Security     BearerAuth
// @Router       /v1/img/{tenant_id}/r2_config [get]
func (h *HttpHandler) GetR2Config(ctx *gin.Context) {
	req := new(GetR2ConfigRequest)
	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	res, err := h.service.GetR2Config(req.TenantID)

	if err != nil && !errors.Is(err, codes.ErrImgR2ConfigNotFound) {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainR2ConfigToResponse(res))
}

// SetR2Secret godoc
// @Summary      配置图库R2密钥
// @Tags         tenant
// @Accept       json
// @Produce      json
// @Param        tenant_id      path   string  true  "租户id"
// @Param        request body   handler.SetR2SecretAccessKeyRequest true "请求参数"
// @Success      200 {object} response.successResponse "请求成功"
// @Failure      500 {object} response.errorResponse "服务器错误"
// @Security     BearerAuth
// @Router       /v1/img/{tenant_id}/r2_config/secret [put]
func (h *HttpHandler) SetR2Secret(ctx *gin.Context) {
	req := new(SetR2SecretAccessKeyRequest)
	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	err := h.service.SetR2SecretKey(req.TenantID, req.SecretAccessKey)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx)
}

// IsSetR2Secret godoc
// @Summary      是否已配置图库R2密钥
// @Tags         tenant
// @Accept       json
// @Produce      json
// @Param        tenant_id      path   string  true  "租户id"
// @Success      200 {object} response.successResponse{data=handler.IsSetR2SecretResponse} "请求成功"
// @Failure      500 {object} response.errorResponse "服务器错误"
// @Security     BearerAuth
// @Router       /v1/img/{tenant_id}/r2_config/secret [get]
func (h *HttpHandler) IsSetR2Secret(ctx *gin.Context) {
	req := new(IsSetR2SecretRequest)
	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	isSet, err := h.service.IsSetR2SecretKey(req.TenantID)

	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, IsSetR2SecretResponse{
		IsSet: isSet,
	})
}
