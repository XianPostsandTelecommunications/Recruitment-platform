package config

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	// DB 全局数据库实例
	DB *gorm.DB
)

// InitDatabase 初始化数据库连接
func InitDatabase(config *DatabaseConfig) error {
	dsn := config.GetDSN()
	
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: getGormLogger(config.LogLevel),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})
	
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	// 获取底层的 sql.DB 对象
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("数据库连接测试失败: %w", err)
	}

	DB = db
	return nil
}

// getGormLogger 获取GORM日志配置
func getGormLogger(level string) logger.Interface {
	var logLevel logger.LogLevel
	
	switch level {
	case "silent":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "info":
		logLevel = logger.Info
	default:
		logLevel = logger.Info
	}

	return logger.Default.LogMode(logLevel)
}

// CloseDatabase 关闭数据库连接
func CloseDatabase() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return fmt.Errorf("获取数据库实例失败: %w", err)
		}
		return sqlDB.Close()
	}
	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}

// AutoMigrate 自动迁移数据库表
func AutoMigrate(models ...interface{}) error {
	if DB == nil {
		return fmt.Errorf("数据库未初始化")
	}
	
	return DB.AutoMigrate(models...)
} 