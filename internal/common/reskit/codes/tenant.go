package codes

// 租户模块(1600-1799)
var (
	ErrTenantNotFound = ErrCode{Msg: "租户不存在", Type: ErrorTypeNotFound, Code: 1600}

	// 租户邀请 
	ErrTenantInviteTokenInvalid = ErrCode{Msg: "无效的租户邀请令牌", Type: ErrorTypeNotFound, Code: 1620}
)
