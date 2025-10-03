package codes

// Tenant相关错误 (xx00-xx99)
var (
	ErrTenantNotFound = ErrCode{Msg: "租户不存在", Type: ErrorTypeNotFound, Code: 1601}
)
