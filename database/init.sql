CREATE DATABASE IF NOT EXISTS carbon_tracker DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE carbon_tracker;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(100),
    region VARCHAR(50) DEFAULT '全国',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB;

-- 碳排因子库
CREATE TABLE IF NOT EXISTS emission_factors (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    category VARCHAR(20) NOT NULL COMMENT '分类: transport/electricity/diet',
    item VARCHAR(50) NOT NULL COMMENT '具体项目',
    factor DECIMAL(10,4) NOT NULL COMMENT '碳排因子 (kgCO2/单位)',
    unit VARCHAR(20) NOT NULL COMMENT '单位',
    description VARCHAR(200),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB;

-- 每日碳记录
CREATE TABLE IF NOT EXISTS carbon_records (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    record_date DATE NOT NULL,
    category VARCHAR(20) NOT NULL COMMENT 'transport/electricity/diet',
    item VARCHAR(50) NOT NULL,
    amount DECIMAL(10,2) NOT NULL COMMENT '数量',
    emission DECIMAL(10,4) NOT NULL COMMENT '计算后的碳排放量(kgCO2)',
    note VARCHAR(200),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_date (user_id, record_date)
) ENGINE=InnoDB;

-- 减排目标表
CREATE TABLE IF NOT EXISTS reduction_goals (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    target_emission DECIMAL(10,2) NOT NULL COMMENT '目标碳排放量(kgCO2/天)',
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    status VARCHAR(20) DEFAULT 'active' COMMENT 'active/completed/failed',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB;

-- 区域平均值表
CREATE TABLE IF NOT EXISTS regional_averages (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    region VARCHAR(50) NOT NULL,
    month VARCHAR(7) NOT NULL COMMENT '格式: 2024-01',
    avg_emission DECIMAL(10,2) NOT NULL COMMENT '日均碳排放(kgCO2)',
    UNIQUE KEY uk_region_month (region, month)
) ENGINE=InnoDB;

-- 初始化碳排因子数据
INSERT INTO emission_factors (category, item, factor, unit, description) VALUES
-- 出行方式
('transport', '私家车', 0.21, 'km', '小汽车每公里排放约0.21kgCO2'),
('transport', '公交车', 0.08, 'km', '公交车每人每公里排放约0.08kgCO2'),
('transport', '地铁', 0.05, 'km', '地铁每人每公里排放约0.05kgCO2'),
('transport', '骑行', 0.00, 'km', '自行车零排放'),
('transport', '步行', 0.00, 'km', '步行零排放'),
('transport', '出租车', 0.25, 'km', '出租车每公里排放约0.25kgCO2'),
('transport', '飞机', 0.26, 'km', '飞机每人每公里排放约0.26kgCO2'),
-- 用电
('electricity', '家庭用电', 0.785, 'kWh', '中国电网平均碳排因子'),
('electricity', '办公用电', 0.785, 'kWh', '中国电网平均碳排因子'),
-- 饮食
('diet', '牛肉', 27.0, 'kg', '每公斤牛肉约产生27kgCO2'),
('diet', '猪肉', 12.1, 'kg', '每公斤猪肉约产生12.1kgCO2'),
('diet', '鸡肉', 6.9, 'kg', '每公斤鸡肉约产生6.9kgCO2'),
('diet', '鱼类', 6.1, 'kg', '每公斤鱼类约产生6.1kgCO2'),
('diet', '蛋奶', 4.8, 'kg', '每公斤蛋奶约产生4.8kgCO2'),
('diet', '蔬菜', 2.0, 'kg', '每公斤蔬菜约产生2.0kgCO2'),
('diet', '米饭', 2.7, 'kg', '每公斤米饭约产生2.7kgCO2');

-- 初始化区域平均值(模拟数据)
INSERT INTO regional_averages (region, month, avg_emission) VALUES
('全国', '2024-01', 8.5),
('全国', '2024-02', 8.2),
('全国', '2024-03', 7.8),
('全国', '2024-04', 7.5),
('全国', '2024-05', 7.9),
('全国', '2024-06', 8.1),
('北京', '2024-01', 9.2),
('北京', '2024-02', 8.8),
('北京', '2024-03', 8.5),
('北京', '2024-04', 8.0),
('北京', '2024-05', 8.3),
('北京', '2024-06', 8.6),
('上海', '2024-01', 8.8),
('上海', '2024-02', 8.5),
('上海', '2024-03', 8.1),
('上海', '2024-04', 7.7),
('上海', '2024-05', 8.0),
('上海', '2024-06', 8.3),
('广州', '2024-01', 7.9),
('广州', '2024-02', 7.6),
('广州', '2024-03', 7.3),
('广州', '2024-04', 7.0),
('广州', '2024-05', 7.4),
('广州', '2024-06', 7.7);
