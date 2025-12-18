package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config 应用配置结构
type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Redis     RedisConfig     `mapstructure:"redis"`
	JWT       JWTConfig       `mapstructure:"jwt"`
	WeChat    WeChatConfig    `mapstructure:"wechat"`
	SMS       SMSConfig       `mapstructure:"sms"`
	Log       LogConfig       `mapstructure:"log"`
	Business  BusinessConfig  `mapstructure:"business"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         int           `mapstructure:"port"`
	Mode         string        `mapstructure:"mode"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	DBName          string        `mapstructure:"dbname"`
	Charset         string        `mapstructure:"charset"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// DSN 返回数据库连接字符串
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
		c.Charset,
	)
}

// RedisConfig Redis配置
type RedisConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

// Addr 返回Redis地址
func (c *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret             string        `mapstructure:"secret"`
	AccessTokenExpire  time.Duration `mapstructure:"access_token_expire"`
	RefreshTokenExpire time.Duration `mapstructure:"refresh_token_expire"`
}

// WeChatConfig 微信配置
type WeChatConfig struct {
	AppID     string          `mapstructure:"app_id"`
	AppSecret string          `mapstructure:"app_secret"`
	Subscribe WeChatSubscribe `mapstructure:"subscribe"`
}

type WeChatSubscribe struct {
	AppointmentReminderTemplateID string `mapstructure:"appointment_reminder_template_id"`
	AppointmentReminderPage       string `mapstructure:"appointment_reminder_page"`
}

// SMSConfig 短信配置
type SMSConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Provider string `mapstructure:"provider"` // console | disabled | (预留第三方服务商)
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

// BusinessConfig 业务规则配置
type BusinessConfig struct {
	Appointment AppointmentConfig `mapstructure:"appointment"`
	Schedule    ScheduleConfig    `mapstructure:"schedule"`
	Checkin     CheckinConfig     `mapstructure:"checkin"`
}

// AppointmentConfig 预约规则配置
type AppointmentConfig struct {
	AdvanceDays        int `mapstructure:"advance_days"`
	MinAdvanceDays     int `mapstructure:"min_advance_days"`
	DailyLimit         int `mapstructure:"daily_limit"`
	CancelDeadlineDays int `mapstructure:"cancel_deadline_days"`
	MonthlyCancelLimit int `mapstructure:"monthly_cancel_limit"`
}

// ScheduleConfig 排班规则配置
type ScheduleConfig struct {
	MorningStart   string `mapstructure:"morning_start"`
	MorningEnd     string `mapstructure:"morning_end"`
	AfternoonStart string `mapstructure:"afternoon_start"`
	AfternoonEnd   string `mapstructure:"afternoon_end"`
	SlotDuration   int    `mapstructure:"slot_duration"`
	MorningSlots   int    `mapstructure:"morning_slots"`
	AfternoonSlots int    `mapstructure:"afternoon_slots"`
}

// CheckinConfig 签到规则配置
type CheckinConfig struct {
	EarlyMinutes int `mapstructure:"early_minutes"`
	LateMinutes  int `mapstructure:"late_minutes"`
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Enabled           bool `mapstructure:"enabled"`
	RequestsPerSecond int  `mapstructure:"requests_per_second"`
	Burst             int  `mapstructure:"burst"`
}

var globalConfig *Config

// Load 加载配置文件
func Load(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// 设置默认值
	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	globalConfig = &cfg
	return &cfg, nil
}

// Get 获取全局配置
func Get() *Config {
	return globalConfig
}

// setDefaults 设置默认配置值
func setDefaults() {
	// 服务器默认配置
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.read_timeout", "10s")
	viper.SetDefault("server.write_timeout", "10s")

	// 数据库默认配置
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.charset", "utf8mb4")
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("database.max_open_conns", 100)
	viper.SetDefault("database.conn_max_lifetime", "1h")

	// Redis默认配置
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.pool_size", 100)

	// JWT默认配置
	viper.SetDefault("jwt.access_token_expire", "2h")
	viper.SetDefault("jwt.refresh_token_expire", "168h")

	// 短信默认配置（默认关闭）
	viper.SetDefault("sms.enabled", false)
	viper.SetDefault("sms.provider", "disabled")

	// 日志默认配置
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.filename", "logs/app.log")
	viper.SetDefault("log.max_size", 100)
	viper.SetDefault("log.max_backups", 30)
	viper.SetDefault("log.max_age", 7)
	viper.SetDefault("log.compress", true)

	// 业务规则默认配置
	viper.SetDefault("business.appointment.advance_days", 7)
	viper.SetDefault("business.appointment.min_advance_days", 1)
	viper.SetDefault("business.appointment.daily_limit", 10)
	viper.SetDefault("business.appointment.cancel_deadline_days", 1)
	viper.SetDefault("business.appointment.monthly_cancel_limit", 5)

	viper.SetDefault("business.schedule.morning_start", "08:00")
	viper.SetDefault("business.schedule.morning_end", "12:00")
	viper.SetDefault("business.schedule.afternoon_start", "14:00")
	viper.SetDefault("business.schedule.afternoon_end", "17:30")
	viper.SetDefault("business.schedule.slot_duration", 15)
	viper.SetDefault("business.schedule.morning_slots", 16)
	viper.SetDefault("business.schedule.afternoon_slots", 14)

	viper.SetDefault("business.checkin.early_minutes", 30)
	viper.SetDefault("business.checkin.late_minutes", 15)

	// 限流默认配置
	viper.SetDefault("rate_limit.enabled", true)
	viper.SetDefault("rate_limit.requests_per_second", 100)
	viper.SetDefault("rate_limit.burst", 200)
}
