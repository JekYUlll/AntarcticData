package crawler

import (
	"antarctic/cache"
	"antarctic/models"
	"antarctic/storage"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
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

	var newData []models.WeatherData

	// 处理HTML内容
	c.collector.OnHTML("body", func(e *colly.HTMLElement) {
		// 遍历每个科考站
		e.ForEach(".sssj-rg", func(_ int, el *colly.HTMLElement) {
			data := c.parseWeatherData(el)

			// 检查是否为新数据（基于record_time比较）
			if c.cache.IsNewer(data) {
				newData = append(newData, data)
				log.Printf("发现新数据: %s - %s", data.Station, data.Time.Format("2006-01-02 15:04:05"))
			}
		})

		// 如果有新数据，调用处理函数
		if len(newData) > 0 {
			c.handler(newData)
			// 处理完成后更新缓存
			for _, data := range newData {
				c.cache.UpdateLatestRecordTime(data.Station, data.Time)
			}
			log.Printf("处理了 %d 条新数据", len(newData))
		} else {
			log.Printf("爬取完成，没有新数据")
		}
	})

	// 访问网站并记录结果
	err := c.collector.Visit(baseURL)
	return err
}

// parseWeatherData 解析天气数据
func (c *Crawler) parseWeatherData(el *colly.HTMLElement) models.WeatherData {
	// 设置爬取时间
	crawlTime := time.Now()

	data := models.WeatherData{
		CreatedAt: crawlTime, // 爬取时间
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

	// 获取网站显示的记录时间（这是数据的时间戳，不是爬取时间）
	timeStr := strings.TrimSpace(strings.TrimPrefix(el.ChildText(".sssj-time span:last-child"), "时间："))
	chinaLoc, _ := time.LoadLocation("Asia/Shanghai")
	parsedTime, err := time.ParseInLocation("2006-01-02 15:04:05", timeStr, chinaLoc)
	if err != nil {
		log.Printf("时间解析错误: %v, 原始字符串: %s", err, timeStr)
		parsedTime = time.Now().In(chinaLoc)
	}
	data.Time = parsedTime // 这是网站显示的记录时间，用于数据去重和比较

	// 获取气象数据
	el.ForEach(".ssj-wd-rg-list .ssj-wd-rg-item", func(_ int, item *colly.HTMLElement) {
		value := strings.TrimSpace(item.ChildText("span:last-child"))
		switch {
		case strings.Contains(item.Text, "温度"):
			fmt.Sscanf(value, "%f", &data.Temperature)
		case strings.Contains(item.Text, "湿度"):
			// 去掉 % 符号再解析
			value = strings.TrimSuffix(value, "%")
			fmt.Sscanf(value, "%d", &data.Humidity)
		case strings.Contains(item.Text, "风向"):
			// 获取所有 span，尝试解析包含数字的 span
			spans := item.ChildTexts("span")
			for _, s := range spans {
				s = strings.TrimSpace(s)
				s = strings.TrimSuffix(s, "°") // 去掉单位
				if n, _ := fmt.Sscanf(s, "%d", &data.WindDir); n == 1 {
					break
				}
			}
		case strings.Contains(item.Text, "风速"):
			// 去掉 m/s 再解析
			value = strings.TrimSuffix(value, "m/s")
			fmt.Sscanf(value, "%f", &data.WindSpeed)
		}
	})

	return data
}

// InitCacheFromDB 从数据库初始化缓存
func (c *Crawler) InitCacheFromDB(storage storage.Storage) error {
	return c.cache.InitFromDB(storage)
}
