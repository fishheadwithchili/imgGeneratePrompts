-- Image Generate Prompts 数据库结构 v2.0
-- 这个文件仅供参考，实际表结构由GORM自动创建
-- 运行 go run cmd/db-manager.go -write 后，GORM会自动执行迁移创建这些表

-- =============================================
-- 1. 标签表 (tags)
-- =============================================
CREATE TABLE `tags` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL COMMENT '标签名称',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_tags_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =============================================
-- 2. 提示词表 (prompts) - 主表
-- =============================================
CREATE TABLE `prompts` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT '软删除时间',
  `prompt_text` longtext NOT NULL COMMENT '正面提示词',
  `negative_prompt` longtext COMMENT '负面提示词',
  `model_name` varchar(100) DEFAULT NULL COMMENT '使用的AI模型名称',
  `image_url` varchar(500) NOT NULL COMMENT '生成图片的存储路径或URL',
  `is_public` tinyint(1) DEFAULT '0' COMMENT '是否公开, 0:不公开, 1:公开',
  
  -- === v2.0 新增字段 ===
  `style_description` varchar(500) DEFAULT NULL COMMENT '风格描述',
  `usage_scenario` varchar(500) DEFAULT NULL COMMENT '适用场景描述',
  `atmosphere_description` varchar(500) DEFAULT NULL COMMENT '氛围描述',
  `expressive_intent` varchar(500) DEFAULT NULL COMMENT '表现意图描述',
  `structure_analysis` json DEFAULT NULL COMMENT '提示词结构分析 (JSON格式)',
  
  PRIMARY KEY (`id`),
  KEY `idx_prompts_deleted_at` (`deleted_at`),
  KEY `idx_prompts_is_public` (`is_public`),
  KEY `idx_prompts_model_name` (`model_name`),
  KEY `idx_prompts_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =============================================
-- 3. 提示词标签关联表 (prompt_tags) - 中间表
-- =============================================
CREATE TABLE `prompt_tags` (
  `prompt_id` bigint unsigned NOT NULL COMMENT '提示词ID',
  `tag_id` bigint unsigned NOT NULL COMMENT '标签ID',
  PRIMARY KEY (`prompt_id`,`tag_id`),
  KEY `fk_prompt_tags_tag` (`tag_id`),
  CONSTRAINT `fk_prompt_tags_prompt` FOREIGN KEY (`prompt_id`) REFERENCES `prompts` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_prompt_tags_tag` FOREIGN KEY (`tag_id`) REFERENCES `tags` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =============================================
-- 示例数据插入 (可选)
-- =============================================

-- 插入示例标签
INSERT INTO `tags` (`name`, `created_at`) VALUES 
('风景', NOW()),
('人物', NOW()),
('动物', NOW()),
('建筑', NOW()),
('抽象', NOW()),
('科幻', NOW()),
('复古', NOW()),
('现代', NOW()),
('暖色调', NOW()),
('冷色调', NOW()),
('高质量', NOW()),
('4K', NOW());

-- 插入示例提示词
INSERT INTO `prompts` (
  `prompt_text`, 
  `negative_prompt`, 
  `model_name`, 
  `image_url`, 
  `is_public`,
  `style_description`,
  `usage_scenario`,
  `atmosphere_description`,
  `expressive_intent`,
  `structure_analysis`,
  `created_at`,
  `updated_at`
) VALUES 
(
  'a beautiful sunset over mountains, golden hour, cinematic lighting, high quality', 
  'ugly, blurry, low quality, pixelated, noise', 
  'stable-diffusion-v1-5', 
  '/uploads/sample_sunset.jpg', 
  1,
  '风景摄影风格，温暖的金色调',
  '适用于自然风光、旅游宣传、背景图片',
  '宁静、温暖、壮观的黄昏氛围',
  '表现大自然的壮美和宁静',
  '{"主体":"山峰日落","光照":"黄金时刻","质量":"高质量","风格":"电影感"}',
  NOW(),
  NOW()
),
(
  'portrait of a cat, professional photography, studio lighting, detailed fur texture', 
  'cartoon, anime, low resolution, distorted', 
  'stable-diffusion-v1-5', 
  '/uploads/sample_cat.jpg', 
  1,
  '专业摄影风格，细致的毛发质感',
  '适用于宠物摄影、动物主题设计',
  '温馨、可爱、专业的摄影氛围',
  '突出动物的可爱特征和毛发细节',
  '{"主体":"猫咪肖像","技法":"专业摄影","光照":"工作室灯光","细节":"毛发质感"}',
  NOW(),
  NOW()
),
(
  'futuristic city skyline, neon lights, cyberpunk style, night scene, high-tech architecture',
  'old, vintage, daylight, low quality',
  'stable-diffusion-xl',
  '/uploads/sample_cyberpunk.jpg',
  1,
  '赛博朋克风格，霓虹灯光效果',
  '适用于科幻题材、游戏背景、未来主题设计',
  '神秘、科技感十足的未来夜景',
  '展现未来科技城市的繁华与神秘',
  '{"主体":"未来城市","风格":"赛博朋克","光效":"霓虹灯","时间":"夜景"}',
  NOW(),
  NOW()
);

