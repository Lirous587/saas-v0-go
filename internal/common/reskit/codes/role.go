package codes

// Role相关错误 (xx00-xx99)
var (
    ErrRoleNotFound = ErrCode{Msg: "Role不存在", Type: ErrorTypeNotFound, Code: 0000}
)
