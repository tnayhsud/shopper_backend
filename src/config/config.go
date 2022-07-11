package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type contextKey uint

const UserKey contextKey = 0

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath("src/config")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(fmt.Errorf("file: %w not found", err))
		} else {
			panic(fmt.Errorf("failed to read : %w", err))
		}
	}
}

func ReadConfig(key string) string {
	value, ok := viper.Get(key).(string)
	// If the type is a string then ok will be true
	if !ok {
		log.Fatalf("Invalid type assertion")
	}
	return value
}
