package main

import (
	"github.com/go-ldap/ldap/v3"
	"net"
	"net/url"
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

	results, err := DeepSearch(new(net.Dialer), ldapURL, DeepSearchRequest{
		SearchRequest: req,
		Username:      "cn=admin,dc=test,dc=com",
		Password:      "admin",
	})

	if err != nil {
		panic(err)
	} else {
		results.Print()
	}

}
