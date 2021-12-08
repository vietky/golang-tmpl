package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"

	"git.chotot.org/fse/orchestration-svc/libraries/logger"
)

// Config ...
var (
	Values Schema

	log = logger.GetLogger("config")
)

// Schema struct
type Schema struct {
	App struct {
		Port       int   `mapstructure:"port"`
		MaxConns   int32 `mapstructure:"max_conns"`
		Pprof      bool  `mapstructure:"pprof"`
		Prometheus bool  `mapstructure:"app"`
	} `mapstructure:"app"`
}

// InitConfig ...
func init() {
	// Initialize viper default instance with API base config.
	config := viper.New()
	config.SetConfigName("config")        // Name of config file (without extension).
	config.AddConfigPath(".")             // Look for config in current directory
	config.AddConfigPath("./config")      // Optionally look for config in the working directory.
	config.AddConfigPath("../config/")    // Look for config needed for tests.
	config.AddConfigPath("../../config/") // Look for config needed for tests.
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	config.AutomaticEnv()
	// Initialize map that contains viper configuration objects.
	err := config.ReadInConfig() // Find and read the config file
	if err != nil {              // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	err = config.Unmarshal(&Values)
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	log.Infof("Current Config: %+v", Values)
}
