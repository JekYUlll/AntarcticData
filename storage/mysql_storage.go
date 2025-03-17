package storage

import (
	"antarctic/models"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// TODO 实现 storage 接口逻辑
type MysqlStorage struct {
	db        *gorm.DB
	tableName string // 表名
	batchSize int    // 批量操作阈值（默认200）
}

func NewMysqlStorage(dsn string, opts ...Option) (*MysqlStorage, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true, // 开启预编译提升性能
	})
	if err != nil {
		return nil, fmt.Errorf("连接失败: %w | DSN: %s", err, dsn)
	}
	// 默认配置
	storage := &MysqlStorage{
		db: db,
		// TODO 更改表名
		tableName: "weather_data",
		batchSize: 200,
	}
	// 应用可选参数
	for _, opt := range opts {
		opt(storage)
	}
	// 配置连接池
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return storage, nil
}

// Option 模式配置参数
type Option func(*MysqlStorage)

func WithTableName(name string) Option {
	return func(s *MysqlStorage) {
		s.tableName = name
	}
}

func WithBatchSize(size int) Option {
	return func(s *MysqlStorage) {
		s.batchSize = size
	}
}

// ---------------

// TODO 检查
func (m MysqlStorage) Save(data []models.WeatherData) error {
	return m.db.Table(m.tableName).CreateInBatches(data, m.batchSize).Error
}

func (m *MysqlStorage) GetLatest(station string) (*models.WeatherData, error) {
	var data models.WeatherData
	err := m.db.Table(m.tableName).
		Where("station = ?", station).
		Order("timestamp DESC").
		First(&data).Error
	return &data, err
}

func (m *MysqlStorage) GetRange(station string, start, end time.Time) ([]models.WeatherData, error) {
	var data []models.WeatherData
	query := m.db.Table(m.tableName).
		Where("timestamp BETWEEN ? AND ?", start, end)

	if station != "" {
		query = query.Where("station = ?", station)
	}

	err := query.Order("timestamp ASC").Find(&data).Error
	return data, err
}

func (m *MysqlStorage) Close() error {
	sqlDB, _ := m.db.DB()
	return sqlDB.Close()
}
