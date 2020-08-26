package main

import (
	"encoding/json"
	"fmt"
	ldap2 "github.com/nmcclain/ldap"
	"github.com/shauncampbell/dapper/pkg/console"
	"github.com/shauncampbell/dapper/pkg/ldap"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
	"strings"
)

var (
	searchCmd = &cobra.Command{
		Use:   "search",
		Short: "Search for resources on the LDAP server",
		Long:  "search for resources on the LDAP server",
		Run: func(cmd *cobra.Command, args []string) {
			search(args)
		},
	}
	outputFormat string
)

// UserSearchResult is the output from a search
type UserSearchResult struct {
	DN       string `json:"dn" yaml:"dn"`
	UID      string `json:"uid" yaml:"uid"`
	Forename string `json:"forename" yaml:"forename"`
	Surname  string `json:"surname" yaml:"surname"`
}

// search performs a search against the ldap service.
func search(args []string) {
	// create a new ldap server and load its config but don't start it.
	dapper := ldap.NewServer(baseDN, cfgFile, serverPort)
	dapper.ReloadConfiguration(cfgFile)

	// perform the search using the SearchInternal function.
	// if no argument is specified then search for all dn's
	// otherwise use the query provided as an argument.
	var result []*ldap2.Entry
	var err error
	if len(args) >= 1 {
		result, err = dapper.SearchInternal(args[0])
	} else {
		result, err = dapper.SearchInternal("(dn=*)")
	}

	// if the search errors out then print an error message.
	if err != nil {
		dapper.Logger.Err(err).Msg("search failed")
		return
	}

	// build up an array of results
	var values []interface{}

	for _, res := range result {
		values = append(values, createSearchResult(res))
	}

	if outputFormat == "" {
		out, err := console.Marshal(&values)
		if err != nil {
			dapper.Logger.Err(err).Msg("failed to output to console")
			return
		}
		fmt.Println(string(out))
		fmt.Println("Query returned", len(values), "result(s)")
	} else if outputFormat == "json" {
		out, err := json.Marshal(&values)
		if err != nil {
			dapper.Logger.Err(err).Msg("failed to output to json")
			return
		}
		fmt.Println(string(out))
	} else if outputFormat == "yaml" {
		out, err := yaml.Marshal(&values)
		if err != nil {
			dapper.Logger.Err(err).Msg("failed to output to yaml")
			return
		}
		fmt.Println(string(out))
	}
}

func createSearchResult(res *ldap2.Entry) interface{} {
	vals := res.GetAttributeValues("objectClass")
	for _, v := range vals {
		if strings.ToUpper(v) == strings.ToUpper("posixAccount") {
			return UserSearchResult{UID: res.GetAttributeValue("uid"), DN: res.DN, Forename: res.GetAttributeValue("givenName"), Surname: res.GetAttributeValue("sn")}
		}
	}
	return map[string]string{ "DN": res.DN }
}
