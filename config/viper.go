package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

var v *viper.Viper

func V() *viper.Viper {
	if v == nil {
		err := fmt.Errorf("viper not initialized")
		panic(err)
	}
	return v
}

func Save() error {
	return v.WriteConfig()
}

func updateViper() error {

	_, err := os.Stat(GetConfigDir())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = os.MkdirAll(GetConfigDir(), 0755)
			if err != nil {
				return fmt.Errorf("failed to create config dir: %w", err)
			}
		} else {
			return fmt.Errorf("failed to check for config dir: %w", err)
		}
	}

	v = viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	v.AddConfigPath(GetConfigDir())
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	_, err = os.Stat(GetFilePath("config.yaml"))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = createMissingConfig()
			if err != nil {
				return fmt.Errorf("failed to create config file: %w", err)
			}
		} else {
			return fmt.Errorf("failed to check for config file: %w", err)
		}
	}

	err = v.ReadInConfig()
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	return nil
}

func Viper() *viper.Viper {
	if v == nil {
		err := fmt.Errorf("viper not initialized")
		panic(err)
	}
	return v
}

func createMissingConfig() error {
	err := v.SafeWriteConfig()
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	return nil
}
