package codes

// Plan相关错误 (xx00-xx99)
var (
    ErrPlanNotFound = ErrCode{Msg: "Plan不存在", Type: ErrorTypeNotFound, Code: 0000}
)
