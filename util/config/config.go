package config

import "github.com/spf13/viper"

var (
	localConfig *viper.Viper
)

func init() {
	localConfig = viper.New()

	localConfig.SetConfigName("config")
	localConfig.AddConfigPath("$workdirectory/config/")

	if err := localConfig.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// config not found
			panic("config not found")
		} else {
			panic(err.Error())
		}
	}
}

func Get(key string) string {
	return localConfig.GetString(key)
}
