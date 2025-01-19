# 中国极地研究中心网站爬虫

这是一个用Go语言编写的网络爬虫，用于抓取中国极地研究中心网站(www.pric.org.cn)的内容。

## 功能特点

- 自动抓取网站文章内容
- 每30分钟自动更新一次
- 支持错误处理和重试机制
- 使用友好的User-Agent避免被封禁

## 使用方法

1. 确保已安装Go 1.16或更高版本
2. 克隆本仓库
3. 运行以下命令：

```bash
cd cmd/crawler
go run main.go
```

## 依赖

- github.com/gocolly/colly/v2：网络爬虫框架

