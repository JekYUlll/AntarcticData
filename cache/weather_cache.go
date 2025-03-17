package cache

import (
	"antarctic/models"
	"antarctic/storage"
	"errors"
	"sync"
	"time"

	"gorm.io/gorm"
)

// WeatherCache 天气数据缓存
type WeatherCache struct {
	data  map[string]time.Time // 保存站点的最新时间戳而不是整个数据
	mutex sync.RWMutex
}

// New 创建新的缓存实例
func New() *WeatherCache {
	return &WeatherCache{
		data: make(map[string]time.Time),
	}
}

// 从数据库初始化缓存
func (c *WeatherCache) InitFromDB(storage storage.Storage) error {
	stations := []string{"长城站", "中山站", "昆仑站", "黄河站", "泰山站", "秦岭站"}

	for _, station := range stations {
		data, err := storage.GetLatest(station)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if data != nil {
			c.data[station] = data.Time
		}
	}
	return nil
}

// 检查数据是否新鲜
func (c *WeatherCache) IsNewer(data models.WeatherData) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	lastTime, exists := c.data[data.Station]
	if !exists {
		return true // 没有缓存，视为新数据
	}

	return data.Time.After(lastTime)
}

// 更新缓存的时间戳
func (c *WeatherCache) UpdateTimestamp(data models.WeatherData) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data[data.Station] = data.Time
}

// UpdateStation 更新单个站点的数据，返回是否发生变化
func (c *WeatherCache) UpdateStation(data models.WeatherData) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// 检查是否存在缓存数据
	if cached, exists := c.data[data.Station]; !exists {
		// 新数据
		c.data[data.Station] = data.Time
		return true
	} else {
		// 比较数据是否发生变化
		if !isEqual(cached, data.Time) {
			c.data[data.Station] = data.Time
			return true
		}
	}

	return false
}

// isEqual 比较两个WeatherData是否相等
func isEqual(a, b time.Time) bool {
	return a.Equal(b)
}

// SyncCacheWithDB 同步缓存与数据库
func (c *WeatherCache) SyncCacheWithDB(storage storage.Storage) error {
	stations := []string{"长城站", "中山站", "昆仑站", "黄河站", "泰山站", "秦岭站"}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, station := range stations {
		data, err := storage.GetLatest(station)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if data != nil {
			c.data[station] = data.Time
		}
	}

	return nil
}
