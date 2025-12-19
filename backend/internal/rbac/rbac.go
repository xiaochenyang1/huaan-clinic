package rbac

// Role codes
const (
	RoleSuperAdmin      = "super_admin"
	RoleHospitalAdmin   = "hospital_admin"
	RoleDepartmentAdmin = "department_admin"
)

// Permission codes (module:action)
const (
	PermDepartmentView   = "department:view"
	PermDepartmentCreate = "department:create"
	PermDepartmentUpdate = "department:update"
	PermDepartmentDelete = "department:delete"

	PermDoctorView   = "doctor:view"
	PermDoctorCreate = "doctor:create"
	PermDoctorUpdate = "doctor:update"
	PermDoctorDelete = "doctor:delete"

	PermScheduleView   = "schedule:view"
	PermScheduleCreate = "schedule:create"
	PermScheduleUpdate = "schedule:update"
	PermScheduleDelete = "schedule:delete"
	PermScheduleBatch  = "schedule:batch"

	PermAppointmentView   = "appointment:view"
	PermAppointmentUpdate = "appointment:update"
	PermAppointmentExport = "appointment:export"

	PermPatientView = "patient:view"

	PermStatisticsView = "statistics:view"

	PermLogView = "log:view"

	PermAdminView     = "admin:view"
	PermAdminCreate   = "admin:create"
	PermAdminUpdate   = "admin:update"
	PermAdminPassword = "admin:password"

	PermRoleView       = "role:view"
	PermRoleCreate     = "role:create"
	PermRoleUpdate     = "role:update"
	PermRolePermission = "role:permission"

	PermPermissionView = "permission:view"

	PermUploadAvatar = "upload:avatar"
	PermUploadImage  = "upload:image"
)

type PermissionDef struct {
	Code        string
	Name        string
	Module      string
	Description string
	SortOrder   int
}

// DefaultPermissions 默认权限清单（用于初始化/对齐权限表）
var DefaultPermissions = []PermissionDef{
	// 科室管理
	{Code: PermDepartmentView, Name: "查看科室", Module: "department", Description: "查看科室列表/详情", SortOrder: 1},
	{Code: PermDepartmentCreate, Name: "创建科室", Module: "department", Description: "创建科室", SortOrder: 2},
	{Code: PermDepartmentUpdate, Name: "编辑科室", Module: "department", Description: "更新科室", SortOrder: 3},
	{Code: PermDepartmentDelete, Name: "删除科室", Module: "department", Description: "删除科室", SortOrder: 4},

	// 医生管理
	{Code: PermDoctorView, Name: "查看医生", Module: "doctor", Description: "查看医生列表/详情", SortOrder: 1},
	{Code: PermDoctorCreate, Name: "创建医生", Module: "doctor", Description: "创建医生", SortOrder: 2},
	{Code: PermDoctorUpdate, Name: "编辑医生", Module: "doctor", Description: "更新医生", SortOrder: 3},
	{Code: PermDoctorDelete, Name: "删除医生", Module: "doctor", Description: "删除医生", SortOrder: 4},

	// 排班管理
	{Code: PermScheduleView, Name: "查看排班", Module: "schedule", Description: "查看排班列表/详情", SortOrder: 1},
	{Code: PermScheduleCreate, Name: "创建排班", Module: "schedule", Description: "创建排班", SortOrder: 2},
	{Code: PermScheduleUpdate, Name: "编辑排班", Module: "schedule", Description: "更新排班", SortOrder: 3},
	{Code: PermScheduleDelete, Name: "删除排班", Module: "schedule", Description: "删除排班", SortOrder: 4},
	{Code: PermScheduleBatch, Name: "批量排班", Module: "schedule", Description: "批量创建排班", SortOrder: 5},

	// 预约管理
	{Code: PermAppointmentView, Name: "查看预约", Module: "appointment", Description: "查看预约列表/详情", SortOrder: 1},
	{Code: PermAppointmentUpdate, Name: "处理预约", Module: "appointment", Description: "更新预约状态", SortOrder: 2},
	{Code: PermAppointmentExport, Name: "导出预约", Module: "appointment", Description: "导出预约数据", SortOrder: 3},

	// 患者管理
	{Code: PermPatientView, Name: "查看患者", Module: "patient", Description: "查看患者列表/详情", SortOrder: 1},

	// 数据统计
	{Code: PermStatisticsView, Name: "查看统计", Module: "statistics", Description: "查看仪表盘/统计数据", SortOrder: 1},

	// 系统日志
	{Code: PermLogView, Name: "查看日志", Module: "log", Description: "查看操作日志/登录日志", SortOrder: 1},

	// 系统管理 - 管理员
	{Code: PermAdminView, Name: "查看管理员", Module: "admin", Description: "查看管理员列表", SortOrder: 1},
	{Code: PermAdminCreate, Name: "创建管理员", Module: "admin", Description: "创建管理员", SortOrder: 2},
	{Code: PermAdminUpdate, Name: "编辑管理员", Module: "admin", Description: "更新管理员信息/角色/状态", SortOrder: 3},
	{Code: PermAdminPassword, Name: "重置管理员密码", Module: "admin", Description: "重置管理员密码", SortOrder: 4},

	// 系统管理 - 角色
	{Code: PermRoleView, Name: "查看角色", Module: "role", Description: "查看角色列表", SortOrder: 1},
	{Code: PermRoleCreate, Name: "创建角色", Module: "role", Description: "创建角色", SortOrder: 2},
	{Code: PermRoleUpdate, Name: "编辑角色", Module: "role", Description: "更新角色信息", SortOrder: 3},
	{Code: PermRolePermission, Name: "分配角色权限", Module: "role", Description: "设置角色权限", SortOrder: 4},

	// 系统管理 - 权限
	{Code: PermPermissionView, Name: "查看权限", Module: "permission", Description: "查看权限清单", SortOrder: 1},

	// 文件上传
	{Code: PermUploadAvatar, Name: "上传头像", Module: "upload", Description: "上传头像文件", SortOrder: 1},
	{Code: PermUploadImage, Name: "上传图片", Module: "upload", Description: "上传图片文件", SortOrder: 2},
}

