-- 创建数据库 antarctic_data
CREATE DATABASE IF NOT EXISTS antarctic_data;
USE antarctic_data;

-- 删除现有表（如果存在）
DROP TABLE IF EXISTS `weather_data`;

-- 创建新表
CREATE TABLE `weather_data` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `station` varchar(50) NOT NULL,
  `record_time` datetime NOT NULL,
  `temperature` double NOT NULL,
  `humidity` bigint NOT NULL,
  `wind_dir` bigint NOT NULL,
  `wind_speed` double NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `unique_id` varchar(100) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_station_time` (`unique_id`),
  KEY `idx_weather_data_station` (`station`),
  KEY `idx_weather_data_record_time` (`record_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;