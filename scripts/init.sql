-- Docker MySQL 初始化脚本
-- 这个脚本会在 MySQL 容器首次启动时执行

-- 设置字符集
SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS `img_generate_prompts` 
CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 使用数据库
USE `img_generate_prompts`;

-- 创建用户（如果不存在）
CREATE USER IF NOT EXISTS 'img_prompts_user'@'%' IDENTIFIED BY 'img_prompts_pass';

-- 授权
GRANT ALL PRIVILEGES ON `img_generate_prompts`.* TO 'img_prompts_user'@'%';
FLUSH PRIVILEGES;

-- 注意：实际的表结构将由 GORM 自动迁移创建
-- 这里只是确保数据库和用户准备就绪

-- 输出确认信息
SELECT 'Database img_generate_prompts initialized successfully' AS status;
