package models

import (
	"time"
)

// WeatherData 存储气象数据的结构体
type WeatherData struct {
	Station     string    `gorm:"index;not null" json:"station"`                       // 添加索引加速查询
	Time        string    `gorm:"column:record_time;autoUpdateTime:milli" json:"time"` // 自定义字段名和更新规则
	Temperature float64   `json:"temperature"`
	Humidity    int       `json:"humidity"`
	WindDir     int       `json:"wind_dir"`
	WindSpeed   float64   `json:"wind_speed"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}
