-- =============================================
-- 华安医疗预约系统 - 数据库初始化脚本
-- 数据库：MySQL 8.0+
-- 字符集：utf8mb4
-- =============================================

-- 创建数据库
CREATE DATABASE IF NOT EXISTS huaan_medical
DEFAULT CHARACTER SET utf8mb4
COLLATE utf8mb4_unicode_ci;

USE huaan_medical;

-- =============================================
-- 1. 用户表
-- =============================================
CREATE TABLE IF NOT EXISTS `users` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '用户ID',
    `open_id` VARCHAR(64) NOT NULL COMMENT '微信OpenID',
    `union_id` VARCHAR(64) DEFAULT NULL COMMENT '微信UnionID',
    `nickname` VARCHAR(64) DEFAULT NULL COMMENT '昵称',
    `avatar` VARCHAR(512) DEFAULT NULL COMMENT '头像URL',
    `phone` VARCHAR(20) DEFAULT NULL COMMENT '手机号',
    `gender` TINYINT DEFAULT 0 COMMENT '性别 0未知 1男 2女',
    `status` TINYINT DEFAULT 1 COMMENT '状态 0禁用 1启用',
    `blocked_until` DATETIME DEFAULT NULL COMMENT '封禁截止时间',
    `missed_count` INT DEFAULT 0 COMMENT '累计爽约次数',
    `last_login_at` DATETIME DEFAULT NULL COMMENT '最后登录时间',
    `last_login_ip` VARCHAR(64) DEFAULT NULL COMMENT '最后登录IP',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_open_id` (`open_id`),
    KEY `idx_phone` (`phone`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- =============================================
-- 2. 就诊人表
-- =============================================
CREATE TABLE IF NOT EXISTS `patients` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '就诊人ID',
    `user_id` BIGINT NOT NULL COMMENT '所属用户ID',
    `name` VARCHAR(32) NOT NULL COMMENT '姓名',
    `id_card` VARCHAR(18) NOT NULL COMMENT '身份证号',
    `phone` VARCHAR(20) NOT NULL COMMENT '手机号',
    `gender` TINYINT DEFAULT 0 COMMENT '性别 0未知 1男 2女',
    `birth_date` DATE DEFAULT NULL COMMENT '出生日期',
    `relation` VARCHAR(20) DEFAULT 'self' COMMENT '与用户关系',
    `is_default` TINYINT DEFAULT 0 COMMENT '是否默认 0否 1是',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_id_card` (`id_card`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='就诊人表';

-- =============================================
-- 3. 科室表
-- =============================================
CREATE TABLE IF NOT EXISTS `departments` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '科室ID',
    `name` VARCHAR(64) NOT NULL COMMENT '科室名称',
    `description` VARCHAR(512) DEFAULT NULL COMMENT '科室描述',
    `icon` VARCHAR(256) DEFAULT NULL COMMENT '科室图标',
    `sort_order` INT DEFAULT 0 COMMENT '排序序号',
    `status` TINYINT DEFAULT 1 COMMENT '状态 0停用 1启用',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_status` (`status`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='科室表';

-- =============================================
-- 4. 医生表
-- =============================================
CREATE TABLE IF NOT EXISTS `doctors` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '医生ID',
    `department_id` BIGINT NOT NULL COMMENT '所属科室ID',
    `name` VARCHAR(32) NOT NULL COMMENT '姓名',
    `avatar` VARCHAR(512) DEFAULT NULL COMMENT '头像URL',
    `title` VARCHAR(32) NOT NULL COMMENT '职称',
    `specialty` VARCHAR(256) DEFAULT NULL COMMENT '擅长领域',
    `introduction` TEXT DEFAULT NULL COMMENT '个人简介',
    `sort_order` INT DEFAULT 0 COMMENT '排序序号',
    `status` TINYINT DEFAULT 1 COMMENT '状态 0停诊 1正常',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_department_id` (`department_id`),
    KEY `idx_status` (`status`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='医生表';

-- =============================================
-- 5. 排班表
-- =============================================
CREATE TABLE IF NOT EXISTS `schedules` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '排班ID',
    `doctor_id` BIGINT NOT NULL COMMENT '医生ID',
    `schedule_date` DATE NOT NULL COMMENT '排班日期',
    `period` VARCHAR(20) NOT NULL COMMENT '时段 morning/afternoon',
    `start_time` VARCHAR(10) NOT NULL COMMENT '开始时间',
    `end_time` VARCHAR(10) NOT NULL COMMENT '结束时间',
    `total_slots` INT NOT NULL COMMENT '总号源数',
    `available_slots` INT NOT NULL COMMENT '剩余号源数',
    `status` TINYINT DEFAULT 1 COMMENT '状态 0停诊 1正常',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_doctor_id` (`doctor_id`),
    KEY `idx_schedule_date` (`schedule_date`),
    KEY `idx_doctor_date` (`doctor_id`, `schedule_date`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='排班表';

-- =============================================
-- 6. 预约表
-- =============================================
CREATE TABLE IF NOT EXISTS `appointments` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '预约ID',
    `appointment_no` VARCHAR(32) NOT NULL COMMENT '预约编号',
    `user_id` BIGINT NOT NULL COMMENT '用户ID',
    `patient_id` BIGINT NOT NULL COMMENT '就诊人ID',
    `doctor_id` BIGINT NOT NULL COMMENT '医生ID',
    `department_id` BIGINT NOT NULL COMMENT '科室ID',
    `schedule_id` BIGINT NOT NULL COMMENT '排班ID',
    `appointment_date` DATE NOT NULL COMMENT '预约日期',
    `period` VARCHAR(20) NOT NULL COMMENT '时段',
    `appointment_time` VARCHAR(10) NOT NULL COMMENT '预约时间',
    `slot_number` INT NOT NULL COMMENT '号序',
    `status` VARCHAR(20) DEFAULT 'pending' COMMENT '状态',
    `symptom` VARCHAR(512) DEFAULT NULL COMMENT '症状描述',
    `cancel_reason` VARCHAR(256) DEFAULT NULL COMMENT '取消原因',
    `cancelled_at` DATETIME DEFAULT NULL COMMENT '取消时间',
    `checked_in_at` DATETIME DEFAULT NULL COMMENT '签到时间',
    `completed_at` DATETIME DEFAULT NULL COMMENT '完成时间',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_appointment_no` (`appointment_no`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_patient_id` (`patient_id`),
    KEY `idx_doctor_id` (`doctor_id`),
    KEY `idx_department_id` (`department_id`),
    KEY `idx_schedule_id` (`schedule_id`),
    KEY `idx_appointment_date` (`appointment_date`),
    KEY `idx_status` (`status`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='预约表';

-- =============================================
-- 7. 就诊记录表
-- =============================================
CREATE TABLE IF NOT EXISTS `medical_records` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '记录ID',
    `appointment_id` BIGINT NOT NULL COMMENT '预约ID',
    `patient_id` BIGINT NOT NULL COMMENT '患者ID',
    `doctor_id` BIGINT NOT NULL COMMENT '医生ID',
    `department_id` BIGINT NOT NULL COMMENT '科室ID',
    `visit_date` DATE NOT NULL COMMENT '就诊日期',
    `diagnosis` TEXT DEFAULT NULL COMMENT '诊断结果',
    `prescription` TEXT DEFAULT NULL COMMENT '处方',
    `advice` TEXT DEFAULT NULL COMMENT '医嘱',
    `remark` VARCHAR(512) DEFAULT NULL COMMENT '备注',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_appointment_id` (`appointment_id`),
    KEY `idx_patient_id` (`patient_id`),
    KEY `idx_doctor_id` (`doctor_id`),
    KEY `idx_department_id` (`department_id`),
    KEY `idx_visit_date` (`visit_date`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='就诊记录表';

-- =============================================
-- 8. 管理员表
-- =============================================
CREATE TABLE IF NOT EXISTS `admins` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '管理员ID',
    `username` VARCHAR(64) NOT NULL COMMENT '用户名',
    `password` VARCHAR(128) NOT NULL COMMENT '密码',
    `nickname` VARCHAR(64) DEFAULT NULL COMMENT '昵称',
    `avatar` VARCHAR(512) DEFAULT NULL COMMENT '头像',
    `phone` VARCHAR(20) DEFAULT NULL COMMENT '手机号',
    `email` VARCHAR(128) DEFAULT NULL COMMENT '邮箱',
    `status` TINYINT DEFAULT 1 COMMENT '状态 0禁用 1启用',
    `last_login_at` DATETIME DEFAULT NULL COMMENT '最后登录时间',
    `last_login_ip` VARCHAR(64) DEFAULT NULL COMMENT '最后登录IP',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_username` (`username`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='管理员表';

-- =============================================
-- 9. 角色表
-- =============================================
CREATE TABLE IF NOT EXISTS `roles` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '角色ID',
    `code` VARCHAR(32) NOT NULL COMMENT '角色编码',
    `name` VARCHAR(64) NOT NULL COMMENT '角色名称',
    `description` VARCHAR(256) DEFAULT NULL COMMENT '角色描述',
    `sort_order` INT DEFAULT 0 COMMENT '排序序号',
    `status` TINYINT DEFAULT 1 COMMENT '状态 0禁用 1启用',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_code` (`code`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';

-- =============================================
-- 10. 权限表
-- =============================================
CREATE TABLE IF NOT EXISTS `permissions` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '权限ID',
    `code` VARCHAR(64) NOT NULL COMMENT '权限编码',
    `name` VARCHAR(64) NOT NULL COMMENT '权限名称',
    `module` VARCHAR(32) NOT NULL COMMENT '所属模块',
    `description` VARCHAR(256) DEFAULT NULL COMMENT '权限描述',
    `sort_order` INT DEFAULT 0 COMMENT '排序序号',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_code` (`code`),
    KEY `idx_module` (`module`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='权限表';

-- =============================================
-- 11. 管理员角色关联表
-- =============================================
CREATE TABLE IF NOT EXISTS `admin_roles` (
    `admin_id` BIGINT NOT NULL COMMENT '管理员ID',
    `role_id` BIGINT NOT NULL COMMENT '角色ID',
    PRIMARY KEY (`admin_id`, `role_id`),
    KEY `idx_role_id` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='管理员角色关联表';

-- =============================================
-- 12. 角色权限关联表
-- =============================================
CREATE TABLE IF NOT EXISTS `role_permissions` (
    `role_id` BIGINT NOT NULL COMMENT '角色ID',
    `permission_id` BIGINT NOT NULL COMMENT '权限ID',
    PRIMARY KEY (`role_id`, `permission_id`),
    KEY `idx_permission_id` (`permission_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色权限关联表';

-- =============================================
-- 13. 操作日志表
-- =============================================
CREATE TABLE IF NOT EXISTS `operation_logs` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '日志ID',
    `admin_id` BIGINT DEFAULT NULL COMMENT '管理员ID',
    `admin_name` VARCHAR(64) DEFAULT NULL COMMENT '管理员名称',
    `module` VARCHAR(32) DEFAULT NULL COMMENT '操作模块',
    `action` VARCHAR(32) DEFAULT NULL COMMENT '操作动作',
    `method` VARCHAR(10) DEFAULT NULL COMMENT '请求方法',
    `path` VARCHAR(256) DEFAULT NULL COMMENT '请求路径',
    `query` TEXT DEFAULT NULL COMMENT '请求参数',
    `body` TEXT DEFAULT NULL COMMENT '请求体',
    `response` TEXT DEFAULT NULL COMMENT '响应内容',
    `ip` VARCHAR(64) DEFAULT NULL COMMENT 'IP地址',
    `user_agent` VARCHAR(512) DEFAULT NULL COMMENT 'User-Agent',
    `status` INT DEFAULT NULL COMMENT '响应状态码',
    `latency` BIGINT DEFAULT NULL COMMENT '耗时(ms)',
    `error_msg` TEXT DEFAULT NULL COMMENT '错误信息',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_admin_id` (`admin_id`),
    KEY `idx_module` (`module`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='操作日志表';

-- =============================================
-- 14. 登录日志表
-- =============================================
CREATE TABLE IF NOT EXISTS `login_logs` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '日志ID',
    `user_type` VARCHAR(20) DEFAULT NULL COMMENT '用户类型',
    `user_id` BIGINT DEFAULT NULL COMMENT '用户ID',
    `username` VARCHAR(64) DEFAULT NULL COMMENT '用户名',
    `login_type` VARCHAR(20) DEFAULT NULL COMMENT '登录方式',
    `ip` VARCHAR(64) DEFAULT NULL COMMENT 'IP地址',
    `location` VARCHAR(128) DEFAULT NULL COMMENT '登录地点',
    `device` VARCHAR(256) DEFAULT NULL COMMENT '设备信息',
    `os` VARCHAR(64) DEFAULT NULL COMMENT '操作系统',
    `browser` VARCHAR(64) DEFAULT NULL COMMENT '浏览器',
    `status` TINYINT DEFAULT NULL COMMENT '登录状态',
    `message` VARCHAR(256) DEFAULT NULL COMMENT '登录消息',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_type` (`user_type`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='登录日志表';
