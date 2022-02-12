package config

import (
	"github.com/spf13/viper"
)

func Init() (err error) {
	defer func() {
		if ok := recover(); ok != nil {
			err = ok.(error)
		}
	}()
	viper.SetConfigFile("./config.yaml")

	if err = viper.ReadInConfig(); err != nil {
		return
	}
	if err = viper.Unmarshal(Configs); err != nil {
		return
	}
	return
}
