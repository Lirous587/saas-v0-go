package codes

// Role相关错误 (1800-1999)
var (
	ErrRoleNotFound = ErrCode{Msg: "Role不存在", Type: ErrorTypeNotFound, Code: 1800}

	ErrRoleInvalid = ErrCode{Msg: "无效的Role", Type: ErrorTypeExternal, Code: 1820}
)
