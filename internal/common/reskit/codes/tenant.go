package codes

// 租户模块(1600-1799)
var (
	ErrTenantNotFound = ErrCode{Msg: "租户不存在", Type: ErrorTypeNotFound, Code: 1600}

	ErrTenantNotCreator = ErrCode{Msg: "当前用户不为租户创建则", Type: ErrorTypeUnauthorized, Code: 1610}
)
