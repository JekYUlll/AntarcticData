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
	refreshMinutes = 5 // 刷新间隔（分钟）
)

func main() {

	dsn := "root:ZZYzzy4771430///@tcp(127.0.0.1:3306)/antarctic_data?charset=utf8mb4&parseTime=True&loc=Local"
	ms, err := storage.NewMysqlStorage(dsn)
	if err != nil {
		panic(err.Error())
	}

	// 创建组合处理器
	h := handler.NewMultiHandler(
		[]handler.WeatherHandler{
			handler.NewJSONHandler(),
			handler.NewDBHandler(ms),
		})

	// 创建爬虫实例
	c := crawler.New(h.Handle)

	// 初始化缓存
	if err := c.InitCacheFromDB(ms); err != nil {
		log.Printf("初始化缓存失败: %v", err)
	}

	// 定时任务
	ticker := time.NewTicker(refreshMinutes * time.Minute)
	defer ticker.Stop()

	// 缓存同步任务（每小时执行一次）
	syncTicker := time.NewTicker(1 * time.Hour)
	defer syncTicker.Stop()

	// 首次运行
	if err := c.Start(baseURL); err != nil {
		log.Printf("访问失败: %v", err)
	}

	// 定时任务处理
	for {
		select {
		case <-ticker.C:
			if err := c.Start(baseURL); err != nil {
				log.Printf("访问失败: %v", err)
			}
		case <-syncTicker.C:
			// 重新初始化缓存以确保与数据库同步
			if err := c.InitCacheFromDB(ms); err != nil {
				log.Printf("重新初始化缓存失败: %v", err)
			} else {
				log.Printf("缓存已重新同步")
			}
		}
	}
}
