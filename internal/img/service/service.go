package service

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"saas/internal/common/reskit/codes"
	"saas/internal/common/utils"
	"saas/internal/img/domain"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type tenantR2Config struct {
	accountID       string
	accessKeyID     string
	secretAccessKey string
	publicBucket    bucket
	publicURLPrefix string
	deleteBucket    bucket
	s3Client        *s3.Client
	presignClient   *s3.PresignClient
	expireAt        time.Time
}

type tenantR2ConfigWithOnce struct {
	config *tenantR2Config
	once   sync.Once
	err    error
}

type service struct {
	repo            domain.ImgRepository
	msgQueue        domain.ImgMsgQueue
	tenantR2        sync.Map // key: int64 (tenant_id), value: *TenantR2Config
	imgMutex        sync.Map // key: int64 (imgID), value: *sync.Mutex
	ace256Encryptor *utils.AES256Encryptor
}

const tenantR2ConfigTTL = 1 * time.Hour

func NewImgService(repo domain.ImgRepository, msgQueue domain.ImgMsgQueue) domain.ImgService {
	encryptKey := os.Getenv("R2_AES256_ENCRYPTION_KEY")
	if encryptKey == "" {
		panic("R2_AES256_ENCRYPTION_KEY环境变量加载失败")
	}
	ace256Encryptor, err := utils.NewAES256Encryptor(encryptKey)
	if err != nil {
		panic(err)
	}

	svc := &service{
		repo:            repo,
		msgQueue:        msgQueue,
		ace256Encryptor: ace256Encryptor,
	}

	go svc.cleanupExpiredConfigs()

	return svc
}

func (s *service) getTenantR2Config(tenantID domain.TenantID) (*tenantR2Config, error) {
	// 用 LoadOrStore 获取或创建 wrapper
	value, _ := s.tenantR2.LoadOrStore(tenantID, &tenantR2ConfigWithOnce{})
	wrapper := value.(*tenantR2ConfigWithOnce)

	// 用 once.Do 确保只加载一次
	wrapper.once.Do(func() {
		cfg, err := s.loadTenantR2Config(tenantID)
		if err != nil {
			wrapper.err = err
			return
		}
		wrapper.config = cfg
	})

	// 如果加载失败，返回错误
	if wrapper.err != nil {
		return nil, wrapper.err
	}

	// 检查是否过期
	if time.Now().After(wrapper.config.expireAt) {
		// 过期：删除旧的，递归重新加载（会创建新 once）
		s.tenantR2.Delete(tenantID)
		return s.getTenantR2Config(tenantID)
	}

	// 未过期：延长 TTL 并返回
	wrapper.config.expireAt = time.Now().Add(tenantR2ConfigTTL)
	return wrapper.config, nil
}

func (s *service) loadTenantR2Config(tenantID domain.TenantID) (*tenantR2Config, error) {
	cfg, err := s.repo.GetTenantR2Config(tenantID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get tenant R2 config")
	}

	// 解密 secret key
	decryptedSecret, err := s.ace256Encryptor.Decrypt(cfg.GetSecretAccessKey())
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt secret key")
	}

	// 配置 S3 客户端
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, decryptedSecret, "")),
		config.WithRegion("auto"), // R2 不使用区域，但 SDK 需要
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load AWS config")
	}

	// 使用服务特定端点创建 S3 客户端（适用于 Cloudflare R2）
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.AccountID))
	})
	presignClient := s3.NewPresignClient(client)

	return &tenantR2Config{
		accountID:       cfg.AccountID,
		accessKeyID:     cfg.AccessKeyID,
		secretAccessKey: decryptedSecret,
		publicBucket:    bucket(cfg.PublicBucket),
		publicURLPrefix: cfg.PublicURLPrefix,
		deleteBucket:    bucket(cfg.DeleteBucket),
		s3Client:        client,
		presignClient:   presignClient,
		expireAt:        time.Now().Add(tenantR2ConfigTTL),
	}, nil
}

// cleanupExpiredConfigs 定期清理过期配置
func (s *service) cleanupExpiredConfigs() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.tenantR2.Range(func(key, value interface{}) bool {
			wrapper := value.(*tenantR2ConfigWithOnce)
			if wrapper.config != nil && time.Now().After(wrapper.config.expireAt) {
				s.tenantR2.Delete(key)
			}
			return true
		})
	}
}

const compressQuality = 60