-- 插入标签关联关系 (需要在插入提示词后执行)
-- 这些INSERT语句需要根据实际的ID进行调整

-- 示例：为第一个提示词添加标签 (风景、暖色调、高质量、4K)
-- INSERT INTO `prompt_tags` (`prompt_id`, `tag_id`) VALUES 
-- (1, 1), -- 风景
-- (1, 9), -- 暖色调  
-- (1, 11), -- 高质量
-- (1, 12); -- 4K

-- 示例：为第二个提示词添加标签 (动物、现代、高质量)
-- INSERT INTO `prompt_tags` (`prompt_id`, `tag_id`) VALUES 
-- (2, 3), -- 动物
-- (2, 8), -- 现代
-- (2, 11); -- 高质量

-- 示例：为第三个提示词添加标签 (科幻、现代、冷色调、高质量)
-- INSERT INTO `prompt_tags` (`prompt_id`, `tag_id`) VALUES 
-- (3, 6), -- 科幻
-- (3, 8), -- 现代
-- (3, 10), -- 冷色调
-- (3, 11); -- 高质量

-- =============================================
-- 常用查询示例
-- =============================================

-- 1. 获取提示词及其标签
-- SELECT 
--   p.id,
--   p.prompt_text,
--   p.style_description,
--   GROUP_CONCAT(t.name) as tags
-- FROM prompts p
-- LEFT JOIN prompt_tags pt ON p.id = pt.prompt_id
-- LEFT JOIN tags t ON pt.tag_id = t.id
-- WHERE p.deleted_at IS NULL
-- GROUP BY p.id;

-- 2. 查找包含特定标签的提示词
-- SELECT DISTINCT p.*
-- FROM prompts p
-- JOIN prompt_tags pt ON p.id = pt.prompt_id
-- JOIN tags t ON pt.tag_id = t.id
-- WHERE t.name IN ('风景', '高质量')
-- AND p.deleted_at IS NULL;

-- 3. 标签使用统计
-- SELECT 
--   t.name,
--   COUNT(pt.prompt_id) as usage_count
-- FROM tags t
-- LEFT JOIN prompt_tags pt ON t.id = pt.tag_id
-- GROUP BY t.id, t.name
-- ORDER BY usage_count DESC;

-- 4. 模型使用统计
-- SELECT 
--   model_name,
--   COUNT(*) as count
-- FROM prompts
-- WHERE model_name IS NOT NULL 
-- AND model_name != ''
-- AND deleted_at IS NULL
-- GROUP BY model_name
-- ORDER BY count DESC;

-- =============================================
-- 索引优化建议
-- =============================================

-- 如果数据量大，可以考虑添加以下索引：

-- 复合索引用于复杂查询
-- CREATE INDEX idx_prompts_public_model ON prompts(is_public, model_name);
-- CREATE INDEX idx_prompts_public_created ON prompts(is_public, created_at);

-- 全文索引用于文本搜索 (MySQL 5.7+)
-- ALTER TABLE prompts ADD FULLTEXT(prompt_text, negative_prompt, style_description, usage_scenario);
