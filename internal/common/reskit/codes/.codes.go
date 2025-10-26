package codes

// 错误码范围定义
const (
	// 系统级错误 (0-199)
	SystemErrorStart = 1
	SystemErrorEnd   = 199

	// 验证码模块 (200-399)
	CaptchaErrorStart = 200
	CaptchaErrorEnd   = 399

	// 用户模块 (1000-1199)
	UserErrorStart = 1000
	UserErrorEnd   = 1199


	// 租户模块(1600-1799)
	TenantErrorStart = 1600
	TenantErrorEnd   = 1799

	// 图库模块 (2000-2199)
	ImgErrorStart = 2000
	ImgErrorEnd   = 2199

	// 评论模块 (2200-2399)
	CommentErrorStart = 2200
	CommentErrorEnd   = 2399
)
