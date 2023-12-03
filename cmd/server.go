/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/rahn-it/svalin/config"
	"github.com/rahn-it/svalin/rpc"
	"github.com/rahn-it/svalin/server"

	"github.com/spf13/cobra"
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

		profile, err := config.OpenProfile("server")
		if err != nil {
			panic(err)
		}

		config := profile.Config()
		config.BindFlags(cmd.Flags())

		server, err := server.Open(profile)
		if err != nil {
			panic(err)
		}

		// catch interrupt and gracefully shut down server
		wg := sync.WaitGroup{}

		wg.Add(1)
		go func() {
			err := server.Run()
			if err != nil {
				if errors.Is(err, rpc.ErrRpcNotRunning) {
					log.Printf("Server was stopped")
				} else {
					logErr := fmt.Errorf("error running server: %w", err)
					log.Println(logErr)
				}
			}
			wg.Done()
		}()

		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

		wg.Add(1)
		go func() {
			<-interrupt
			err := server.Close(200, "OK")
			if err != nil {
				err := fmt.Errorf("error shutting down program: error closing server: %w", err)
				log.Println(err)
			} else {
				log.Println("Server closed without errors")
			}
			wg.Done()
		}()

		wg.Wait()

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
