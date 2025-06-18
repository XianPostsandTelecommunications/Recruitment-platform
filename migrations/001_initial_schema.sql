-- 实验室招新平台数据库初始化脚本
-- 版本: 1.0.0
-- 创建时间: 2024-01-01

-- 创建数据库
CREATE DATABASE IF NOT EXISTS lab_recruitment CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE lab_recruitment;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '用户ID',
    username VARCHAR(50) UNIQUE NOT NULL COMMENT '用户名',
    email VARCHAR(100) UNIQUE NOT NULL COMMENT '邮箱',
    password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希',
    role ENUM('student', 'admin') DEFAULT 'student' COMMENT '用户角色',
    avatar VARCHAR(255) COMMENT '头像URL',
    phone VARCHAR(20) COMMENT '手机号',
    student_id VARCHAR(20) COMMENT '学号',
    major VARCHAR(100) COMMENT '专业',
    grade VARCHAR(20) COMMENT '年级',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_email (email),
    INDEX idx_role (role),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 实验室表
CREATE TABLE IF NOT EXISTS labs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '实验室ID',
    name VARCHAR(100) NOT NULL COMMENT '实验室名称',
    description TEXT COMMENT '实验室描述',
    requirements TEXT COMMENT '招新要求',
    max_members INT DEFAULT 10 COMMENT '最大成员数',
    current_members INT DEFAULT 0 COMMENT '当前成员数',
    contact_email VARCHAR(100) COMMENT '联系邮箱',
    contact_phone VARCHAR(20) COMMENT '联系电话',
    location VARCHAR(200) COMMENT '实验室位置',
    tags JSON COMMENT '标签数组',
    cover_image VARCHAR(255) COMMENT '封面图片',
    status ENUM('active', 'inactive') DEFAULT 'active' COMMENT '状态',
    created_by BIGINT COMMENT '创建者ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
    INDEX idx_status (status),
    INDEX idx_created_by (created_by),
    INDEX idx_created_at (created_at),
    FULLTEXT idx_search (name, description, requirements)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='实验室表';

