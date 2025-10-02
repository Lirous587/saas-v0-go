package codes

// Tenant相关错误 (xx00-xx99)
var (
	ErrTenantNotFound = ErrCode{Msg: "租户不存在", Type: ErrorTypeNotFound, Code: 1601}

	ErrTenantIDInvalid = ErrCode{Msg: "无效的租户id", Type: ErrorTypeNotFound, Code: 1621}
)
