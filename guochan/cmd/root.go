package cmd

import (
	"log"

	"practisego/guochan/国产"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	Secret []byte
	M      国产.Method `yaml:method`
	Port   int
	Server string
}

var cfg Config

var cfgFile string

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(clientCmd, serverCmd)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")

	rootCmd.PersistentFlags().IntVarP(&cfg.Port, "port", "p", 0, "listen port")
	flagm := rootCmd.PersistentFlags().UintP("method", "m", 1, "crypto method")
	cfg.M = 国产.Method(*flagm)

	clientCmd.Flags().StringVarP(&cfg.Server, "server", "s", "127.0.0.1:12345", "server address")

	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("server", rootCmd.PersistentFlags().Lookup("server"))
	viper.BindPFlag("method", rootCmd.PersistentFlags().Lookup("method"))

}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	viper.AddConfigPath("/etc/guochan")
	viper.AddConfigPath("usr/local/etc/guochan")
	viper.AddConfigPath("$HOME/.guochan")
	viper.AddConfigPath(".")

	viper.SetConfigName("Guochan")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Println("ReadInConfig Error:", err)
	}

}

var rootCmd = &cobra.Command{Use: "国产"}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
