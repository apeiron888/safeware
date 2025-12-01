package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Email    EmailConfig
	Server   ServerConfig
	Audit    AuditConfig
	Backup   BackupConfig
	Captcha  CaptchaConfig
}

type DatabaseConfig struct {
	URI      string
	Database string
}

type RedisConfig struct {
	URL string
}

type JWTConfig struct {
	Secret        string
	RefreshSecret string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

type EmailConfig struct {
	SMTPHost    string
	SMTPPort    int
	SMTPUser    string
	SMTPPass    string
	SMTPFrom    string
	FrontendURL string
}

type ServerConfig struct {
	Port           string
	GinMode        string
	TrustedProxies []string
}

type AuditConfig struct {
	EncryptionKey string
}

type BackupConfig struct {
	Path          string
	RetentionDays int
	EncryptionKey string
}

type CaptchaConfig struct {
	Secret string
}

func Load() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Error reading config file: %v. Using environment variables.", err)
	}

	accessExpiry, _ := time.ParseDuration(viper.GetString("JWT_ACCESS_EXPIRY"))
	if accessExpiry == 0 {
		accessExpiry = 15 * time.Minute
	}

	refreshExpiry, _ := time.ParseDuration(viper.GetString("JWT_REFRESH_EXPIRY"))
	if refreshExpiry == 0 {
		refreshExpiry = 168 * time.Hour // 7 days
	}

	return &Config{
		Database: DatabaseConfig{
			URI:      viper.GetString("MONGODB_URI"),
			Database: viper.GetString("MONGODB_DATABASE"),
		},
		Redis: RedisConfig{
			URL: viper.GetString("REDIS_URL"),
		},
		JWT: JWTConfig{
			Secret:        viper.GetString("JWT_SECRET"),
			RefreshSecret: viper.GetString("JWT_REFRESH_SECRET"),
			AccessExpiry:  accessExpiry,
			RefreshExpiry: refreshExpiry,
		},
		Email: EmailConfig{
			SMTPHost:    viper.GetString("SMTP_HOST"),
			SMTPPort:    viper.GetInt("SMTP_PORT"),
			SMTPUser:    viper.GetString("SMTP_USER"),
			SMTPPass:    viper.GetString("SMTP_PASS"),
			SMTPFrom:    viper.GetString("SMTP_FROM"),
			FrontendURL: viper.GetString("FRONTEND_URL"),
		},
		Server: ServerConfig{
			Port:           viper.GetString("PORT"),
			GinMode:        viper.GetString("GIN_MODE"),
			TrustedProxies: viper.GetStringSlice("TRUSTED_PROXIES"),
		},
		Audit: AuditConfig{
			EncryptionKey: viper.GetString("AUDIT_ENCRYPTION_KEY"),
		},
		Backup: BackupConfig{
			Path:          viper.GetString("BACKUP_PATH"),
			RetentionDays: viper.GetInt("BACKUP_RETENTION_DAYS"),
			EncryptionKey: viper.GetString("BACKUP_ENCRYPTION_KEY"),
		},
		Captcha: CaptchaConfig{
			Secret: viper.GetString("CAPTCHA_SECRET"),
		},
	}
}
