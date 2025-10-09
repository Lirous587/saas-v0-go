package codes

// 验证码模块错误码 (201-399)
var (
	// 生成类错误 (200-219)
	ErrCaptchaGenerateFailed = ErrCode{Msg: "验证码生成失败", Type: ErrorTypeInternal, Code: 201}
	ErrCaptchaImageEmpty     = ErrCode{Msg: "验证码图片为空", Type: ErrorTypeInternal, Code: 202}

	// 验证类错误 (220-239)
	ErrCaptchaVerifyFailed = ErrCode{Msg: "验证码验证失败", Type: ErrorTypeUnauthorized, Code: 220}
	ErrCaptchaInvalid      = ErrCode{Msg: "验证码无效", Type: ErrorTypeUnauthorized, Code: 221}
	ErrCaptchaExpired      = ErrCode{Msg: "验证码已过期", Type: ErrorTypeUnauthorized, Code: 222}
	ErrCaptchaNotFound     = ErrCode{Msg: "验证码不存在", Type: ErrorTypeNotFound, Code: 223}

	// 参数类错误 (240-259)
	ErrCaptchaFormatInvalid = ErrCode{Msg: "验证码格式错误", Type: ErrorTypeValidation, Code: 240}
	ErrCaptchaHeaderMissing = ErrCode{Msg: "缺少验证码请求头", Type: ErrorTypeValidation, Code: 241}

	// 存储类错误 (260-279)
	ErrCaptchaCacheError = ErrCode{Msg: "验证码缓存操作失败", Type: ErrorTypeInternal, Code: 260}

	// 频率限制类错误 (280-299)
	ErrCaptchaRateLimit = ErrCode{Msg: "验证码请求过于频繁", Type: ErrorTypeRateLimit, Code: 280}
)
