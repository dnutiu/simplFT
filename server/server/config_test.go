package server

import (
	"testing"

	"github.com/spf13/viper"
)

func TestLoadConfigFromFile(t *testing.T) {
	// SetDefaultConfiguration must be called BEFORE LoadConfigFromFile.
	InitializedConfiguration()

	Address := viper.GetString("Address")
	if Address == "" {
		t.Error("TestLoadConfigFromFile: Can't get Address!")
	}

	Port := viper.GetString("Port")
	if Port == "" {
		t.Error("TestLoadConfigFromFile: Can't get Port!")
	}

	ConfigPath := viper.GetString("ConfigPath")
	if ConfigPath == "" {
		t.Error("TestLoadConfigFromFile: Can't get ConfigPath!")
	}

	DirDepth := viper.GetInt("MaxDirDepth")
	if DirDepth == 0 {
		t.Error("TestLoadConfigFromFile: Can't get DirDepth!")
	}

	BasePath := viper.GetString("AbsoluteServePath")
	if BasePath == "" {
		t.Error("TestLoadConfigFromFile: Can't get AbsoluteServePath!")
	}
}
