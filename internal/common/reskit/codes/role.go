package codes

// Role相关错误 (1200-1399)
var (
	ErrRoleNotFound = ErrCode{Msg: "角色不存在", Type: ErrorTypeNotFound, Code: 1200}

	ErrRoleInvalid = ErrCode{Msg: "无效的角色", Type: ErrorTypeExternal, Code: 1220}
)
