package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{}

var clientCmd = &cobra.Command{
	Use: "client",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(server)
	},
}

var serverCmd = &cobra.Command{
	Use: "server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(addr)
	},
}

var addr string
var server string

// func (c *Command) AddCommand(cmds ...*Command)
func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(serverCmd, clientCmd)

	serverCmd.PersistentFlags().StringVarP(&addr, "listen", "l", "", "listen address")
	serverCmd.MarkPersistentFlagRequired("listen")

	clientCmd.Flags().StringVarP(&server, "server", "s", "", "server address")
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
