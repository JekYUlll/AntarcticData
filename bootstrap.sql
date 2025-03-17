-- 创建数据库 antarctic_data
CREATE DATABASE IF NOT EXISTS antarctic_data;
USE antarctic_data;

-- 创建表 weather_data
CREATE TABLE IF NOT EXISTS weather_data (
    id INT AUTO_INCREMENT PRIMARY KEY, -- 添加自增主键
    station VARCHAR(255) NOT NULL,     -- 站点名称，非空
    record_time DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3), -- 记录时间，精确到毫秒
    temperature DOUBLE NOT NULL,       -- 温度，非空
    humidity INT NOT NULL,             -- 湿度，非空
    wind_dir INT NOT NULL,             -- 风向，非空
    wind_speed DOUBLE NOT NULL,        -- 风速，非空
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3), -- 创建时间，精确到毫秒
    INDEX idx_station (station)        -- 为 station 字段添加索引
);