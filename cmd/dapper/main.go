package main

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	rootCmd = &cobra.Command{
		Use:   "dapper",
		Short: "Dapper is a very fast simple ldap server",
		Long:  `A fast easy to use LDAP server for use with home labs`,
	}

	cfgFile    string

)

func init() {
	cobra.OnInitialize(initConfig)

	// Add flags to server command
	serverCmd.PersistentFlags().IntVarP(&serverPort, "port", "p", 389, "port to run LDAP server on")

	// Add flags to the root command
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "f", "", "config file (default is $HOME/.dapper.yaml)")

	// Add the sub commands to the root
	rootCmd.AddCommand(serverCmd)
}

func main() {
	_ = rootCmd.Execute()
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			er(err)
		}

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".dapper")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}
