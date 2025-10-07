package codes

// Comment相关错误 (xx00-xx99)
var (
    ErrCommentNotFound = ErrCode{Msg: "Comment不存在", Type: ErrorTypeNotFound, Code: 0000}
)