// Compress 压缩图片质量，返回压缩后的图片数据
func (s *service) Compress(src io.Reader) (io.Reader, error) {
	// 注册 PNG 解码器
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)

	// 解码图片
	img, _, err := image.Decode(src)
	if err != nil {
		return nil, fmt.Errorf("解码输入图片失败: %v", err)
	}

	// 创建一个缓冲区来存储压缩后的图片数据
	output := &bytes.Buffer{}

	// 将图片编码为 JPEG 格式并以指定质量写入缓冲区
	err = jpeg.Encode(output, img, &jpeg.Options{Quality: compressQuality})
	if err != nil {
		return nil, fmt.Errorf("编码图片失败: %v", err)
	}

	return bytes.NewReader(output.Bytes()), nil
}

func (s *service) Upload(src io.Reader, img *domain.Img, categoryID int64) (*domain.Img, error) {

	// 压缩图片
	compressed, err := s.Compress(src)
	if err != nil {
		return nil, codes.ErrImgCompress.WithCause(err)
	}

	// 1. 若有 categoryID 则需要先检查分类是否存在
	var category *domain.Category
	if categoryID != 0 {
		var err error
		category, err = s.repo.FindCategoryByID(img.TenantID, categoryID)
		if err != nil {
			return nil, err
		}
	}

	// 2.查询是否有相同路径
	nowPath := ""
	if category != nil {
		nowPath = category.Prefix + "/" + img.Path
	}

	exist, err := s.repo.ExistByPath(img.TenantID, nowPath)
	if err != nil {
		return nil, err
	}

	if exist {
		return nil, codes.ErrImgPathRepeat
	}

	// 3.入库
	res, err := s.repo.Create(img, categoryID)
	if err != nil {
		return nil, err
	}

	// 加载配置
	r2Config, err := s.getTenantR2Config(img.TenantID)
	if err != nil {
		return nil, err
	}

	// 后续不要再使用 img 使用res！
	// 4.上传s3
	uploadOk := true
	if err = s.UploadFile(r2Config.s3Client, r2Config.publicBucket, compressed, res.Path); err != nil {
		uploadOk = false
		err = codes.ErrImgUploadToR3Failed.WithCause(err)
	}

	// 5.如果第4步发生错误 则删除已入库的记录
	if !uploadOk {
		if err := s.repo.Delete(img.TenantID, res.ID, true); err != nil {
			zap.L().Error("数据库入库成功但图片上传失败，尝试回滚删除数据库记录时出错",
				zap.Int64("tenant_id:", int64(img.TenantID)),
				zap.Int64("id:", res.ID),
				zap.Int64("category_id:", categoryID),
				zap.String("path", res.Path),
				zap.Error(err),
			)
		}
	}

	res.SetPublicPreURL(r2Config.publicURLPrefix)

	return res, err
}

// Delete 删除逻辑
// 硬删除 -> 直接删除 publicBucket 中的对象
// 软删除 -> 复制原有对象到不可公共访问的 deleteBucket 删除 publicBucket 中的对象 -> 类似于回收站功能
func (s *service) Delete(tenantID domain.TenantID, id int64, hard ...bool) error {
	// 为每个图片创建或获取锁
	value, _ := s.imgMutex.LoadOrStore(id, &sync.Mutex{})
	mu := value.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()

	img, err := s.repo.FindByID(tenantID, id)
	if err != nil {
		return err
	}

	if !img.CanDeleted() {
		return codes.ErrImgIllegalOperation
	}

	// 加载配置
	r2Config, err := s.getTenantR2Config(tenantID)
	if err != nil {
		return err
	}

	isHardDelete := len(hard) > 0 && hard[0]

	if isHardDelete {
		// 1.删除r3
		if err := s.DeleteFile(r2Config.s3Client, r2Config.publicBucket, img.Path); err != nil {
			return err
		}
		// 2.删除记录
		if err := s.repo.Delete(tenantID, img.ID, true); err != nil {
			return err
		}
	} else {
		// 1.从 publicBucket 移到 deleteBucket
		if err := s.CopyFileToAnotherBucket(r2Config.s3Client, r2Config.publicBucket, r2Config.deleteBucket, img.Path); err != nil {
			return err
		}
		//2.删除 publicBucket 中的对象
		if err := s.DeleteFile(r2Config.s3Client, r2Config.publicBucket, img.Path); err != nil {
			return err
		}

		// 3.软删除记录
		if err := s.repo.Delete(tenantID, img.ID, false); err != nil {
			return err
		}

		// 4.将id记录到消息队列
		if err := s.msgQueue.AddToDeleteQueue(tenantID, img.ID); err != nil {
			zap.L().Error("图片软删除：添加到定时删除队列失败",
				zap.Int64("img_id", img.ID),
				zap.String("path", img.Path),
				zap.Error(err),
			)
		}
	}

	return nil
}

