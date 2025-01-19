package cache

import (
	"antarctic/models"
	"sync"
)

// WeatherCache 天气数据缓存
type WeatherCache struct {
	data  map[string]models.WeatherData // 使用站点名称作为key
	mutex sync.RWMutex
}

// New 创建新的缓存实例
func New() *WeatherCache {
	return &WeatherCache{
		data: make(map[string]models.WeatherData),
	}
}

// Update 更新缓存，返回发生变化的数据
func (c *WeatherCache) Update(newData []models.WeatherData) []models.WeatherData {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var changedData []models.WeatherData

	for _, data := range newData {
		// 检查是否存在缓存数据
		if cached, exists := c.data[data.Station]; !exists {
			// 新数据
			changedData = append(changedData, data)
			c.data[data.Station] = data
		} else {
			// 比较数据是否发生变化
			if !isEqual(cached, data) {
				changedData = append(changedData, data)
				c.data[data.Station] = data
			}
		}
	}

	return changedData
}

// isEqual 比较两个WeatherData是否相等
func isEqual(a, b models.WeatherData) bool {
	return a.Time == b.Time &&
		a.Temperature == b.Temperature &&
		a.Humidity == b.Humidity &&
		a.WindDir == b.WindDir &&
		a.WindSpeed == b.WindSpeed
}
