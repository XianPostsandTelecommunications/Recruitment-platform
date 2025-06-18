package config

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

// Config 应用配置结构
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Log      LogConfig      `mapstructure:"log"`
	Upload   UploadConfig   `mapstructure:"upload"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         string        `mapstructure:"port"`
	Mode         string        `mapstructure:"mode"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	Database     string `mapstructure:"database"`
	Charset      string `mapstructure:"charset"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	LogLevel     string `mapstructure:"log_level"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret     string        `mapstructure:"secret"`
	ExpireTime time.Duration `mapstructure:"expire_time"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	Output     string `mapstructure:"output"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

// UploadConfig 文件上传配置
type UploadConfig struct {
	Path       string `mapstructure:"path"`
	MaxSize    int64  `mapstructure:"max_size"`
	AllowedExt string `mapstructure:"allowed_ext"`
}

var (
	// GlobalConfig 全局配置实例
	GlobalConfig *Config
)

// LoadConfig 加载配置
func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()

	// 设置默认值
	setDefaults()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("读取配置文件失败: %w", err)
		}
		log.Println("配置文件未找到，使用环境变量和默认值")
	}

	// 绑定环境变量
	bindEnvs()

	// 解析配置
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}

	// 验证配置
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	GlobalConfig = &config
	return &config, nil
}

// setDefaults 设置默认值
func setDefaults() {
	// 服务器默认配置
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.read_timeout", "30s")
	viper.SetDefault("server.write_timeout", "30s")
	viper.SetDefault("server.idle_timeout", "60s")

	// 数据库默认配置
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.charset", "utf8mb4")
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("database.max_open_conns", 100)
	viper.SetDefault("database.log_level", "info")

	// Redis默认配置
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.db", 0)

	// JWT默认配置
	viper.SetDefault("jwt.expire_time", "24h")

	// 日志默认配置
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")
	viper.SetDefault("log.output", "stdout")
	viper.SetDefault("log.max_size", 100)
	viper.SetDefault("log.max_backups", 3)
	viper.SetDefault("log.max_age", 28)
	viper.SetDefault("log.compress", true)

	// 文件上传默认配置
	viper.SetDefault("upload.path", "uploads")
	viper.SetDefault("upload.max_size", 10485760) // 10MB
	viper.SetDefault("upload.allowed_ext", "jpg,jpeg,png,gif,pdf,doc,docx")
}

// bindEnvs 绑定环境变量
func bindEnvs() {
	// 服务器环境变量
	viper.BindEnv("server.port", "SERVER_PORT")
	viper.BindEnv("server.mode", "SERVER_MODE")

	// 数据库环境变量
	viper.BindEnv("database.host", "DB_HOST")
	viper.BindEnv("database.port", "DB_PORT")
	viper.BindEnv("database.username", "DB_USER")
	viper.BindEnv("database.password", "DB_PASSWORD")
	viper.BindEnv("database.database", "DB_NAME")

	// Redis环境变量
	viper.BindEnv("redis.host", "REDIS_HOST")
	viper.BindEnv("redis.port", "REDIS_PORT")
	viper.BindEnv("redis.password", "REDIS_PASSWORD")
	viper.BindEnv("redis.db", "REDIS_DB")

	// JWT环境变量
	viper.BindEnv("jwt.secret", "JWT_SECRET")
	viper.BindEnv("jwt.expire_time", "JWT_EXPIRE_HOURS")

	// 日志环境变量
	viper.BindEnv("log.level", "LOG_LEVEL")
	viper.BindEnv("log.output", "LOG_FILE")

	// 文件上传环境变量
	viper.BindEnv("upload.path", "UPLOAD_PATH")
	viper.BindEnv("upload.max_size", "MAX_FILE_SIZE")
}

// validateConfig 验证配置
func validateConfig(config *Config) error {
	// 验证服务器配置
	if config.Server.Port == "" {
		return fmt.Errorf("服务器端口不能为空")
	}

	// 验证数据库配置
	if config.Database.Host == "" {
		return fmt.Errorf("数据库主机不能为空")
	}
	if config.Database.Username == "" {
		return fmt.Errorf("数据库用户名不能为空")
	}
	if config.Database.Database == "" {
		return fmt.Errorf("数据库名称不能为空")
	}

	// 验证JWT配置
	if config.JWT.Secret == "" {
		return fmt.Errorf("JWT密钥不能为空")
	}

	return nil
}

// GetDSN 获取数据库连接字符串
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		c.Username, c.Password, c.Host, c.Port, c.Database, c.Charset)
}

// GetRedisAddr 获取Redis地址
func (c *RedisConfig) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// GetJWTExpireHours 获取JWT过期时间（小时）
func (c *JWTConfig) GetJWTExpireHours() int {
	return int(c.ExpireTime.Hours())
}

// IsDevelopment 是否为开发环境
func (c *ServerConfig) IsDevelopment() bool {
	return c.Mode == "debug" || c.Mode == "development"
}

// IsProduction 是否为生产环境
func (c *ServerConfig) IsProduction() bool {
	return c.Mode == "release" || c.Mode == "production"
} 