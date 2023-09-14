/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"rahnit-rmm/connection"
	"rahnit-rmm/rpc"
	"time"

	"github.com/spf13/cobra"
)

// cliCmd represents the cli command
var cliCmd = &cobra.Command{
	Use:   "cli",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("cli called")

		addr := "localhost:1234"

		conn, err := connection.CreateClient(context.Background(), addr)

		stream, err := conn.OpenStreamSync(context.Background())
		if err != nil {
			panic(err)
		}

		session := rpc.NewRpcSession(stream)

		rpcCmd := &rpc.PingCmd{}

		session.SendCommand(rpcCmd)

		err = stream.Close()
		if err != nil {
			panic(err)
		}

		time.Sleep(time.Second)

	},
}

func init() {
	rootCmd.AddCommand(cliCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cliCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cliCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
