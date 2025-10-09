package codes

// 图库相关错误 (1400-1599)
var (
	// 图库基础操作 (1400-1419)
	ErrImgNotFound            = ErrCode{Msg: "图片不存在", Type: ErrorTypeNotFound, Code: 2000}
	ErrImgPathRepeat          = ErrCode{Msg: "图片路径重复,请检查图库和回收站", Type: ErrorTypeAlreadyExists, Code: 2001}
	ErrImgCategoryNotFound    = ErrCode{Msg: "图片分类不存在", Type: ErrorTypeNotFound, Code: 2002}
	ErrImgCategoryTitleRepeat = ErrCode{Msg: "图片分类名重复", Type: ErrorTypeAlreadyExists, Code: 2003}
	ErrImgCategoryToMany      = ErrCode{Msg: "图片分类过多", Type: ErrorTypeExternal, Code: 2004}
	ErrImgCategoryExistImg    = ErrCode{
		Msg: "当前图片分类下存在图片,请检查图库和回收站", Type: ErrorTypeExternal, Code: 2005,
	}
	ErrImgIllegalOperation = ErrCode{Msg: "非法的图片操作", Type: ErrorTypeExternal, Code: 2006}

	// 图片处理 (1420-1439)
	ErrImgCompress         = ErrCode{Msg: "压缩图片失败", Type: ErrorTypeInternal, Code: 2020}
	ErrImgUploadToR3Failed = ErrCode{Msg: "上传图片到R3失败", Type: ErrorTypeInternal, Code: 2021}

	// 图库配置 (1440-1459)
	ErrImgR2ConfigNotFound = ErrCode{Msg: "图库R2配置不存在", Type: ErrorTypeNotFound, Code: 2040}
)
