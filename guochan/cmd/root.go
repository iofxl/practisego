package cmd

import (
	"crypto/rand"
	"log"

	"practisego/guochan/国产"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	Secret string    `yaml: secret`
	M      国产.Method `yaml: method`
	Port   int       `yaml: port`
	Server string    `yaml: server`
}

var cfg Config

var cfgFile string
var serverString string

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(clientCmd, serverCmd)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")

	rootCmd.PersistentFlags().IntVarP(&cfg.Port, "port", "p", 443, "listen port")
	flagm := rootCmd.PersistentFlags().UintP("method", "m", 1, "crypto method")
	cfg.M = 国产.Method(*flagm)

	clientCmd.Flags().StringVarP(&cfg.Server, "server", "s", "127.0.0.1:443", "server address")
	clientCmd.Flags().StringVarP(&serverString, "", "S", "", "server string")

	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("method", rootCmd.PersistentFlags().Lookup("method"))
	viper.BindPFlag("server", clientCmd.Flags().Lookup("server"))

}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {

		viper.AddConfigPath("/etc/guochan")
		viper.AddConfigPath("usr/local/etc/guochan")
		viper.AddConfigPath("$HOME/.guochan")
		viper.AddConfigPath(".")

		viper.SetConfigName("Guochan")
		viper.SetConfigType("yaml")
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Println("ReadInConfig Error:", err)
	}

	err := viper.Unmarshal(&cfg)
	if err != nil {
		log.Println("Unmarshal:", err)
	}

	if cfg.Secret == "" {
		b := make([]byte, 128)
		rand.Read(b)
		cfg.Secret = string(b)
		viper.Set("secret", cfg.Secret)
	}

	if viper.ConfigFileUsed() == "" {
		log.Println("Since there is no configuration file, I will create one for you!")
		err := viper.WriteConfigAs("Guochan.yaml")
		if err != nil {
			log.Println("WriteConfigAs:", err)
		}
		log.Println("WriteConfigAs Guochan.yaml done")
	} else {
		log.Println("ConfigFileUsed:", viper.ConfigFileUsed())
	}
}

var rootCmd = &cobra.Command{Use: "guochan"}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
