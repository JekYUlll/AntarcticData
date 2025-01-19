package models

// WeatherData 存储气象数据的结构体
type WeatherData struct {
	Station     string  `json:"station"`     // 科考站名称
	Time        string  `json:"time"`        // 时间
	Temperature float64 `json:"temperature"` // 温度
	Humidity    int     `json:"humidity"`    // 湿度
	WindDir     int     `json:"wind_dir"`    // 风向
	WindSpeed   float64 `json:"wind_speed"`  // 风速
}
