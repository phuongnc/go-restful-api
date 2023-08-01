package util

import (
	"os"
	"reflect"
	"strings"

	"github.com/google/uuid"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

const (
	// CommonFile is common config file prefix.
	CommonFile = "app"
)

// repository: https://github.com/spf13/viper
func init() {
	// path to look for the config file in
	viper.AddConfigPath(os.Getenv("API_DIR"))
	viper.AddConfigPath(".")
	viper.SetConfigName(CommonFile)

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
}

// Set set value to config file.
func Set(key string, value interface{}) {
	viper.Set(key, value)
}

// Get get valur from config file
func Get(key string) interface{} {
	return viper.Get(key)
}

// GetString get string from config file.
func GetString(key string) string {
	return viper.GetString(key)
}

// GetInt get int from config file.
func GetInt(key string) int {
	return viper.GetInt(key)
}

// GetInt64 get int64 from config file.
func GetInt64(key string) int64 {
	return viper.GetInt64(key)
}

// GetBool get bool from config file.
func GetBool(key string) bool {
	return viper.GetBool(key)
}

// GetStringMap get bool from config file.
func GetStringMap(key string) interface{} {
	return viper.GetStringMap(key)
}

// GetBytes get []byte from config file.
func GetBytes(key string) []byte {
	return []byte(viper.GetString(key))
}

// UnmarshalKey unmarshal config from config file by key.
func UnmarshalKey(key string, i interface{}) error {
	return viper.UnmarshalKey(key, i, func(config *mapstructure.DecoderConfig) {
		config.TagName = "config"
	})
}

// Unmarshal config
func Unmarshal(i interface{}) error {
	return viper.GetViper().Unmarshal(i, func(config *mapstructure.DecoderConfig) {
		config.TagName = "config"
		// custom decode hook for handle uuid string from configuration and parse to github.com/google/uuid.UUID type
		config.DecodeHook = func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
			if t == reflect.TypeOf(uuid.UUID{}) {
				if f.Kind() == reflect.String {
					var uid uuid.UUID
					val, ok := data.(string)
					if ok {
						uid, _ = uuid.Parse(val)
						return uid, nil
					}
				}
			}
			return data, nil
		}
	})
}

// LoadConfig func
func LoadConfig(i interface{}) error {
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		return err
	}
	return Unmarshal(i)
}
