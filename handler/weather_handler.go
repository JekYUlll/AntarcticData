package handler

import (
	"antarctic/models"
	"antarctic/storage"
	"fmt"
	"log"
)

// WeatherHandler 定义处理气象数据的接口
type WeatherHandler interface {
	Handle(data []models.WeatherData)
}

// JSONHandler JSON格式处理器
type JSONHandler struct{}

// NewJSONHandler 创建新的JSON处理器
func NewJSONHandler() *JSONHandler {
	return &JSONHandler{}
}

// Handle 实现WeatherHandler接口，将数据格式化输出
func (h *JSONHandler) Handle(data []models.WeatherData) {
	// 展示新获取的数据
	for _, d := range data {
		fmt.Printf("新数据: %s - %s\n", d.Station, d.Time.Format("2006-01-02 15:04:05"))
	}
}

// ConsoleHandler 控制台格式处理器
type ConsoleHandler struct{}

// NewConsoleHandler 创建新的控制台处理器
func NewConsoleHandler() *ConsoleHandler {
	return &ConsoleHandler{}
}

// Handle 实现WeatherHandler接口，将数据打印到控制台
func (h *ConsoleHandler) Handle(data []models.WeatherData) {
	for _, d := range data {
		fmt.Printf("\n科考站: %s\n", d.Station)
		fmt.Printf("时间: %s\n", d.Time)
		fmt.Printf("温度: %.1f°C\n", d.Temperature)
		fmt.Printf("湿度: %d%%\n", d.Humidity)
		fmt.Printf("风向: %d°\n", d.WindDir)
		fmt.Printf("风速: %.1fm/s\n", d.WindSpeed)
	}
}

// DBHandler 数据库处理器
type DBHandler struct {
	storage storage.Storage
}

// NewDBHandler 创建新的数据库处理器
func NewDBHandler(storage storage.Storage) *DBHandler {
	return &DBHandler{
		storage: storage,
	}
}

// Handle 实现WeatherHandler接口，将数据保存到数据库
func (h *DBHandler) Handle(data []models.WeatherData) {
	if err := h.storage.Save(data); err != nil {
		log.Printf("保存数据失败: %v", err)
	}
}

// MultiHandler 组合多个处理器
type MultiHandler struct {
	handlers []WeatherHandler
}

// NewMultiHandler 创建新的组合处理器
func NewMultiHandler(handlers []WeatherHandler) *MultiHandler {
	return &MultiHandler{
		handlers: handlers,
	}
}

// Handle 实现WeatherHandler接口，调用所有处理器
func (h *MultiHandler) Handle(data []models.WeatherData) {
	for _, handler := range h.handlers {
		handler.Handle(data)
	}
}
