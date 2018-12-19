package cmd

import "github.com/spf13/cobra"

var clientCmd = &cobra.Command{
	Use: "client",
	Run: client,
}

func client(cmd *cobra.Command, args []string) {

}
