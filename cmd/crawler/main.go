package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gocolly/colly/v2"
)

const (
	baseURL = "https://www.pric.org.cn"
)

func main() {
	// 创建爬虫实例
	c := colly.NewCollector(
		colly.AllowedDomains("www.pric.org.cn"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	// 处理HTML内容
	c.OnHTML("article", func(e *colly.HTMLElement) {
		title := e.ChildText("h2")
		date := e.ChildText(".date")
		content := e.ChildText("p")
		
		fmt.Printf("标题: %s\n日期: %s\n内容: %s\n\n", title, date, content)
	})

	// 错误处理
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("请求 %v 失败: %v", r.Request.URL, err)
	})

	// 访问网站
	for {
		err := c.Visit(baseURL)
		if err != nil {
			log.Printf("访问失败: %v", err)
		}

		// 等待30分钟后再次爬取
		time.Sleep(30 * time.Minute)
	}
} 