package config

import (
	"github.com/spf13/viper"
	"path/filepath"
)

func InitConfig(root string) error {
	viper.AddConfigPath(filepath.Join(root, "configs"))
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
