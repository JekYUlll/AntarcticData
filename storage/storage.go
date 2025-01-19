package storage

import (
	"antarctic/models"
	"time"
)

// Storage 定义数据存储接口
type Storage interface {
	// Save 保存天气数据
	Save(data []models.WeatherData) error
	// GetLatest 获取最新的天气数据
	GetLatest(station string) (*models.WeatherData, error)
	// GetRange 获取指定时间范围的数据
	GetRange(station string, start, end time.Time) ([]models.WeatherData, error)
	// Close 关闭存储连接
	Close() error
}