-- 申请表
CREATE TABLE IF NOT EXISTS applications (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '申请ID',
    user_id BIGINT NOT NULL COMMENT '申请人ID',
    lab_id BIGINT NOT NULL COMMENT '实验室ID',
    motivation TEXT COMMENT '申请动机',
    skills JSON COMMENT '技能列表',
    experience TEXT COMMENT '相关经验',
    available_time VARCHAR(100) COMMENT '可用时间',
    resume_url VARCHAR(255) COMMENT '简历文件URL',
    status ENUM('pending', 'accepted', 'rejected') DEFAULT 'pending' COMMENT '申请状态',
    feedback TEXT COMMENT '审核反馈',
    reviewed_by BIGINT COMMENT '审核人ID',
    reviewed_at TIMESTAMP NULL COMMENT '审核时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (lab_id) REFERENCES labs(id) ON DELETE CASCADE,
    FOREIGN KEY (reviewed_by) REFERENCES users(id) ON DELETE SET NULL,
    INDEX idx_user_id (user_id),
    INDEX idx_lab_id (lab_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at),
    UNIQUE KEY uk_user_lab (user_id, lab_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='申请表';

-- 通知表
CREATE TABLE IF NOT EXISTS notifications (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '通知ID',
    user_id BIGINT NOT NULL COMMENT '接收用户ID',
    title VARCHAR(200) NOT NULL COMMENT '通知标题',
    content TEXT COMMENT '通知内容',
    type ENUM('system', 'application', 'lab') DEFAULT 'system' COMMENT '通知类型',
    is_read BOOLEAN DEFAULT FALSE COMMENT '是否已读',
    related_id BIGINT COMMENT '关联ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_is_read (is_read),
    INDEX idx_created_at (created_at),
    INDEX idx_type (type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='通知表';

-- 文件上传表
CREATE TABLE IF NOT EXISTS file_uploads (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '文件ID',
    filename VARCHAR(255) NOT NULL COMMENT '文件名',
    original_name VARCHAR(255) NOT NULL COMMENT '原始文件名',
    file_path VARCHAR(500) NOT NULL COMMENT '文件路径',
    file_size BIGINT NOT NULL COMMENT '文件大小(字节)',
    mime_type VARCHAR(100) NOT NULL COMMENT 'MIME类型',
    uploader_id BIGINT NOT NULL COMMENT '上传者ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    FOREIGN KEY (uploader_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_uploader_id (uploader_id),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文件上传表';

-- 系统配置表
CREATE TABLE IF NOT EXISTS system_configs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '配置ID',
    config_key VARCHAR(100) UNIQUE NOT NULL COMMENT '配置键',
    config_value TEXT COMMENT '配置值',
    description VARCHAR(255) COMMENT '配置描述',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_config_key (config_key)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统配置表';

-- 操作日志表
CREATE TABLE IF NOT EXISTS operation_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '日志ID',
    user_id BIGINT COMMENT '操作用户ID',
    action VARCHAR(100) NOT NULL COMMENT '操作动作',
    resource_type VARCHAR(50) COMMENT '资源类型',
    resource_id BIGINT COMMENT '资源ID',
    details JSON COMMENT '操作详情',
    ip_address VARCHAR(45) COMMENT 'IP地址',
    user_agent TEXT COMMENT '用户代理',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
    INDEX idx_user_id (user_id),
    INDEX idx_action (action),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='操作日志表';

-- 插入默认管理员用户
INSERT INTO users (username, email, password_hash, role, created_at) VALUES 
('admin', 'admin@lab-recruitment.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin', NOW())
ON DUPLICATE KEY UPDATE updated_at = NOW();

-- 插入默认系统配置
INSERT INTO system_configs (config_key, config_value, description) VALUES 
('site_name', '实验室招新平台', '网站名称'),
('site_description', '高校实验室招新管理平台', '网站描述'),
('max_file_size', '10485760', '最大文件上传大小(字节)'),
('allowed_file_types', '["jpg","jpeg","png","gif","pdf","doc","docx"]', '允许上传的文件类型'),
('application_deadline', '2024-12-31', '申请截止日期'),
('notification_enabled', 'true', '是否启用通知功能')
ON DUPLICATE KEY UPDATE config_value = VALUES(config_value), updated_at = NOW();

-- 创建视图：实验室申请统计
CREATE OR REPLACE VIEW lab_application_stats AS
SELECT 
    l.id as lab_id,
    l.name as lab_name,
    COUNT(a.id) as total_applications,
    SUM(CASE WHEN a.status = 'pending' THEN 1 ELSE 0 END) as pending_applications,
    SUM(CASE WHEN a.status = 'accepted' THEN 1 ELSE 0 END) as accepted_applications,
    SUM(CASE WHEN a.status = 'rejected' THEN 1 ELSE 0 END) as rejected_applications,
    CASE 
        WHEN COUNT(a.id) > 0 THEN 
            ROUND(SUM(CASE WHEN a.status = 'accepted' THEN 1 ELSE 0 END) * 100.0 / COUNT(a.id), 2)
        ELSE 0 
    END as acceptance_rate
FROM labs l
LEFT JOIN applications a ON l.id = a.lab_id
WHERE l.status = 'active'
GROUP BY l.id, l.name;

-- 创建视图：用户申请统计
CREATE OR REPLACE VIEW user_application_stats AS
SELECT 
    u.id as user_id,
    u.username,
    u.email,
    COUNT(a.id) as total_applications,
    SUM(CASE WHEN a.status = 'pending' THEN 1 ELSE 0 END) as pending_applications,
    SUM(CASE WHEN a.status = 'accepted' THEN 1 ELSE 0 END) as accepted_applications,
    SUM(CASE WHEN a.status = 'rejected' THEN 1 ELSE 0 END) as rejected_applications,
    MAX(a.created_at) as last_application_date
FROM users u
LEFT JOIN applications a ON u.id = a.user_id
WHERE u.role = 'student'
GROUP BY u.id, u.username, u.email; 