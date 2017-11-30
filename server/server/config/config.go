// This file contains the configuration settings for the server.
package config

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// loadConfigFromFile tries to load the configuration file from the disk.
func loadConfigFromFile(configName string) error {
	viper.SetConfigName(configName)
	viper.AddConfigPath(viper.GetString("ConfigPath"))

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Printf("Fatal error reading config file: %s \n", err)
	}
	return err
}

// setDefaultConfiguration will set the default configuration settings.
func setDefaultConfiguration(configPath string) {
	viper.SetDefault("address", "localhost")
	viper.SetDefault("port", 8080)
	viper.SetDefault("configPath", configPath)
	viper.SetDefault("maxDirDepth", 30)
	viper.SetDefault("absoluteServePath", "./")
	viper.SetDefault("pic.x", 0)
	viper.SetDefault("pic.y", 0)
	viper.SetDefault("pic.color", false)
	viper.SetDefault("upload.enabled", false)
	viper.SetDefault("upload.directory", "upload")
	viper.SetDefault("upload.timeout", 3)
	viper.SetDefault("upload.address", "localhost")
	viper.SetDefault("upload.port", 8081)
}

// InitializeConfiguration initializes the configuration for the application.
func InitializeConfiguration(configName string, configPath string) {
	setDefaultConfiguration(configPath)
	loadConfigFromFile(configName)

	viper.WatchConfig()
}

func ChangeCallback(cb func(event fsnotify.Event)) {
	viper.OnConfigChange(cb)
}
