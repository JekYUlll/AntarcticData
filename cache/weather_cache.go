package cache

import (
	"antarctic/models"
	"antarctic/storage"
	"errors"
	"sync"
	"time"

	"gorm.io/gorm"
)

// WeatherCache 天气数据缓存 - 简化版本，只存储每个站点的最新record_time
type WeatherCache struct {
	latestRecordTime map[string]time.Time // 保存每个站点的最新record_time
	mutex            sync.RWMutex
}

// New 创建新的缓存实例
func New() *WeatherCache {
	return &WeatherCache{
		latestRecordTime: make(map[string]time.Time),
	}
}

// InitFromDB 从数据库初始化缓存
func (c *WeatherCache) InitFromDB(storage storage.Storage) error {
	stations := []string{"长城站", "中山站", "昆仑站", "黄河站", "泰山站", "秦岭站"}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, station := range stations {
		data, err := storage.GetLatest(station)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if data != nil {
			c.latestRecordTime[station] = data.Time // data.Time 对应数据库的 record_time
		}
	}
	return nil
}

// IsNewer 检查数据是否比缓存中的更新（基于record_time比较）
func (c *WeatherCache) IsNewer(data models.WeatherData) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	lastRecordTime, exists := c.latestRecordTime[data.Station]
	if !exists {
		return true // 没有缓存，视为新数据
	}

	// 比较record_time，如果新数据的record_time更晚，则认为是新数据
	return data.Time.After(lastRecordTime)
}

// UpdateLatestRecordTime 更新站点的最新record_time
func (c *WeatherCache) UpdateLatestRecordTime(station string, recordTime time.Time) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.latestRecordTime[station] = recordTime
}

// GetLatestRecordTime 获取站点的最新record_time
func (c *WeatherCache) GetLatestRecordTime(station string) (time.Time, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	recordTime, exists := c.latestRecordTime[station]
	return recordTime, exists
}
