package config

import (
	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	Environment       string `mapstructure:"ENVIRONMENT"`
	DBDriver          string `mapstructure:"DB_DRIVER"`
	DBSource          string `mapstructure:"DB_SOURCE"`
	HTTPServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`
	MonolithUrl       string `mapstructure:"MONOLITH_URL"`
	ApiKey            string `mapstructure:"API_KEY"`
	BaseUrl           string `mapstructure:"BASE_URL"`
	IpWhitelist       string `mapstructure:"IP_WHITELIST"`
	ServiceKeyPath    string `mapstructure:"SERVICE_KEY_PATH"`

	RabbitmqUrl                  string `mapstructure:"RABBITMQ_URL"`
	RabbitmqTopic                string `mapstructure:"RABBITMQ_TOPIC"`
	GooglePubSubProjectId        string `mapstructure:"GOOGLE_PUBSUB_PROJECT_ID"`
	GooglePubSubTopic            string `mapstructure:"GOOGLE_PUBSUB_TOPIC"`
	GoogleApplicationCredentials string `mapstructure:"GOOGLE_APPLICATION_CREDENTIALS"`

	EmailSMTPHost     string `mapstructure:"EMAIL_SMTP_HOST"`
	EmailSMTPPort     int    `mapstructure:"EMAIL_SMTP_PORT"`
	EmailAUTHUsername string `mapstructure:"EMAIL_AUTH_USERNAME"`
	EmailAUTHPassword string `mapstructure:"EMAIL_AUTH_PASSWORD"`
	EmailSenderName   string `mapstructure:"EMAIL_SENDER_NAME"`

	VonageUsername              string `mapstructure:"VONAGE_USERNAME"`
	VonagePassword              string `mapstructure:"VONAGE_PASSWORD"`
	VonageWhatsAppUrl           string `mapstructure:"VONAGE_WHATSAPP_URL"`
	VonageWhatsAppDefaultSender string `mapstructure:"VONAGE_WHATSAPP_DEFAULT_SENDER"`

	VonageSmsUrl    string `mapstructure:"VONAGE_SMS_URL"`
	VonageSmsSender string `mapstructure:"VONAGE_SMS_SENDER"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
