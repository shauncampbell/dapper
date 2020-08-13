package main

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "dapper",
		Short: "Dapper is a very fast simple ldap server",
		Long:  `A fast easy to use LDAP server for use with home labs`,
	}

	cfgFile    string
	baseDN	   string

)

func init() {
	cobra.OnInitialize()

	// Add flags to server command
	serverCmd.PersistentFlags().IntVarP(&serverPort, "port", "p", 389, "port to run LDAP server on")
	serverCmd.PersistentFlags().StringVarP(&baseDN, "base", "b", "", "the base DN for which this server will accept requests")
	serverCmd.MarkPersistentFlagRequired("base")
	// Add flags to the root command
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "f", "", "config file (default is $HOME/.dapper.yaml)")
	rootCmd.MarkPersistentFlagRequired("config")

	// Add the sub commands to the root
	rootCmd.AddCommand(serverCmd)
}

func main() {
	_ = rootCmd.Execute()
}
