-- =============================================
-- 回滚脚本：删除多种登录方式支持
-- 创建时间：2025-12-17
-- 说明：回滚 migration_add_multi_login.sql 的变更
-- =============================================

USE huaan_medical;

-- 1. 删除新增字段
ALTER TABLE `users` DROP COLUMN IF EXISTS `username`;
ALTER TABLE `users` DROP COLUMN IF EXISTS `password`;
ALTER TABLE `users` DROP COLUMN IF EXISTS `login_type`;

-- 2. 恢复open_id为NOT NULL
ALTER TABLE `users`
MODIFY COLUMN `open_id` VARCHAR(64) NOT NULL COMMENT '微信OpenID';

-- 3. 恢复phone为普通索引
ALTER TABLE `users` DROP INDEX IF EXISTS `uk_phone`;
ALTER TABLE `users` ADD INDEX `idx_phone` (`phone`);

-- 回滚完成提示
SELECT '回滚完成！已删除 username, password, login_type 字段' AS status;
