package ldapx

import (
	"context"
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"net"
	"net/url"
	"strings"
)

type CompositeError struct {
	Errs []error
}

func (err CompositeError) Error() string {
	if err.Errs == nil || len(err.Errs) < 1 {
		return "no errors found"
	}
	return "multiple errors found"
}

type DeepSearchRequest struct {
	*ldap.SearchRequest
	Username, Password string
	AllowAnonymousBind bool
	MaxDepth int
}

// Dialer provides the same interface as net.Dialer and the
// upcoming tls.Dialer.
// This will allow clients to provide TLS configuration if needed.
type Dialer interface {
	Dial(network, address string) (net.Conn, error)
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}

func DeepSearch(ldapURL string, req *DeepSearchRequest) (*ldap.SearchResult, error) {
	refs := []string{ldapURL}

	entries := make([]*ldap.Entry, 0, 100)
	errs := make([]error, 0)

	for {
		if len(refs) < 1 {
			break
		}

		ref := refs[0]
		refs = refs[1:]
		refURL, err := url.Parse(ref)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		// The referral URL might look like ldap://test:123/dc=some,dc=dn
		// where the path (dc=some,dc=dn) is the DN of the referenced entry
		req.BaseDN = strings.TrimLeft(refURL.Path, "/")

		result, err := submitRequest(refURL, req)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if result.Entries != nil {
			entries = append(entries, result.Entries...)
		}
		refs = append(refs, result.Referrals...)
	}

	if len(errs) > 0 {
		return nil, CompositeError{Errs:errs}
	}

	result := new(ldap.SearchResult)
	result.Entries = entries
	return result, nil
}

func submitRequest(ldapURL *url.URL, req *DeepSearchRequest) (*ldap.SearchResult, error) {
	addr := fmt.Sprintf("%s:%s", ldapURL.Hostname(), ldapURL.Port())
	conn, err := ldap.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	if req.Username != "" {
		err = conn.Bind(req.Username, req.Password)
		if err != nil {
			return nil, err
		}
	}
	return conn.Search(req.SearchRequest)
}
