package main

import (
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"net/url"
)

func requestAll(ldapURL *url.URL, baseDn string) (*ldap.SearchResult, error) {
	fmt.Printf("searching %s for %s\n", ldapURL.String(), baseDn)
	conn, err := ldap.Dial("tcp", fmt.Sprintf("%s:%s", ldapURL.Hostname(), ldapURL.Port()))
	if err != nil {
		fmt.Printf("error dialing LDAP server: %+v\n", err)
		return nil, err
	}

	err = conn.Bind("cn=admin,dc=test,dc=com", "admin")
	if err != nil {
		fmt.Printf("error authenticating LDAP connection: %+v\n", err)
		return nil, err
	}

	return conn.Search(ldap.NewSearchRequest(
		baseDn,
		ldap.ScopeWholeSubtree,
		ldap.DerefAlways,
		0,
		0,
		false,
		"(objectClass=inetOrgPerson)",
		[]string{"uid", "dn"},
		nil))
}

func main() {
	ldapURL, _ := url.Parse("ldap://localhost:3898/dc=test,dc=com");
	refs := []string{ldapURL.String()}

	entries := make([]*ldap.Entry, 0, 100)

	for {
		if len(refs) < 1 {
			break
		}
		ref := refs[0]
		refs = refs[1:]
		refURL, err := url.Parse(ref)
		if err != nil {
			fmt.Printf("unable to parse referral: %s\n", refURL)
			continue
		}
		result, err := requestAll(refURL, refURL.Path[1:])
		if err != nil {
			fmt.Printf("error during request: %+v\n", err)
			continue
		}
		if result.Entries != nil {
			entries = append(entries, result.Entries...)
		}
		refs = append(refs, result.Referrals...)
	}

	for _, entry := range entries {
		fmt.Printf("%s\n", entry.DN)
	}
}
