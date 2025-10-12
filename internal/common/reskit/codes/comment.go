package codes

// Comment相关错误 (2200-2399)
var (
	//  评论错误
	ErrCommentNotFound            = ErrCode{Msg: "当前评论不存在", Type: ErrorTypeNotFound, Code: 2200}
	ErrCommentNotFoundInNowPlate  = ErrCode{Msg: "当前板块不存在该评论", Type: ErrorTypeNotFound, Code: 2201}
	ErrCommentRootCommentNotFound = ErrCode{Msg: "当前评论评论根评论不存在", Type: ErrorTypeNotFound, Code: 2203}

	ErrCommentIllegalReply         = ErrCode{Msg: "不合法的回复评论", Type: ErrorTypeExternal, Code: 2210}
	ErrCommentNoPermissionToDelete = ErrCode{Msg: "无权限删除该评论", Type: ErrorTypeUnauthorized, Code: 2211}
	ErrCommentHierarchyInvalid     = ErrCode{Msg: "父级评论与根评论不处于同一根树", Type: ErrorTypeUnauthorized, Code: 2212}

	//  评论板块错误 2320-2339
	ErrCommentPlateNotFound = ErrCode{Msg: "评论板块不存在", Type: ErrorTypeNotFound, Code: 2320}
	ErrCommentPlateExist    = ErrCode{Msg: "该评论板块已存在", Type: ErrorTypeAlreadyExists, Code: 2321}

	//  评论配置错误 2340-2359
	ErrCommentTenantConfigNotFound = ErrCode{Msg: "该租户下的评论配置不存在", Type: ErrorTypeNotFound, Code: 2340}
	ErrCommentPlateConfigNotFound  = ErrCode{Msg: "该板块的评论配置不存在", Type: ErrorTypeNotFound, Code: 2359}
)
