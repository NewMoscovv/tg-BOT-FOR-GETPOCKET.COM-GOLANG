package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"os"
)

type Config struct {
	TelegramToken     string
	PocketConsumerKey string
	AuthServerURL     string
	TelegramBotURL    string `mapstructure:"bot_url"`
	DBPath            string `mapstructure:"db_file"`

	Messages Messages
}
type Messages struct {
	Errors    Errors
	Responses Responses
}

type Errors struct {
	Default      string `mapstructure:"default"`
	InvalidURL   string `mapstructure:"invalid_url"`
	Unauthorized string `mapstructure:"unauthorized"`
	UnableToSave string `mapstructure:"unable_to_save"`
}

type Responses struct {
	Start                 string `mapstructure:"start"`
	AlreadyAuthorized     string `mapstructure:"already_authorized"`
	LinkSavedSuccessfully string `mapstructure:"link_saved_successfully"`
	UnknownCommand        string `mapstructure:"unknown_command"`
}

func Init() (*Config, error) {
	viper.AddConfigPath("configs")
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	if err := viper.UnmarshalKey("messages.responses", &config.Messages.Responses); err != nil {
		return nil, err
	}
	if err := viper.UnmarshalKey("messages.errors", &config.Messages.Errors); err != nil {
		return nil, err
	}
	if err := parseEnv(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

func parseEnv(cfg *Config) error {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}
	if value, exists := os.LookupEnv("TOKEN"); exists {
		cfg.TelegramToken = value
	}
	if value, exists := os.LookupEnv("POCKET_CONSUMER_KEY"); exists {
		cfg.PocketConsumerKey = value
	}
	if value, exists := os.LookupEnv("AUTH_SERVER_URL"); exists {
		cfg.AuthServerURL = value
	}
	return nil

}
