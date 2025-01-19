package crawler

import (
	"antarctic/cache"
	"antarctic/models"
	"fmt"
	"strings"

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
	)

	return &Crawler{
		collector: c,
		handler:   handler,
		cache:     cache.New(),
	}
}

// Start 开始爬取数据
func (c *Crawler) Start(baseURL string) error {
	var weatherDataList []models.WeatherData

	// 处理HTML内容
	c.collector.OnHTML("body", func(e *colly.HTMLElement) {
		// 遍历每个科考站
		e.ForEach(".sssj-rg", func(_ int, el *colly.HTMLElement) {
			data := c.parseWeatherData(el)
			weatherDataList = append(weatherDataList, data)
		})

		// 检查数据变化并调用处理函数
		if changedData := c.cache.Update(weatherDataList); len(changedData) > 0 {
			c.handler(changedData)
		}
	})

	// 访问网站
	return c.collector.Visit(baseURL)
}

// parseWeatherData 解析天气数据
func (c *Crawler) parseWeatherData(el *colly.HTMLElement) models.WeatherData {
	data := models.WeatherData{}

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
	data.Time = strings.TrimSpace(strings.TrimPrefix(el.ChildText(".sssj-time span:last-child"), "时间："))

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
