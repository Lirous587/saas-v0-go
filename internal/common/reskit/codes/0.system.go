package codes

// 系统级错误码 (0-199)
var (
	ErrAPIForbidden     = ErrCode{Msg: "当前接口禁止访问", Type: ErrorTypeForbidden, Code: 1}
	ErrPermissionDenied = ErrCode{Msg: "当前接口无权访问", Type: ErrorTypeUnauthorized, Code: 2}
	ErrInvalidRequest   = ErrCode{Msg: "无效的请求", Type: ErrorTypeExternal, Code: 3}
)
