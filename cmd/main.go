package main

import (
	"antarctic/crawler"
	"antarctic/handler"
	"antarctic/storage"
	"log"
	"time"
)

const (
	baseURL        = "https://www.pric.org.cn"
	refreshMinutes = 10 // 刷新间隔（分钟）
)

func main() {

	dsn := "root:donotpanic@tcp(127.0.0.1:3306)/antarctic_data?charset=utf8mb4&parseTime=True&loc=Local"
	ms, err := storage.NewMysqlStorage(dsn)
	if err != nil {
		panic(err.Error())
	}

	// 创建数据处理器
	handlers := []handler.WeatherHandler{
		handler.NewJSONHandler(), // JSON输出
		// TODO 添加 db 操作的 handler
		handler.NewDBHandler(ms),
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
