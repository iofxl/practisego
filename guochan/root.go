package main

import (
	"crypto/rand"
	"log"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Config is Config
type Config struct {
	Method uint   `yaml:"method"`
	Listen int    `yaml:"listen"`
	Secret string `yaml:"secret"`
	Server string `yaml:"server"`
}

var cfg Config

var cfgFile string
var serverString string

var logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(clientCmd, serverCmd)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")

	rootCmd.PersistentFlags().IntVarP(&cfg.Listen, "listen", "l", 443, "listen port")
	rootCmd.PersistentFlags().UintVarP(&cfg.Method, "method", "m", 1, "crypto method")

	clientCmd.Flags().StringVarP(&cfg.Server, "server", "s", "127.0.0.1:443", "server address")
	clientCmd.Flags().StringVarP(&serverString, "string", "S", "", "server string")

	viper.BindPFlag("listen", rootCmd.PersistentFlags().Lookup("listen"))
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
		logger.Println("ReadInConfig Error:", err)
	}

	err := viper.Unmarshal(&cfg)
	if err != nil {
		logger.Println("Unmarshal:", err)
	}

	if cfg.Secret == "" {
		b := make([]byte, 128)
		_, _ = rand.Read(b)
		cfg.Secret = string(b)
		viper.Set("secret", cfg.Secret)
	}

	if viper.ConfigFileUsed() == "" {
		logger.Println("Since there is no configuration file, I will create one for you!")
		err := viper.WriteConfigAs("Guochan.yaml")
		if err != nil {
			logger.Println("WriteConfigAs:", err)
		}
		logger.Println("WriteConfigAs Guochan.yaml done")
	} else {
		logger.Println("ConfigFileUsed:", viper.ConfigFileUsed())
	}
}

var rootCmd = &cobra.Command{Use: "guochan"}

var serverCmd = &cobra.Command{
	Use: "server",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Printf("Listen: %v\n", cfg.Listen)

		addr := ":" + strconv.Itoa(cfg.Listen)
		if err := ListenAndServeSS(addr, cfg, logger); err != nil {
			logger.Fatal(err)
		}
	},
}

var clientCmd = &cobra.Command{
	Use: "client",
	Run: func(cmd *cobra.Command, args []string) {

		logger.Printf("Listen: %v\n", cfg.Listen)
		logger.Printf("Server: %v\n", cfg.Server)
		addr := ":" + strconv.Itoa(cfg.Listen)
		if err := ListenAndServe(addr, cfg, logger); err != nil {
			logger.Fatal(err)
		}
	},
}

// Execute is Execute
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Fatal(err)
	}
}
