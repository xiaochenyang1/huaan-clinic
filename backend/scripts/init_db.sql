-- ====================================================
-- 华安医疗预约系统 - 数据库初始化脚本
-- ====================================================

-- 1. 创建数据库
DROP DATABASE IF EXISTS huaan_medical;
CREATE DATABASE huaan_medical CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 2. 使用数据库
USE huaan_medical;

-- 提示：表结构将由 GORM AutoMigrate 自动创建
-- 本脚本仅用于创建数据库和初始化基础数据
