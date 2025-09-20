-- 德州扑克数据库表结构创建脚本
-- 作用：创建所有必要的数据库表，包括用户、房间、游戏、管理员等

-- 创建数据库
CREATE DATABASE IF NOT EXISTS texas_poker DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE texas_poker;

-- 用户表
CREATE TABLE users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL COMMENT '用户名',
    email VARCHAR(100) UNIQUE NOT NULL COMMENT '邮箱',
    password_hash VARCHAR(255) NOT NULL COMMENT '加密后的密码',
    chips INT DEFAULT 1000 COMMENT '当前筹码数',
    total_games INT DEFAULT 0 COMMENT '总游戏局数',
    total_wins INT DEFAULT 0 COMMENT '总获胜局数',
    avatar_url VARCHAR(255) COMMENT '头像URL',
    status ENUM('active', 'disabled') DEFAULT 'active' COMMENT '账号状态',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_username (username),
    INDEX idx_email (email),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 房间表
CREATE TABLE rooms (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL COMMENT '房间名称',
    chip_level ENUM('low', 'medium', 'high') NOT NULL COMMENT '筹码级别',
    min_chips INT NOT NULL COMMENT '进入最低筹码要求',
    small_blind INT NOT NULL COMMENT '小盲注',
    big_blind INT NOT NULL COMMENT '大盲注',
    max_players INT DEFAULT 6 COMMENT '最大玩家数',
    is_private BOOLEAN DEFAULT FALSE COMMENT '是否私人房间',
    password_hash VARCHAR(255) COMMENT '私人房间密码哈希',
    status ENUM('waiting', 'playing', 'closed') DEFAULT 'waiting' COMMENT '房间状态',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_chip_level (chip_level),
    INDEX idx_status (status),
    INDEX idx_is_private (is_private)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='房间表';

-- 游戏记录表
CREATE TABLE games (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    room_id BIGINT NOT NULL COMMENT '所属房间ID',
    winner_id BIGINT COMMENT '获胜者ID',
    pot_amount INT NOT NULL COMMENT '底池金额',
    start_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '开始时间',
    end_time TIMESTAMP NULL COMMENT '结束时间',
    game_log JSON COMMENT '游戏日志(JSON格式)',
    FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE,
    FOREIGN KEY (winner_id) REFERENCES users(id) ON DELETE SET NULL,
    INDEX idx_room_id (room_id),
    INDEX idx_winner_id (winner_id),
    INDEX idx_start_time (start_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='游戏记录表';

-- 游戏玩家表
CREATE TABLE game_players (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    game_id BIGINT NOT NULL COMMENT '游戏ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    chips_change INT NOT NULL COMMENT '筹码变化',
    position INT NOT NULL COMMENT '座位位置',
    FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_game_id (game_id),
    INDEX idx_user_id (user_id),
    UNIQUE KEY uk_game_user (game_id, user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='游戏玩家表';

-- 管理员表
CREATE TABLE admins (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL COMMENT '管理员用户名',
    email VARCHAR(100) UNIQUE NOT NULL COMMENT '邮箱',
    password_hash VARCHAR(255) NOT NULL COMMENT '加密后的密码',
    role ENUM('super', 'normal') DEFAULT 'normal' COMMENT '角色',
    status ENUM('active', 'inactive') DEFAULT 'active' COMMENT '状态',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_username (username),
    INDEX idx_role (role),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='管理员表';

-- 插入默认管理员账号
INSERT INTO admins (username, email, password_hash, role) VALUES 
('admin', 'admin@texaspoker.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'super'); 