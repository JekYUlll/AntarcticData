package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

const (
	baseURL = "https://www.pric.org.cn"
)

// WeatherData 存储气象数据的结构体
type WeatherData struct {
	Station     string  // 科考站名称
	Time        string  // 时间
	Temperature float64 // 温度
	Humidity    int     // 湿度
	WindDir     int     // 风向
	WindSpeed   float64 // 风速
}

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("www.pric.org.cn"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	fmt.Println("开始获取各科考站气象数据：")

	// 处理HTML内容
	c.OnHTML("body", func(e *colly.HTMLElement) {
		// 遍历每个科考站
		e.ForEach(".sssj-rg", func(_ int, el *colly.HTMLElement) {
			data := WeatherData{}

			// 获取科考站名称
			stationClass := el.Attr("class")
			// 直接从class名称判断站点
			switch {
			case strings.Contains(stationClass, "ccz"):
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

			// 打印数据
			fmt.Printf("\n科考站: %s\n", data.Station)
			fmt.Printf("时间: %s\n", data.Time)
			fmt.Printf("温度: %.1f°C\n", data.Temperature)
			fmt.Printf("湿度: %d%%\n", data.Humidity)
			fmt.Printf("风向: %d°\n", data.WindDir)
			fmt.Printf("风速: %.1fm/s\n", data.WindSpeed)
		})
	})

	// 错误处理
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("请求 %v 失败: %v", r.Request.URL, err)
	})

	// 访问网站并定时更新
	for {
		err := c.Visit(baseURL)
		if err != nil {
			log.Printf("访问失败: %v", err)
		}

		// 每30分钟更新一次数据
		time.Sleep(30 * time.Minute)
	}
}