const deletedPresignExpired = 1 * time.Minute

func (s *service) List(query *domain.ImgQuery) (*domain.ImgList, error) {
	// 加载配置
	r2Config, err := s.getTenantR2Config(query.TenantID)
	if err != nil {
		return nil, err
	}

	res, err := s.repo.List(query)
	if err != nil {
		return nil, err
	}
	if query.Deleted {
		g := new(errgroup.Group)

		for i := range res.List {
			g.Go(func() error {
				presignUrl, err := s.getPresignURL(r2Config.presignClient, r2Config.deleteBucket, res.List[i].Path, deletedPresignExpired)
				if err != nil {
					return nil
				}
				res.List[i].Path = presignUrl
				return nil
			})
		}

		if err := g.Wait(); err != nil {
			return nil, err
		}

	}

	for i := range res.List {
		res.List[i].SetPublicPreURL(r2Config.publicURLPrefix)
	}

	return res, nil
}

func (s *service) ListenDeleteQueue() {
	s.msgQueue.ListenDeleteQueue(func(tenantID domain.TenantID, imgID int64) {
		// 加载配置
		r2Config, err := s.getTenantR2Config(tenantID)
		if err != nil {
			zap.L().Error("加载租户R2配置失败",
				zap.Int64("tenant_id:", int64(tenantID)),
				zap.Error(err),
			)
			return
		}

		//1.先查询img
		img, err := s.repo.FindByID(tenantID, imgID, true)
		if err != nil {
			zap.L().Error("定时删除队列：查询图片失败",
				zap.Int64("tenant_id:", int64(tenantID)),
				zap.Int64("img_id", imgID),
				zap.Error(err),
			)
			return
		}

		//2.删除R3
		if err := s.repo.Delete(tenantID, imgID, true); err != nil {
			zap.L().Error("定时删除队列：删除数据库记录失败",
				zap.Int64("tenant_id:", int64(tenantID)),
				zap.Int64("img_id", imgID),
				zap.String("path", img.Path),
				zap.Error(err),
			)
			return
		}

		//	3.删除记录
		if err := s.DeleteFile(r2Config.s3Client, r2Config.deleteBucket, img.Path); err != nil {
			zap.L().Error("定时删除队列：删除存储文件失败",
				zap.Int64("img_id", imgID),
				zap.String("path", img.Path),
				zap.String("bucket", r2Config.deleteBucket.string()),
				zap.Error(err),
			)
		}

		zap.L().Info("定时删除队列：图片删除成功",
			zap.Int64("img_id", imgID),
			zap.String("path", img.Path),
		)
	})
}

// ClearRecycleBin 删除被软删除的数据
// 此时删除 deleteBucket对象 数据库记录 消息队列key
func (s *service) ClearRecycleBin(tenantID domain.TenantID, id int64) error {
	value, _ := s.imgMutex.LoadOrStore(id, &sync.Mutex{})
	mu := value.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()

	// 1.先查询图片信息
	img, err := s.repo.FindByID(tenantID, id, true)
	if err != nil {
		return err
	}

	if !img.IsDeleted() {
		return codes.ErrImgIllegalOperation
	}

	// 加载配置
	r2Config, err := s.getTenantR2Config(tenantID)
	if err != nil {
		return err
	}

	// 2.删除 deleteBucket 中的文件
	if err := s.DeleteFile(r2Config.s3Client, r2Config.deleteBucket, img.Path); err != nil {
		return err
	}

	// 3.硬删除数据库记录
	if err := s.repo.Delete(tenantID, id, true); err != nil {
		return err
	}

	// 4.从删除队列中移除（防止定时器重复删除）
	if err := s.msgQueue.RemoveFromDeleteQueue(tenantID, id); err != nil {
		zap.L().Error("清空回收站：移除定时删除任务失败",
			zap.Int64("imgID", id),
			zap.Error(err),
		)
	}

	return nil
}

