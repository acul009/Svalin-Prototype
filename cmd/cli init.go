/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		// // address is required
		// addr := cmd.Flag("addr").Value.String()
		// if len(addr) == 0 {
		// 	fmt.Println("Address is required (--addr)")
		// 	return
		// }

		// // server-name is required
		// nameForServer := cmd.Flag("server-name").Value.String()
		// if len(nameForServer) == 0 {
		// 	fmt.Println("Server name is required (--server-name)")
		// 	return
		// }

		// err := config.SetSubdir("client")
		// if err != nil {
		// 	panic(err)
		// }

		// // create root user if missing
		// var rootPassword []byte

		// _, err = pki.Root.Get()
		// if err != nil {
		// 	if errors.Is(err, pki.ErrNoRootCert) {
		// 		fmt.Println("No root certificate found, generating one")

		// 		rootUser, err := util.AskForString("Enter username for root")
		// 		if err != nil {
		// 			panic(err)
		// 		}

		// 		rootPassword, err = util.AskForNewPassword("Enter password to encrypt the root certificate")
		// 		if err != nil {
		// 			panic(err)
		// 		}

		// 		err = pki.InitRoot(rootUser, rootPassword)
		// 		if err != nil {
		// 			panic(err)
		// 		}
		// 	} else {
		// 		panic(err)
		// 	}
		// } else {
		// 	fmt.Println("Root certificate found, skipping CA generation")
		// 	rootPassword, err = util.AskForPassword("Enter password to decrypt the root certificate")
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// }

		// pki.Unlock(rootPassword)

		// // err = rpc.SetupServer(addr, rootPassword, nameForServer)
		// // if err != nil {
		// // 	panic(err)
		// // }

		// // rootCert, err := pki.Root.Get()
		// if err != nil {
		// 	panic(err)
		// }

		// // rootKey, err := pki.GetCurrentKey()
		// if err != nil {
		// 	panic(err)
		// }

		// // rootTotp, current, err := util.AskForNewTotp(rootCert.Subject.CommonName)
		// if err != nil {
		// 	panic(err)
		// }

		// // reg, err := rpc.NewRegisterUserCmd(rootCert, rootKey, rootPassword, rootTotp, current)
		// if err != nil {
		// 	panic(err)
		// }

		// // ep, err := rpc.ConnectToUpstream(context.Background())
		// if err != nil {
		// 	panic(err)
		// }

		// // session, err := ep.Session(context.Background())
		// if err != nil {
		// 	panic(err)
		// }

		// // err = session.SendCommand(reg)
		// if err != nil {
		// 	panic(err)
		// }

	},
}

func init() {
	cliCmd.AddCommand(initCmd)

	initCmd.Flags().StringP("addr", "a", "", "example-rmm.com:1234")
	initCmd.Flags().StringP("server-name", "n", "", "The name you want to assign to the server")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
