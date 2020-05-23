package main

import (
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"ldap-referral/ldapx"
	"net"
	"net/url"
	"time"
)

func main() {
	ldapURL, _ := url.Parse("ldap://localhost:3898/dc=test,dc=com");

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

	dialer := new(net.Dialer)

	dialer.Timeout = time.Second * 4

	results, err := ldapx.DeepSearch(dialer, ldapURL, ldapx.DeepSearchRequest{
		SearchRequest: req,
		Username:      "cn=admin,dc=test,dc=com",
		Password:      "admin",
	})

	fmt.Printf("result: %+v\nerrors: %+v\n", results, err)

	if compErr, ok := err.(ldapx.CompositeError); ok {
		for _, e := range compErr.Errs {
			fmt.Println("\t" + e.Error())
		}
	}
}
