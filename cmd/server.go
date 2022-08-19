/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	barkserver "github.com/c1emon/barkbridge/src/barkserver"
	"github.com/spf13/cobra"
)

var serverAddress string
var daemon bool

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start brak bridge server",
	Run: func(cmd *cobra.Command, args []string) {
		barkserver.Push("a", barkserver.Message{})
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")
	serverCmd.PersistentFlags().StringVarP(&serverAddress, "bark-server", "a", "http://127.0.0.1:8080", "Bark server address")
	serverCmd.PersistentFlags().BoolVarP(&daemon, "daemon", "d", false, "Server in daemon mode")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
