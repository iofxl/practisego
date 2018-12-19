package cmd

import "github.com/spf13/cobra"

var serverCmd = &cobra.Command{
	Use: "server",
	Run: server,
}

func server(cmd *cobra.Command, args []string) {
}
