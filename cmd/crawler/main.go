package main

import (
	"antarctic/crawler"
	"antarctic/handler"
	"log"
	"time"
)

const (
	baseURL        = "https://www.pric.org.cn"
	refreshMinutes = 10 // 刷新间隔（分钟）
)

func main() {
	// 创建数据处理器
	handlers := []handler.WeatherHandler{
		handler.NewJSONHandler(), // JSON输出
	}

	// 创建组合处理器
	h := handler.NewMultiHandler(handlers)

	// 创建爬虫实例
	c := crawler.New(h.Handle)

	// 定时任务
	ticker := time.NewTicker(refreshMinutes * time.Minute)
	defer ticker.Stop()

	// 首次运行
	if err := c.Start(baseURL); err != nil {
		log.Printf("访问失败: %v", err)
	}

	// 定时运行
	for range ticker.C {
		if err := c.Start(baseURL); err != nil {
			log.Printf("访问失败: %v", err)
		}
	}
}
