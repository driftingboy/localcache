package main

import (
	"fmt"

	"github.com/spf13/viper"
)

var GlobalConfig *Config

// var globalMutex sync.RWMutex

type Config struct {
	// TODO
}

func InitConfig(confPath string) error {
	var (
		err       error
		confViper *viper.Viper
	)

	if confViper, err = initViper(confPath); err != nil {
		return fmt.Errorf("Load sdk config failed, %s", err)
	}

	GlobalConfig = &Config{}
	if err = confViper.Unmarshal(&GlobalConfig); err != nil {
		return fmt.Errorf("Unmarshal config file failed, %s", err)
	}

	return nil
}

func initViper(confPath string) (*viper.Viper, error) {
	cmViper := viper.New()
	cmViper.SetConfigFile(confPath)
	if err := cmViper.ReadInConfig(); err != nil {
		return nil, err
	}

	return cmViper, nil
}
