-- =============================================
-- 迁移脚本：添加多种登录方式支持
-- 创建时间：2025-12-17
-- 说明：支持用户名密码、手机号验证码、微信登录
-- =============================================

USE huaan_medical;

-- 1. 新增字段
ALTER TABLE `users`
ADD COLUMN `username` VARCHAR(64) DEFAULT NULL COMMENT '用户名（用于密码登录）' AFTER `open_id`,
ADD COLUMN `password` VARCHAR(128) DEFAULT NULL COMMENT '密码（bcrypt加密）' AFTER `username`,
ADD COLUMN `login_type` VARCHAR(20) DEFAULT 'wechat' COMMENT '登录类型：wechat/password/phone' AFTER `password`;

-- 2. 修改open_id为可空
ALTER TABLE `users`
MODIFY COLUMN `open_id` VARCHAR(64) DEFAULT NULL COMMENT '微信OpenID';

-- 3. 调整索引
-- 删除原phone普通索引，改为唯一索引
ALTER TABLE `users` DROP INDEX IF EXISTS `idx_phone`;
ALTER TABLE `users` ADD UNIQUE KEY `uk_phone` (`phone`);

-- 新增username唯一索引
ALTER TABLE `users` ADD UNIQUE KEY `uk_username` (`username`);

-- 新增login_type索引（用于统计）
ALTER TABLE `users` ADD INDEX `idx_login_type` (`login_type`);

-- 4. 更新现有数据的login_type为wechat
UPDATE `users` SET `login_type` = 'wechat' WHERE `login_type` IS NULL OR `login_type` = '';

-- 5. 验证迁移结果
SELECT
    COLUMN_NAME,
    DATA_TYPE,
    IS_NULLABLE,
    COLUMN_DEFAULT,
    COLUMN_COMMENT
FROM INFORMATION_SCHEMA.COLUMNS
WHERE TABLE_SCHEMA = 'huaan_medical'
  AND TABLE_NAME = 'users'
  AND COLUMN_NAME IN ('username', 'password', 'login_type', 'open_id');

-- 迁移完成提示
SELECT '迁移完成！已添加 username, password, login_type 字段' AS status;
