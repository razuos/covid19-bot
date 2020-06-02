package config

import (
	"log"

	"github.com/spf13/viper"
)

// AppConfig holds a global instance of config
var AppConfig Config

// Load loads the config from the file and env vars to the global instance.
func Load() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %s", err.Error())
	}

	viper.SetDefault("Twitter.TweetIntervalMinutes", "15")
	viper.SetDefault("Gatherer.UpdateIntervalMinutes", "30")

	err = viper.Unmarshal(&AppConfig)
	if err != nil {
		log.Fatalf("Error parsing config: %s", err.Error())
	}
}

// Check if a config is defined
func Check(keys []string) bool {
	allSettings := viper.AllSettings()
	for key, val := range allSettings {
		if key == "twitter" {
			v := val.(map[string]interface{})
			for _, key := range keys {
				if v[key] == nil {
					return false
				}
			}
		}
	}
	return true
}
