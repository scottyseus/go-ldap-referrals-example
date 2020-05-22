package main

import (
	"context"
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"net"
	"net/url"
)

type DeepSearchRequest struct {
	*ldap.SearchRequest
	Username, Password string
	AllowAnonymousBind bool
}

// Dialer provides the same interface as net.Dialer and the
// upcoming tls.Dialer.
// This will allow clients to provide TLS configuration if needed.
type Dialer interface {
	Dial(network, address string) (net.Conn, error)
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}

func DeepSearch(dialer Dialer, ldapURL *url.URL, req DeepSearchRequest) (*ldap.SearchResult, error) {
	refs := []string{ldapURL.String()}

	entries := make([]*ldap.Entry, 0, 100)
	result := new(ldap.SearchResult)

	for {
		if len(refs) < 1 {
			break
		}

		ref := refs[0]
		refs = refs[1:]
		refURL, err := url.Parse(ref)
		if err != nil {
			continue
		}
		result, err := requestAll(dialer, refURL, req)
		if err != nil {
			continue
		}
		if result.Entries != nil {
			entries = append(entries, result.Entries...)
		}
		refs = append(refs, result.Referrals...)
	}

	return result, nil
}

func requestAll(dialer Dialer, ldapURL *url.URL, req DeepSearchRequest) (*ldap.SearchResult, error) {
	netConn, err:= dialer.Dial("tcp", fmt.Sprintf("%s:%s", ldapURL.Hostname(), ldapURL.Port()))
	if err != nil {
		return nil, err
	}
	conn := ldap.NewConn(netConn, ldapURL.Scheme == "ldaps")

	if req.Username != "" {
		err = conn.Bind(req.Username, req.Password)
		bindReq := ldap.NewSimpleBindRequest(req.Username, req.Password, req.Controls)
		bindReq.AllowEmptyPassword = req.AllowAnonymousBind
		conn.SimpleBind(bindReq)
		if err != nil {
			return nil, err
		}
	}

	return conn.Search(req.SearchRequest)
}
