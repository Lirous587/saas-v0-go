package codes

// Comment相关错误 (2200-2399)
var (
	//  评论错误
	ErrCommentNotFound = ErrCode{Msg: "当前评论不存在", Type: ErrorTypeNotFound, Code: 2200}

	//  评论配置错误 2300-2399
	ErrCommentTenantConfigNotFound = ErrCode{Msg: "该租户下的评论配置不存在", Type: ErrorTypeNotFound, Code: 2300}
	ErrCommentConfigNotFound       = ErrCode{Msg: "该板块的评论配置不存在", Type: ErrorTypeNotFound, Code: 230}
)
