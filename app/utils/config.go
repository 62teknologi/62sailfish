package utils

import (
	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	Environment       string `mapstructure:"ENVIRONMENT"`
	DBSource          string `mapstructure:"DB_SOURCE"`
	HTTPServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`
	MonolithUrl       string `mapstructure:"MONOLITH_URL"`
	ApiKey            string `mapstructure:"API_KEY"`

	EmailSMTPHost     string `mapstructure:"EMAIL_SMTP_HOST"`
	EmailSMTPPort     int    `mapstructure:"EMAIL_SMTP_PORT"`
	EmailAUTHUsername string `mapstructure:"EMAIL_AUTH_USERNAME"`
	EmailAUTHPassword string `mapstructure:"EMAIL_AUTH_PASSWORD"`
	EmailSenderName   string `mapstructure:"EMAIL_SENDER_NAME"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")

	viper.SetDefault("ENVIRONMENT", "development")
	viper.SetDefault("DB_SOURCE", "root:password@tcp(127.0.0.1:3306)/dolphin?charset=utf8mb4&parseTime=True&loc=Local")
	viper.SetDefault("HTTP_SERVER_ADDRESS", "0.0.0.0:8080")
	viper.SetDefault("MONOLITH_URL", "http://localhost:8000")
	viper.SetDefault("API_KEY", "sadjaskdjlkasjdlkasdj")

	viper.SetDefault("EMAIL_SMTP_HOST", "sandbox.smtp.mailtrap.io")
	viper.SetDefault("EMAIL_SMTP_PORT", "587")
	viper.SetDefault("EMAIL_AUTH_USERNAME", "123455av123231")
	viper.SetDefault("EMAIL_AUTH_PASSWORD", "12312asdas1231")
	viper.SetDefault("EMAIL_SENDER_NAME", "Dolphin <dolphin@62teknologi.com>")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
