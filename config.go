package main

import (
	"fmt"

	"github.com/spf13/viper"
)

// Create private data struct to hold config options.
type Config struct {
	Hostname  string `yaml:"hostname"`
	Subdomain string `yaml:"subdomain"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
}

// Create a new config instance.
var (
	conf *Config
)

// Read the config file from the current directory and marshal
// into the conf config struct.
func getConf() *Config {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()

	if err != nil {
		fmt.Printf("%v", err)
	}

	conf := &Config{}
	err = viper.Unmarshal(conf)
	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
	}

	return conf
}

// Initialization routine.
func init() {
	// Retrieve config options.
	conf = getConf()
}
