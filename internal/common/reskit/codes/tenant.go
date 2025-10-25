package codes

// 租户模块(1600-1799)
var (
	ErrTenantNotFound = ErrCode{Msg: "租户不存在", Type: ErrorTypeNotFound, Code: 1600}

	ErrTenantHasSameName = ErrCode{Msg: "存在相同的租户名", Type: ErrorTypeConflict, Code: 1601}

	ErrTenantNotCreator = ErrCode{Msg: "当前用户不为租户创建者", Type: ErrorTypeUnauthorized, Code: 1610}
)
