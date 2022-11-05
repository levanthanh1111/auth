package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/tpp/msf/shared/log"
)

const (
	configCommonFile = "config"
)

func init() {
	viper.SetConfigName(configCommonFile)
	viper.SetConfigType("yaml")
	// viper.AddConfigPath("$HOME/config")
	viper.AddConfigPath(".")
	viper.MergeInConfig()

	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Logger.Warn().Msg("Config file changed: " + e.Name)
	})
	viper.WatchConfig()
	viper.AutomaticEnv()
	viper.SetEnvPrefix("CONFIG")
}

func Load(c any) {
	if err := viper.Unmarshal(&c, func(dc *mapstructure.DecoderConfig) {
		dc.TagName = "config"
	}); err != nil {
		panic(err)
	}

	log.Logger.Info().Interface("config", c).Msg("")
}

func GetConfig[V int64 | string | float64](key string) V {
	var t V
	var ret any
	switch any(t).(type) {
	case string:
		ret = viper.GetString(key)
	case int64:
		ret = viper.GetInt64(key)
	case float64:
		ret = viper.GetFloat64(key)
	}

	return ret.(V)
}

// GetConfigByte get []byte from config file.
func GetConfigByte(key string) []byte {
	return []byte(viper.GetString(key))
}
