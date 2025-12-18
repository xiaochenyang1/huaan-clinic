package errorcode

// 错误码定义
// 格式说明：
// - 200xxx: 成功
// - 400xxx: 客户端错误（参数错误、业务校验失败等）
// - 401xxx: 认证错误
// - 403xxx: 权限错误
// - 404xxx: 资源不存在
// - 500xxx: 服务端错误

const (
	// 成功
	Success = 200000

	// 通用错误 400xxx
	ErrInvalidParams     = 400001 // 参数错误
	ErrBindJSON          = 400002 // JSON绑定失败
	ErrValidation        = 400003 // 参数校验失败
	ErrInvalidIDFormat   = 400004 // ID格式错误
	ErrInvalidDateFormat = 400005 // 日期格式错误
	ErrInvalidPageParams = 400006 // 分页参数错误
	ErrFileTooLarge      = 400007 // 文件过大
	ErrInvalidFileType   = 400008 // 文件类型错误
	ErrFrequentRequest   = 400009 // 请求过于频繁

	// 认证错误 401xxx
	ErrUnauthorized       = 401001 // 未登录
	ErrTokenExpired       = 401002 // Token已过期
	ErrTokenInvalid       = 401003 // Token无效
	ErrTokenMalformed     = 401004 // Token格式错误
	ErrRefreshTokenFailed = 401005 // 刷新Token失败
	ErrWeChatLoginFailed  = 401006 // 微信登录失败
	ErrAccountDisabled    = 401007 // 账号已禁用
	ErrPasswordWrong      = 401008 // 密码错误

	// 权限错误 403xxx
	ErrForbidden         = 403001 // 无权限访问
	ErrPermissionDenied  = 403002 // 权限不足
	ErrResourceForbidden = 403003 // 资源禁止访问

	// 资源不存在 404xxx
	ErrNotFound           = 404001 // 资源不存在
	ErrUserNotFound       = 404002 // 用户不存在
	ErrPatientNotFound    = 404003 // 就诊人不存在
	ErrDepartmentNotFound = 404004 // 科室不存在
	ErrDoctorNotFound     = 404005 // 医生不存在
	ErrScheduleNotFound   = 404006 // 排班不存在
	ErrAppointmentNotFound = 404007 // 预约不存在
	ErrRecordNotFound     = 404008 // 就诊记录不存在
	ErrAdminNotFound      = 404009 // 管理员不存在

	// 业务错误 - 用户相关 410xxx
	ErrPhoneExists        = 410001 // 手机号已存在
	ErrIDCardExists       = 410002 // 身份证号已存在
	ErrPatientLimitExceed = 410003 // 就诊人数量超限
	ErrPatientHasAppt     = 410004 // 就诊人有待就诊预约，无法删除
	ErrUsernameExists          = 410005 // 用户名已存在
	ErrUsernameInvalid         = 410006 // 用户名格式错误
	ErrPasswordInvalid         = 410007 // 密码格式错误
	ErrSMSCodeInvalid          = 410008 // 验证码错误
	ErrSMSCodeExpired          = 410009 // 验证码已过期
	ErrSMSCodeSendTooFrequent  = 410010 // 验证码发送频繁

	// 业务错误 - 预约相关 420xxx
	ErrScheduleUnavailable     = 420001 // 该时段不可预约
	ErrNoAvailableSlots        = 420002 // 号源已满
	ErrDuplicateAppointment    = 420003 // 重复预约（同医生同天）
	ErrTimeConflict            = 420004 // 时间冲突（同时段已有预约）
	ErrDailyLimitExceed        = 420005 // 当日预约次数已达上限
	ErrCannotCancelToday       = 420006 // 就诊当天无法取消
	ErrCancelLimitExceed       = 420007 // 本月取消次数已达上限
	ErrAppointmentCancelled    = 420008 // 预约已取消
	ErrAppointmentCompleted    = 420009 // 预约已完成
	ErrCheckinTooEarly         = 420010 // 签到时间未到
	ErrCheckinTooLate          = 420011 // 签到已过时
	ErrAlreadyCheckedIn        = 420012 // 已签到
	ErrIdempotentTokenInvalid  = 420013 // 幂等Token无效
	ErrIdempotentTokenUsed     = 420014 // 幂等Token已使用
	ErrAppointmentDateInvalid  = 420015 // 预约日期无效
	ErrUserBlocked             = 420016 // 用户已被限制预约（爽约惩罚）

	// 业务错误 - 排班相关 430xxx
	ErrScheduleConflict     = 430001 // 排班时间冲突
	ErrScheduleHasAppt      = 430002 // 排班下有预约，无法删除
	ErrInvalidSchedulePeriod = 430003 // 无效的排班时段
	ErrScheduleDatePassed   = 430004 // 排班日期已过

	// 业务错误 - 科室/医生相关 440xxx
	ErrDepartmentHasDoctor = 440001 // 科室下有医生，无法删除
	ErrDepartmentDisabled  = 440002 // 科室已停用
	ErrDoctorDisabled      = 440003 // 医生已停诊
	ErrDoctorHasSchedule   = 440004 // 医生有排班，无法删除

	// 服务端错误 500xxx
	ErrInternalServer = 500001 // 服务器内部错误
	ErrDatabase       = 500002 // 数据库错误
	ErrRedis          = 500003 // Redis错误
	ErrThirdParty     = 500004 // 第三方服务错误
	ErrFileUpload     = 500005 // 文件上传失败
)

