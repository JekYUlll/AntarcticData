CREATE DATABASE IF NOT EXISTS antarctic_data;
USE antarctic_data;

CREATE TABLE IF NOT EXISTS `weather_data` (
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
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;