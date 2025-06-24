package env

import (
	"github.com/spf13/viper"
	"github.com/stxkxs/ok-cli/logger"
	"os"
	"path/filepath"
)

const self = "~/.ok"

func Decode[T any](input, or string) (T, error) {
	var result T

	viper.Reset()

	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	paths := []string{
		".",
		expandPath(self),
		"/etc",
	}

	configName := input
	if configName == "" {
		configName = or
	}

	if input != "" {
		if _, err := os.Stat(input); err == nil {
			viper.SetConfigFile(input)
		} else {
			viper.SetConfigName(configName)
		}
	} else {
		viper.SetConfigName(configName)
	}

	var lastError error
	var found bool

	for _, path := range paths {
		viper.AddConfigPath(path)
		if err := viper.ReadInConfig(); err == nil {
			logger.Logger.Info().
				Str("conf", path).
				Msg("configuration found")
			found = true
			break
		} else {
			lastError = err
			logger.Logger.Warn().
				Err(err).
				Str("conf", path).
				Msg("read conf failed")
		}
	}

	if !found {
		logger.Logger.Error().
			Err(lastError).
			Msg("unable to read configuration from any path")
		return result, lastError
	}

	if err := viper.Unmarshal(&result); err != nil {
		logger.Logger.Error().
			Err(err).
			Msg("unable to decode configuration")
		return result, err
	}

	logger.Logger.Debug().
		Interface("decoded", result).
		Msg("decoded")

	return result, nil
}

func expandPath(path string) string {
	expandedPath := path
	if homeDir, err := os.UserHomeDir(); err == nil {
		expandedPath = filepath.Join(homeDir, path[2:])
	}
	return expandedPath
}
