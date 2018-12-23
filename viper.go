package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type server struct {
	Host string
	Port int
}

type servers []server

var s servers
var cfgFile string

func main() {
	viper.SetConfigName("viper")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Println(err)
	}

	fmt.Println("ConfigFileUsed:", viper.ConfigFileUsed())

	fmt.Println(viper.AllSettings())

	if err := viper.Unmarshal(&s); err != nil {
		log.Fatal(err)
	}

	for k, v := range s {
		fmt.Println(k, v)
	}

}
