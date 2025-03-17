package models

import (
	"time"

	"gorm.io/gorm"
)

// WeatherData 存储气象数据的结构体
type WeatherData struct {
	ID          uint      `gorm:"primaryKey" json:"-"`
	Station     string    `gorm:"index;not null" json:"station"`
	Time        time.Time `gorm:"column:record_time;index;not null" json:"time"`
	Temperature float64   `json:"temperature"`
	Humidity    int       `json:"humidity"`
	WindDir     int       `json:"wind_dir"`
	WindSpeed   float64   `json:"wind_speed"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`

	// GORM 唯一索引标记
	UniqueID string `gorm:"uniqueIndex:idx_station_time;type:varchar(100);default:''" json:"-"`
}

// BeforeSave GORM钩子，设置联合唯一索引
func (w *WeatherData) BeforeSave(tx *gorm.DB) error {
	// 生成唯一ID: 站点名+时间戳
	w.UniqueID = w.Station + "_" + w.Time.Format("20060102150405")
	return nil
}
