package codes

// Plan相关错误 (1400-1599)
var (
	ErrPlanNotFound = ErrCode{Msg: "计划不存在", Type: ErrorTypeNotFound, Code: 1400}

	ErrPlanUserLimit = ErrCode{Msg: "每个用户最多只能创建一个此类型的计划", Type: ErrorTypeConflict, Code: 1401}
)
