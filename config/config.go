package config

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Viper global configuration
var Viper = viper.New()

func arrayValue(s string) []string {
	return strings.Split(s, ",")
}

func init() {
	Viper.SetDefault("metal_log_level", "error")
	Viper.SetDefault("metal_log_formatter", "json")
	Viper.SetDefault("metal_log_destination", "stderr")
	Viper.AutomaticEnv()
	logLevel := Viper.GetString("metal_log_level")
	level, _ := log.ParseLevel(logLevel)
	log.SetLevel(level)
	output := Viper.GetString("metal_log_destination")
	switch output {
	case "stdout":
		log.SetOutput(os.Stdout)
	case "stderr":
		log.SetOutput(os.Stderr)
	default:
		f, err := os.OpenFile(output, os.O_CREATE|os.O_APPEND, 0755)
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(f)
	}
	switch Viper.GetString("metal_log_formatter") {
	case "json":
		// TODO: Add options to customize this
		log.SetFormatter(&log.JSONFormatter{})
	case "text":
		// TODO: Add options to customize this
		log.SetFormatter(&log.TextFormatter{})
	}
}
