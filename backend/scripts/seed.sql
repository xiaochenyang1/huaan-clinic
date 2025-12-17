-- =============================================
-- 华安医疗预约系统 - 初始化数据脚本
-- =============================================

USE huaan_medical;

-- =============================================
-- 1. 初始化角色
-- =============================================
INSERT INTO `roles` (`code`, `name`, `description`, `sort_order`, `status`) VALUES
('super_admin', '超级管理员', '拥有系统所有权限', 1, 1),
('admin', '管理员', '拥有大部分管理权限', 2, 1),
('operator', '操作员', '拥有基本操作权限', 3, 1);

-- =============================================
-- 2. 初始化权限
-- =============================================
INSERT INTO `permissions` (`code`, `name`, `module`, `description`, `sort_order`) VALUES
-- 预约管理
('appointment:list', '预约列表', 'appointment', '查看预约列表', 1),
('appointment:detail', '预约详情', 'appointment', '查看预约详情', 2),
('appointment:update', '更新预约', 'appointment', '更新预约状态', 3),
('appointment:export', '导出预约', 'appointment', '导出预约数据', 4),

-- 患者管理
('patient:list', '患者列表', 'patient', '查看患者列表', 1),
('patient:detail', '患者详情', 'patient', '查看患者详情', 2),

-- 科室管理
('department:list', '科室列表', 'department', '查看科室列表', 1),
('department:create', '添加科室', 'department', '添加新科室', 2),
('department:update', '编辑科室', 'department', '编辑科室信息', 3),
('department:delete', '删除科室', 'department', '删除科室', 4),

-- 医生管理
('doctor:list', '医生列表', 'doctor', '查看医生列表', 1),
('doctor:create', '添加医生', 'doctor', '添加新医生', 2),
('doctor:update', '编辑医生', 'doctor', '编辑医生信息', 3),
('doctor:delete', '删除医生', 'doctor', '删除医生', 4),

-- 排班管理
('schedule:list', '排班列表', 'schedule', '查看排班列表', 1),
('schedule:create', '创建排班', 'schedule', '创建新排班', 2),
('schedule:update', '编辑排班', 'schedule', '编辑排班信息', 3),
('schedule:delete', '删除排班', 'schedule', '删除排班', 4),
('schedule:batch', '批量排班', 'schedule', '批量创建排班', 5),

-- 数据统计
('statistics:view', '查看统计', 'statistics', '查看数据统计', 1),
('statistics:export', '导出统计', 'statistics', '导出统计数据', 2),

-- 系统管理
('system:admin', '管理员管理', 'system', '管理员增删改查', 1),
('system:role', '角色管理', 'system', '角色增删改查', 2),
('system:log', '日志查看', 'system', '查看操作日志', 3);

-- =============================================
-- 3. 初始化角色权限关联
-- =============================================
-- 超级管理员拥有所有权限
INSERT INTO `role_permissions` (`role_id`, `permission_id`)
SELECT 1, id FROM `permissions`;

-- 管理员拥有除系统管理外的所有权限
INSERT INTO `role_permissions` (`role_id`, `permission_id`)
SELECT 2, id FROM `permissions` WHERE `module` != 'system';

-- 操作员拥有查看和基本操作权限
INSERT INTO `role_permissions` (`role_id`, `permission_id`)
SELECT 3, id FROM `permissions`
WHERE `code` IN (
    'appointment:list', 'appointment:detail', 'appointment:update',
    'patient:list', 'patient:detail',
    'department:list',
    'doctor:list',
    'schedule:list',
    'statistics:view'
);