// 错误码消息映射
var codeMessages = map[int]string{
	Success: "成功",

	// 通用错误
	ErrInvalidParams:     "参数错误",
	ErrBindJSON:          "请求数据格式错误",
	ErrValidation:        "参数校验失败",
	ErrInvalidIDFormat:   "ID格式错误",
	ErrInvalidDateFormat: "日期格式错误",
	ErrInvalidPageParams: "分页参数错误",
	ErrFileTooLarge:      "文件过大",
	ErrInvalidFileType:   "不支持的文件类型",
	ErrFrequentRequest:   "请求过于频繁，请稍后再试",

	// 认证错误
	ErrUnauthorized:       "请先登录",
	ErrTokenExpired:       "登录已过期，请重新登录",
	ErrTokenInvalid:       "登录凭证无效",
	ErrTokenMalformed:     "登录凭证格式错误",
	ErrRefreshTokenFailed: "刷新登录凭证失败",
	ErrWeChatLoginFailed:  "微信登录失败",
	ErrAccountDisabled:    "账号已被禁用",
	ErrPasswordWrong:      "密码错误",

	// 权限错误
	ErrForbidden:         "无权访问",
	ErrPermissionDenied:  "权限不足",
	ErrResourceForbidden: "资源禁止访问",

	// 资源不存在
	ErrNotFound:           "资源不存在",
	ErrUserNotFound:       "用户不存在",
	ErrPatientNotFound:    "就诊人不存在",
	ErrDepartmentNotFound: "科室不存在",
	ErrDoctorNotFound:     "医生不存在",
	ErrScheduleNotFound:   "排班信息不存在",
	ErrAppointmentNotFound: "预约不存在",
	ErrRecordNotFound:     "就诊记录不存在",
	ErrAdminNotFound:      "管理员不存在",

	// 用户相关
	ErrPhoneExists:        "手机号已被使用",
	ErrIDCardExists:       "身份证号已被使用",
	ErrPatientLimitExceed: "就诊人数量已达上限",
	ErrPatientHasAppt:     "该就诊人有待就诊预约，无法删除",
	ErrUsernameExists:          "用户名已被使用",
	ErrUsernameInvalid:         "用户名格式错误，请输入4-20位字母、数字或下划线，必须以字母开头",
	ErrPasswordInvalid:         "密码格式错误，请输入6-20位字符",
	ErrSMSCodeInvalid:          "验证码错误",
	ErrSMSCodeExpired:          "验证码已过期，请重新获取",
	ErrSMSCodeSendTooFrequent:  "验证码发送过于频繁，请稍后再试",

	// 预约相关
	ErrScheduleUnavailable:     "该时段暂不可预约",
	ErrNoAvailableSlots:        "该时段号源已满",
	ErrDuplicateAppointment:    "您当天已预约过该医生",
	ErrTimeConflict:            "该时段您已有其他预约",
	ErrDailyLimitExceed:        "当日预约次数已达上限",
	ErrCannotCancelToday:       "就诊当天无法取消预约",
	ErrCancelLimitExceed:       "本月取消次数已达上限，如需取消请联系客服",
	ErrAppointmentCancelled:    "该预约已取消",
	ErrAppointmentCompleted:    "该预约已完成",
	ErrCheckinTooEarly:         "签到时间未到，请在就诊前30分钟内签到",
	ErrCheckinTooLate:          "签到已过时，预约已自动作废",
	ErrAlreadyCheckedIn:        "您已签到，请耐心等待叫号",
	ErrIdempotentTokenInvalid:  "请求已失效，请刷新后重试",
	ErrIdempotentTokenUsed:     "请勿重复提交",
	ErrAppointmentDateInvalid:  "预约日期无效，请选择明天至7天后的日期",
	ErrUserBlocked:             "您因多次爽约已被暂时限制预约",

	// 排班相关
	ErrScheduleConflict:     "排班时间存在冲突",
	ErrScheduleHasAppt:      "该排班下有预约记录，无法删除",
	ErrInvalidSchedulePeriod: "无效的排班时段",
	ErrScheduleDatePassed:   "排班日期已过",

	// 科室/医生相关
	ErrDepartmentHasDoctor: "该科室下有医生，请先处理医生信息",
	ErrDepartmentDisabled:  "该科室已停用",
	ErrDoctorDisabled:      "该医生已停诊",
	ErrDoctorHasSchedule:   "该医生有排班记录，无法删除",

	// 服务端错误
	ErrInternalServer: "服务器开小差了，请稍后再试",
	ErrDatabase:       "数据处理失败，请稍后再试",
	ErrRedis:          "缓存服务异常",
	ErrThirdParty:     "第三方服务异常",
	ErrFileUpload:     "文件上传失败",
}

// GetMessage 获取错误码对应的消息
func GetMessage(code int) string {
	if msg, ok := codeMessages[code]; ok {
		return msg
	}
	return "未知错误"
}

// AppError 应用错误
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

// New 创建应用错误
func New(code int) *AppError {
	return &AppError{
		Code:    code,
		Message: GetMessage(code),
	}
}

// NewWithMessage 创建带自定义消息的应用错误
func NewWithMessage(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Wrap 包装错误
func Wrap(code int, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: GetMessage(code) + ": " + err.Error(),
	}
}
