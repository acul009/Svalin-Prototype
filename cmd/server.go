/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"rahnit-rmm/rpc"

	"github.com/spf13/cobra"

	"github.com/quic-go/quic-go"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("server called")

		port := 1234

		udpConn, err := net.ListenUDP("udp4", &net.UDPAddr{Port: port})
		if err != nil {
			panic(err)
		}
		// ... error handling
		tr := quic.Transport{
			Conn: udpConn,
		}

		tlsConf := &tls.Config{}

		quicConf := &quic.Config{}

		ln, err := tr.Listen(tlsConf, quicConf)

		if err != nil {
			panic(err)
		}

		fmt.Printf("\nListening on localhost:%d\n", port)

		for {
			conn, err := ln.Accept(context.Background())
			if err != nil {
				log.Printf("Error accepting QUIC connection: %v", err)
				continue
			}
			// ... error handling

			commands := rpc.NewCommandCollection()

			go rpc.ServeSession(conn, commands)
		}

	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
