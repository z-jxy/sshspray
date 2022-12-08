/*
Copyright © 2022 zjxy
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sshspray",
	Short: "Tool for testing SSH access",
	Long: `
sshspray is a penetration testing tool for checking SSH access with a set of credentials. 
You can supply individual credentials of files of hosts/users/passwords to test across a range of hosts.
`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cmdSpray.Flags().IntVar(&_sshport, "port", 22, "Specify a port | DEFAULT=22")
	cmdSpray.Flags().StringVarP(&_cmd, "cmd", "x", "", "Execute a command on successful connection | DEFAULT=None")
	cmdSpray.Flags().IntVarP(&_delay, "delay", "d", 10, "Set delay (miliseconds) | DEFAULT=100")

	cmdSpray.Flags().StringVarP(&_passwords, "passwords", "p", "", "Target password")
	cmdSpray.Flags().StringVarP(&_users, "users", "l", "", "Target username")
	cmdSpray.Flags().StringVarP(&_hosts, "hosts", "t", "", "Target host")

	cmdSpray.MarkFlagRequired("hosts")
	cmdSpray.MarkFlagRequired("users")
	cmdSpray.MarkFlagRequired("passwords")

	rootCmd.AddCommand(cmdSpray)
}
