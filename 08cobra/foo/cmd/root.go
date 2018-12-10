package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{Use: "app"}

var printCmd = &cobra.Command{
	Use: "print",
	Run: printfunc,
}

var cfgFile string

type Config struct {
	Port int
}

var cfg Config

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(printCmd)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")

	printCmd.Flags().IntVarP(&cfg.Port, "port", "p", 0, "port")

	viper.BindPFlag("port", printCmd.Flags().Lookup("port"))

}

func initConfig() {

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {

		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/." + rootCmd.Use)
		viper.SetConfigName(rootCmd.Use)
		viper.SetConfigType("yaml")
	}

	err := viper.ReadInConfig()
	if err != nil {
		log.Println(err)
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatal("Unmarshal:", err)
	}

	err = viper.WriteConfigAs("app.yaml")
	if err != nil {
		log.Fatal("WriteConfig:", err)
	}

}

func printfunc(cmd *cobra.Command, args []string) {

	fmt.Println("In printCmd:", cfg.Port)

}

func Execute() {

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
