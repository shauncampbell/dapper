package main

import (
	"fmt"
	"github.com/shauncampbell/dapper/pkg/ldap"
	"github.com/spf13/cobra"
)

var (
	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Start the LDAP server",
		Long:  "Start the LDAP server on the specified port",
		Run: func(cmd *cobra.Command, args []string) {
			dapper := ldap.NewServer(baseDN, cfgFile, serverPort)
			if err := dapper.Listen(); err != nil {
				fmt.Println(err.Error())
			}
		},
	}
	serverPort int
)
