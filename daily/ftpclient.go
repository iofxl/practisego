package main

import (
	"fmt"
	"log"

	"github.com/jlaffaye/ftp"
	"github.com/spf13/cobra"
)

var (
	server, user, passwd string
)

func init() {
	rootCmd.PersistentFlags().StringVar(&server, "server", "", "ftp server")
	rootCmd.Flags().StringVar(&user, "user", "", "user")
	rootCmd.Flags().StringVar(&passwd, "passwd", "", "user")
	rootCmd.MarkPersistentFlagRequired("server")
	rootCmd.MarkFlagRequired("user")
	rootCmd.MarkFlagRequired("passwd")
}

var rootCmd = &cobra.Command{
	Use: "",
	Run: func(cmd *cobra.Command, args []string) {
		serverconn, err := ftp.Dial(server)
		if err != nil {
			log.Fatal(err)
		}

		err = serverconn.Login(user, passwd)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Login succesfull")
		pwd, err := serverconn.CurrentDir()

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(pwd)

		names, err := serverconn.NameList(pwd)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(names)

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	Execute()
}
