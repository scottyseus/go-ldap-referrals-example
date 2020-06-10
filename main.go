package main

import (
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"ldap-referral/ldapx"
)

func main() {
	ldapURL := "ldap://localhost:3898/dc=test,dc=com"

	req := ldap.NewSearchRequest(
		"",
		ldap.ScopeWholeSubtree,
		ldap.DerefAlways,
		0,
		0,
		false,
		"(objectClass=inetOrgPerson)",
		[]string{"uid", "dn"},
		nil)

	results, err := ldapx.DeepSearch(ldapURL, &ldapx.DeepSearchRequest{
		SearchRequest: req,
		Username:      "cn=admin,dc=test,dc=com",
		Password:      "admin",
	})

	if results != nil {
		fmt.Println("results:")
		for _, entry := range results.Entries {
			fmt.Printf("\t%+v\n", entry.DN)
		}
	}

	if err != nil {
		fmt.Println("errors:")
		if compErr, ok := err.(ldapx.CompositeError); ok {
			for _, e := range compErr.Errs {
				fmt.Println("\t" + e.Error())
			}
		} else {
			fmt.Printf("\t%+v\n", err)
		}
	}
}