func (s *service) RestoreFromRecycleBin(tenantID domain.TenantID, id int64) (*domain.Img, error) {
	value, _ := s.imgMutex.LoadOrStore(id, &sync.Mutex{})
	mu := value.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()

	// 1.查询已软删除的图片信息
	img, err := s.repo.FindByID(tenantID, id, true)
	if err != nil {
		return nil, err
	}

	if !img.IsDeleted() {
		return nil, codes.ErrImgIllegalOperation
	}

	// 加载配置
	r2Config, err := s.getTenantR2Config(tenantID)
	if err != nil {
		return nil, err
	}

	// 2.从 deleteBucket 复制回 publicBucket
	if err := s.CopyFileToAnotherBucket(r2Config.s3Client, r2Config.deleteBucket, r2Config.publicBucket, img.Path); err != nil {
		return nil, err
	}

	// 3.删除 deleteBucket 中的文件
	if err := s.DeleteFile(r2Config.s3Client, r2Config.deleteBucket, img.Path); err != nil {
		return nil, err
	}

	// 4.恢复数据库记录（取消软删除）
	res, err := s.repo.Restore(tenantID, id)
	if err != nil {
		return nil, err
	}

	// 5.从删除队列中移除
	if err := s.msgQueue.RemoveFromDeleteQueue(tenantID, id); err != nil {
		zap.L().Error("从回收站恢复：移除定时删除任务失败",
			zap.Int64("img_id", id),
			zap.Error(err),
		)
	}

	res.SetPublicPreURL(r2Config.publicURLPrefix)

	return res, nil
}

const maxCategory = 10

func (s *service) CreateCategory(category *domain.Category) (*domain.Category, error) {
	exist, err := s.repo.CategoryExistByTitle(category.TenantID, category.Title)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if exist {
		return nil, codes.ErrImgCategoryTitleRepeat
	}

	count, err := s.repo.CountCategory(category.TenantID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if count >= maxCategory {
		return nil, codes.ErrImgCategoryToMany
	}

	return s.repo.CreateCategory(category)
}

func (s *service) isCategoryExistImg(tenantID domain.TenantID, id int64) error {
	existing, err := s.repo.IsCategoryExistImg(tenantID, id)
	if err != nil {
		return err
	}

	if existing {
		return codes.ErrImgCategoryExistImg
	}
	return nil
}

func (s *service) UpdateCategory(category *domain.Category) (*domain.Category, error) {
	// 1.除开自己以外 是否有与修改之后title相同的数据
	stored, err := s.repo.FindCategoryByTitle(category.TenantID, category.Title)
	if err != nil && !errors.Is(err, codes.ErrImgCategoryNotFound) {
		return nil, err
	}
	if stored != nil && stored.ID != category.ID {
		return nil, codes.ErrImgCategoryTitleRepeat
	}

	// 2.再去查询原先数据 比对path是否一致
	// 若一致则允许更新
	// 若不一致 则需该分类下无图片关联方可进行修改
	old, err := s.repo.FindCategoryByID(category.TenantID, category.ID)
	if err != nil {
		return nil, err
	}
	if old.Prefix == category.Prefix {
		return s.repo.UpdateCategory(category)
	}

	// 若一致
	if err := s.isCategoryExistImg(category.TenantID, category.ID); err != nil {
		return nil, err
	}

	return s.repo.UpdateCategory(category)
}

func (s *service) DeleteCategory(tenantID domain.TenantID, id int64) error {
	// 检验当前分类下是否存在图片
	if err := s.isCategoryExistImg(tenantID, id); err != nil {
		return err
	}

	return s.repo.DeleteCategory(tenantID, id)
}

func (s *service) ListCategories(tenantID domain.TenantID) (categories []*domain.Category, err error) {
	return s.repo.ListCategories(tenantID)
}

func (s *service) SetR2Config(secretAccessKey string, config *domain.R2Config) error {
	encryptSecret, err := s.ace256Encryptor.Encrypt(secretAccessKey)
	if err != nil {
		return errors.WithStack(err)
	}

	config.SetSecretAccessKey(encryptSecret)

	if err := s.repo.SetTenantR2Config(config); err != nil {
		return err
	}

	// 删除缓存中的旧配置，强制下次重新加载
	s.tenantR2.Delete(config.TenantID)

	return nil

}

func (s *service) GetR2Config(tenantID domain.TenantID) (*domain.R2Config, error) {
	return s.repo.GetTenantR2Config(tenantID)
}
