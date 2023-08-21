package configuration

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"io/fs"
	"log"
	"os"
)

type configurationFactory struct {
	viper *viper.Viper
}

var factory = configurationFactory{viper: viper.New()}

func Build(prefix string, result interface{}) {
	c := factory.viper.Sub(prefix)
	if c == nil {
		panic(fmt.Errorf("failed to construct config '%s'", prefix))
	}

	if err := c.Unmarshal(&result, func(dc *mapstructure.DecoderConfig) { dc.ErrorUnset = true }); err != nil {
		panic(fmt.Errorf("failed to construct '%s' config: %w", prefix, err))
	}
}

func BuildSkipErrors(prefix string, result interface{}) {
	c := factory.viper.Sub(prefix)
	if c == nil {
		log.Println(fmt.Sprintf("failed to construct config '%s'", prefix))
		return
	}

	if err := c.Unmarshal(&result); err != nil {
		log.Println(fmt.Errorf("failed to construct '%s' config: %w", prefix, err))
	}
}

// SetKeyValue sets the value for the key in the override register, ONLY use for testing purpose.
// May lead to unexpected dropping of configs.
func SetKeyValue(key string, value interface{}) {
	factory.viper.Set(key, value)
}

func Init(rPaths ...string) {
	factory.viper.AutomaticEnv()

	gp, _ := os.Getwd()

	// Use project dir as source
	factory.viper.SetConfigFile(fmt.Sprintf("%s/configs/config.yaml", gp))
	if err := factory.viper.MergeInConfig(); err != nil {
		if _, ok := err.(*fs.PathError); ok {
			// Config file not found; ignore error if desired
			log.Println("Failed to read config file: config file not found")
		} else {
			// Config file was found but another error was produced
			panic(fmt.Errorf("Fatal error config file: %w \n", err))
		}
	}

	// Use input dir as source
	if len(rPaths) > 0 {
		for _, rPath := range rPaths {
			factory.viper.AddConfigPath(rPath)
		}
		factory.viper.SetConfigName("config")
		factory.viper.SetConfigType("yaml")
		if err := factory.viper.MergeInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				// Config file not found; ignore error if desired
				log.Println("Failed to read config file: config file not found")
			} else {
				// Config file was found but another error was produced
				panic(fmt.Errorf("Fatal error config file: %w \n", err))
			}
		}
	}

	// Use vault dir as source
	factory.viper.SetConfigFile("/vault/secrets/config.yaml")
	if err := factory.viper.MergeInConfig(); err != nil {
		if _, ok := err.(*fs.PathError); ok {
			// Config file not found; ignore error if desired
			log.Println("Failed to read config file: vault config file not found")
		} else {
			// Config file was found but another error was produced
			panic(fmt.Errorf("Fatal error config file: %w \n", err))
		}
	}

}
