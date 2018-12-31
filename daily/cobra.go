package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{}

var serverCmd = &cobra.Command{
	Use: "server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(listen)
	},
}

var clientCmd = &cobra.Command{
	Use: "client",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(server)
	},
}

var server string
var listen string

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(serverCmd, clientCmd)

	serverCmd.PersistentFlags().StringVarP(&listen, "listen", "l", ":12345", "listen address")
	serverCmd.MarkPersistentFlagRequired("listen")

	clientCmd.Flags().StringVarP(&server, "server", "s", ":12345", "server address")
	clientCmd.MarkFlagRequired("server")

}

func initConfig() {
	fmt.Println("initConfig")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	Execute()
}