-- =============================================
-- 4. 初始化超级管理员
-- 默认密码: admin123 (BCrypt加密)
-- =============================================
INSERT INTO `admins` (`username`, `password`, `nickname`, `status`) VALUES
('admin', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iAt6Z5EH', '超级管理员', 1);

-- 关联超级管理员角色
INSERT INTO `admin_roles` (`admin_id`, `role_id`) VALUES (1, 1);

-- =============================================
-- 5. 初始化科室数据
-- =============================================
INSERT INTO `departments` (`name`, `description`, `icon`, `sort_order`, `status`) VALUES
('内科', '诊治内脏器官疾病，包括心血管、呼吸、消化、内分泌等系统疾病', '/static/icons/neike.png', 1, 1),
('外科', '诊治需要手术治疗的疾病，包括普外、骨科、泌尿外科等', '/static/icons/waike.png', 2, 1),
('妇产科', '女性生殖系统疾病诊治及产前检查、分娩服务', '/static/icons/fuke.png', 3, 1),
('儿科', '0-14岁儿童常见病、多发病的诊治', '/static/icons/erke.png', 4, 1),
('眼科', '眼部疾病的诊断和治疗，包括近视、白内障等', '/static/icons/yanke.png', 5, 1),
('口腔科', '牙齿及口腔疾病的诊治，包括补牙、拔牙、正畸等', '/static/icons/kouqiang.png', 6, 1),
('皮肤科', '皮肤病的诊断和治疗，包括湿疹、皮炎、痤疮等', '/static/icons/pifu.png', 7, 1),
('中医科', '运用中医理论和方法诊治疾病', '/static/icons/zhongyi.png', 8, 1);

-- =============================================
-- 6. 初始化医生数据
-- =============================================
INSERT INTO `doctors` (`department_id`, `name`, `title`, `specialty`, `introduction`, `sort_order`, `status`) VALUES
-- 内科医生
(1, '张明华', 'chief_physician', '心血管疾病、高血压、冠心病', '主任医师，从事内科临床工作30余年，擅长心血管疾病的诊治。', 1, 1),
(1, '李秀英', 'associate_chief_physician', '呼吸系统疾病、慢性支气管炎', '副主任医师，呼吸内科专家，对慢性呼吸道疾病有丰富经验。', 2, 1),
(1, '王建国', 'attending_physician', '消化系统疾病、胃肠炎', '主治医师，擅长消化系统常见病的诊治。', 3, 1),

-- 外科医生
(2, '刘志强', 'chief_physician', '普通外科、腹腔镜手术', '主任医师，普外科专家，擅长微创手术。', 1, 1),
(2, '陈伟', 'associate_chief_physician', '骨科、关节置换', '副主任医师，骨科专家，擅长关节疾病治疗。', 2, 1),

-- 妇产科医生
(3, '赵丽娟', 'chief_physician', '妇科肿瘤、宫颈疾病', '主任医师，妇产科专家，从事临床工作25年。', 1, 1),
(3, '孙晓燕', 'associate_chief_physician', '产科、高危妊娠', '副主任医师，产科专家，擅长高危妊娠管理。', 2, 1),

-- 儿科医生
(4, '周小明', 'associate_chief_physician', '小儿呼吸道疾病、过敏性疾病', '副主任医师，儿科专家，擅长儿童常见病诊治。', 1, 1),
(4, '吴婷婷', 'attending_physician', '新生儿疾病、儿童保健', '主治医师，新生儿专家。', 2, 1),

-- 眼科医生
(5, '郑光明', 'chief_physician', '白内障手术、青光眼', '主任医师，眼科专家，完成白内障手术3000余例。', 1, 1),
(5, '钱视力', 'attending_physician', '近视防控、眼底疾病', '主治医师，擅长青少年近视防控。', 2, 1),

-- 口腔科医生
(6, '孙牙健', 'associate_chief_physician', '口腔种植、牙周病', '副主任医师，口腔种植专家。', 1, 1),
(6, '李洁白', 'attending_physician', '牙齿美容、正畸治疗', '主治医师，擅长牙齿美容修复。', 2, 1),

-- 皮肤科医生
(7, '王美肤', 'associate_chief_physician', '过敏性皮肤病、皮肤美容', '副主任医师，皮肤科专家。', 1, 1),
(7, '张清颜', 'attending_physician', '痤疮、湿疹、皮炎', '主治医师，擅长青春痘治疗。', 2, 1),

-- 中医科医生
(8, '李中和', 'chief_physician', '中医内科、慢性病调理', '主任医师，中医专家，从事中医临床40年。', 1, 1),
(8, '陈养生', 'attending_physician', '针灸推拿、颈肩腰腿痛', '主治医师，擅长针灸治疗。', 2, 1);

-- =============================================
-- 7. 完成提示
-- =============================================
SELECT '初始化数据完成！' AS message;
SELECT '默认管理员账号: admin' AS admin_username;
SELECT '默认管理员密码: admin123' AS admin_password;
