package config

import (
	"github.com/spf13/viper"
)

func InitConfig(filePath, fileName string) (*viper.Viper, error) {
	vp := viper.New()
	vp.SetConfigName(fileName)
	vp.SetConfigType("yaml")
	vp.AddConfigPath(filePath)
	if err := vp.ReadInConfig(); err != nil {
		return nil, err
	}
	return vp, nil
}
