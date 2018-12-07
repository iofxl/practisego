package main

import (
	"fmt"

	"github.com/spf13/viper"
)

func main() {

	//viper.SetDefault("Name", "Adel")
	viper.SetConfigName("foo")
	viper.SetConfigType("yaml")

	viper.AddConfigPath("/etc/foo")
	viper.AddConfigPath("/usr/local/etc/foo")
	viper.AddConfigPath("$HOME/.foo/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

}
