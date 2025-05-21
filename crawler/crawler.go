package crawler

import (
	"antarctic/cache"
	"antarctic/models"
	"antarctic/storage"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"gorm.io/gorm"
)

// Crawler 爬虫结构体
type Crawler struct {
	collector *colly.Collector
	handler   func([]models.WeatherData)
	cache     *cache.WeatherCache
}

// New 创建新的爬虫实例
func New(handler func([]models.WeatherData)) *Crawler {
	c := colly.NewCollector(
		// TODO 此处实际上直接写死链接了，规范的话应该给crawler的New函数里传链接+间隔，但一次性，没必要
		colly.AllowedDomains("www.pric.org.cn"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
		colly.AllowURLRevisit(), // 允许重复访问URL
	)

	return &Crawler{
		collector: c,
		handler:   handler,
		cache:     cache.New(),
	}
}

// Start 开始爬取数据
func (c *Crawler) Start(baseURL string) error {
	log.Printf("开始爬取数据: %s", time.Now().Format("2006-01-02 15:04:05"))

	var changedData []models.WeatherData

	// 处理HTML内容
	c.collector.OnHTML("body", func(e *colly.HTMLElement) {
		// 遍历每个科考站
		e.ForEach(".sssj-rg", func(_ int, el *colly.HTMLElement) {
			data := c.parseWeatherData(el)

			// 检查是否为新数据
			if c.cache.IsNewer(data) {
				changedData = append(changedData, data)
				// 更新缓存时间戳
				c.cache.UpdateTimestamp(data)
			}
		})

		// 如果有新数据，调用处理函数
		if len(changedData) > 0 {
			c.handler(changedData)
		}
	})

	// 访问网站并记录结果
	err := c.collector.Visit(baseURL)
	if len(changedData) == 0 {
		log.Printf("爬取完成，没有新数据")
	}
	return err
}

// parseWeatherData 解析天气数据
func (c *Crawler) parseWeatherData(el *colly.HTMLElement) models.WeatherData {
	data := models.WeatherData{
		CreatedAt: time.Now(), // 手动设置创建时间
	}

	// 获取科考站名称
	stationClass := el.Attr("class")
	switch {
	case strings.Contains(stationClass, "czc"):
		data.Station = "长城站"
	case strings.Contains(stationClass, "zsz"):
		data.Station = "中山站"
	case strings.Contains(stationClass, "klz"):
		data.Station = "昆仑站"
	case strings.Contains(stationClass, "hhz"):
		data.Station = "黄河站"
	case strings.Contains(stationClass, "tsz"):
		data.Station = "泰山站"
	case strings.Contains(stationClass, "qlz"):
		data.Station = "秦岭站"
	default:
		data.Station = "未知站点"
	}

	// 获取时间
	timeStr := strings.TrimSpace(strings.TrimPrefix(el.ChildText(".sssj-time span:last-child"), "时间："))
	// 创建中国时区
	chinaLoc, _ := time.LoadLocation("Asia/Shanghai")
	// 在解析时明确指定使用中国时区
	parsedTime, err := time.ParseInLocation("2006-01-02 15:04:05", timeStr, chinaLoc)
	if err != nil {
		log.Printf("时间解析错误: %v, 原始字符串: %s", err, timeStr)
		// 使用当前时间作为后备
		parsedTime = time.Now().In(chinaLoc)
	}
	data.Time = parsedTime

	// 获取气象数据
	el.ForEach(".ssj-wd-rg-list .ssj-wd-rg-item", func(_ int, item *colly.HTMLElement) {
		value := strings.TrimSpace(item.ChildText("span:last-child"))
		switch {
		case strings.Contains(item.Text, "温度"):
			fmt.Sscanf(value, "%f", &data.Temperature)
		case strings.Contains(item.Text, "湿度"):
			fmt.Sscanf(value, "%d", &data.Humidity)
		case strings.Contains(item.Text, "风向"):
			fmt.Sscanf(value, "%d", &data.WindDir)
		case strings.Contains(item.Text, "风速"):
			fmt.Sscanf(value, "%f", &data.WindSpeed)
		}
	})

	return data
}

// InitCacheFromDB 从数据库初始化缓存
func (c *Crawler) InitCacheFromDB(storage storage.Storage) error {
	return c.cache.InitFromDB(storage)
}

// SyncCacheWithDB 同步缓存与数据库
func (c *Crawler) SyncCacheWithDB(storage storage.Storage) error {
	stations := []string{"长城站", "中山站", "昆仑站", "黄河站", "泰山站", "秦岭站"}

	for _, station := range stations {
		// 获取数据库中最新的数据
		data, err := storage.GetLatest(station)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 忽略未找到的记录
				continue
			}
			return err
		}

		if data == nil {
			continue
		}

		// 更新缓存
		c.cache.UpdateTimestamp(*data)
	}

	return nil
}
