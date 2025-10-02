package codes

// Tenant相关错误 (xx00-xx99)
var (
    ErrTenantNotFound = ErrCode{Msg: "Tenant不存在", Type: ErrorTypeNotFound, Code: 0000}
)