// AdminRoutePermissions 管理后台接口权限映射：key = "METHOD /api/admin/xxx"
// 空数组表示登录后即可访问（不做权限校验）。
var AdminRoutePermissions = map[string][]string{
	// 管理员基础
	"GET /api/admin/info":     {},
	"PUT /api/admin/password": {},

	// 管理员管理
	"GET /api/admin/admins":             {PermAdminView},
	"POST /api/admin/admins":            {PermAdminCreate},
	"PUT /api/admin/admins/:id":         {PermAdminUpdate},
	"PUT /api/admin/admins/:id/password": {PermAdminPassword},

	// 角色管理
	"GET /api/admin/roles":              {PermRoleView},
	"POST /api/admin/roles":             {PermRoleCreate},
	"PUT /api/admin/roles/:id":          {PermRoleUpdate},
	"PUT /api/admin/roles/:id/permissions": {PermRolePermission},

	// 权限清单
	"GET /api/admin/permissions":    {PermPermissionView},
	"GET /api/admin/permissions/me": {},

	// 仪表盘/统计
	"GET /api/admin/dashboard":  {PermStatisticsView},
	"GET /api/admin/statistics": {PermStatisticsView},

	// 预约管理
	"GET /api/admin/appointments":        {PermAppointmentView},
	"GET /api/admin/appointments/:id":    {PermAppointmentView},
	"PUT /api/admin/appointments/:id":    {PermAppointmentUpdate},
	"GET /api/admin/appointments/export": {PermAppointmentExport},

	// 患者管理
	"GET /api/admin/patients":     {PermPatientView},
	"GET /api/admin/patients/:id": {PermPatientView},

	// 科室管理
	"GET /api/admin/departments":     {PermDepartmentView},
	"GET /api/admin/departments/:id": {PermDepartmentView},
	"POST /api/admin/departments":    {PermDepartmentCreate},
	"PUT /api/admin/departments/:id": {PermDepartmentUpdate},
	"DELETE /api/admin/departments/:id": {PermDepartmentDelete},

	// 医生管理
	"GET /api/admin/doctors":     {PermDoctorView},
	"GET /api/admin/doctors/:id": {PermDoctorView},
	"POST /api/admin/doctors":    {PermDoctorCreate},
	"PUT /api/admin/doctors/:id": {PermDoctorUpdate},
	"DELETE /api/admin/doctors/:id": {PermDoctorDelete},

	// 文件上传
	"POST /api/admin/upload/avatar": {PermUploadAvatar},
	"POST /api/admin/upload/image":  {PermUploadImage},

	// 排班管理
	"GET /api/admin/schedules":       {PermScheduleView},
	"GET /api/admin/schedules/:id":   {PermScheduleView},
	"POST /api/admin/schedules":      {PermScheduleCreate},
	"POST /api/admin/schedules/batch": {PermScheduleBatch},
	"PUT /api/admin/schedules/:id":   {PermScheduleUpdate},
	"DELETE /api/admin/schedules/:id": {PermScheduleDelete},

	// 系统日志
	"GET /api/admin/logs/operation": {PermLogView},
	"GET /api/admin/logs/login":     {PermLogView},
}

