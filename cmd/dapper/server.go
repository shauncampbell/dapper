package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Start the LDAP server",
		Long:  "Start the LDAP server on the specified port",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("hello world")
		},
	}
	serverPort int
)
