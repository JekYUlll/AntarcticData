package storage

import (
	"antarctic/models"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// MysqlStorage MySQL存储实现
type MysqlStorage struct {
	db        *gorm.DB
	tableName string // 表名
	batchSize int    // 批量操作阈值
}

func NewMysqlStorage(dsn string, opts ...Option) (*MysqlStorage, error) {
	// 配置日志
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,  // 慢SQL阈值
			LogLevel:                  logger.Error, // 只记录错误
			IgnoreRecordNotFoundError: true,         // 忽略记录未找到错误
			Colorful:                  true,         // 彩色输出
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true,
		Logger:      newLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("连接失败: %w | DSN: %s", err, dsn)
	}

	// TODO WTF
	// if db.Migrator().HasTable(&models.WeatherData{}) {
	// 	// 删除表和索引
	// 	if err := db.Migrator().DropTable(&models.WeatherData{}); err != nil {
	// 		return nil, fmt.Errorf("删除表失败: %w", err)
	// 	}
	// }

	// // 自动迁移表结构
	// if err := db.AutoMigrate(&models.WeatherData{}); err != nil {
	// 	return nil, fmt.Errorf("自动迁移表结构失败: %w", err)
	// }

	// 默认配置
	storage := &MysqlStorage{
		db:        db,
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

// Save 保存天气数据到数据库
func (m MysqlStorage) Save(data []models.WeatherData) error {
	for _, item := range data {
		// 先检查记录是否已存在
		var exists bool
		err := m.db.Table(m.tableName).
			Where("unique_id = ?", item.Station+"_"+item.Time.Format("20060102150405")).
			Select("1").
			Limit(1).
			Find(&exists).Error

		if err != nil {
			return err
		}

		// 只有在记录不存在时才插入
		if !exists {
			if err := m.db.Create(&item).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *MysqlStorage) GetLatest(station string) (*models.WeatherData, error) {
	var data models.WeatherData
	err := m.db.Table(m.tableName).
		Where("station = ?", station).
		Order("record_time DESC").
		First(&data).Error

	if err != nil {
		// 不将"记录未找到"作为错误返回
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 返回nil,nil而不是错误
		}
		return nil, err
	}
	return &data, nil
}

func (m *MysqlStorage) GetRange(station string, start, end time.Time) ([]models.WeatherData, error) {
	var data []models.WeatherData
	query := m.db.Table(m.tableName).
		Where("record_time BETWEEN ? AND ?", start, end)

	if station != "" {
		query = query.Where("station = ?", station)
	}

	err := query.Order("record_time ASC").Find(&data).Error
	return data, err
}

func (m *MysqlStorage) Close() error {
	sqlDB, _ := m.db.DB()
	return sqlDB.Close()
}
