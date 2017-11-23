// This file contains the configuration settings for the server.
package server

import (
	"log"

	"flag"

	"github.com/spf13/viper"
)

// ConfigPath will be used via cmd to set the configuration path for the config file.
var ConfigPath string

// loadConfigFromFile tries to load the configuration file from the disk.
func loadConfigFromFile() error {
	viper.SetConfigName("config")
	viper.AddConfigPath(viper.GetString("ConfigPath"))

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Printf("Fatal error reading config file: %s \n", err)
	}
	return err
}

// setDefaultConfiguration will set the default configuration settings.
func setDefaultConfiguration() {
	viper.SetDefault("Address", "localhost")
	viper.SetDefault("Port", "8080")
	viper.SetDefault("ConfigPath", ConfigPath)
	viper.SetDefault("MaxDirDepth", 30)
	viper.SetDefault("AbsoluteServePath", "./")
}

// InitializedConfiguration initializes the configuration for the application.
func InitializedConfiguration() {
	flag.StringVar(&ConfigPath, "ConfigPath", ".", "Set the location of the config file.")
	flag.Parse()

	setDefaultConfiguration()

	loadConfigFromFile()
	BasePath = viper.GetString("AbsoluteServePath")
}
