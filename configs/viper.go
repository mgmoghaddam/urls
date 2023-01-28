package configs

import (
	"github.com/spf13/viper"
)

var vi = viper.New()

func InitViper() {
	vi.SetConfigFile("config.yml")
	err := vi.ReadInConfig()
	if err != nil {
		return
	}
}
