package main

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/go-ldap/ldap/v3"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"ldap-referral/ldapx"
	"os"
	"testing"
)

const imageName = "osixia/openldap:1.3.0"

var ip string
var port nat.Port

type LDAPFixture struct {
	Name string
	LDIF string
}

func TestMain(m *testing.M) {
	ldapPort, err := nat.NewPort("tcp", "389")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        imageName,
		ExposedPorts: []string{ldapPort.Port() + "/" + ldapPort.Proto()},
		Env: map[string]string{
			"LDAP_ORGANISATION": "Test",
			"LDAP_DOMAIN": "test.com",
			"LDAP_BASE_DN": "dc=test,dc=com",
			"LDAP_TLS": "false",
		},
		//VolumeMounts: map[string]string{
		//	"mount1":
		//		"/home/sweidenk/Workspace/go-projects/sweidenk/ldap-referral/testdata/node1.ldif:/container/service/slapd/assets/config/bootstrap/ldif/custom/node.ldif",
		//},
		WaitingFor:   wait.ForLog("slapd starting"),

	}

	ldapC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}
	defer ldapC.Terminate(ctx)
	ip, err = ldapC.Host(ctx)
	if err != nil {
		panic(err)
	}
	port, err = ldapC.MappedPort(ctx, ldapPort)
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func TestDeepSearch(t *testing.T) {
	ldapURL := fmt.Sprintf("ldap://%s:%s/dc=test,dc=com", ip, port.Port())

	req := ldap.NewSearchRequest(
		"",
		ldap.ScopeWholeSubtree,
		ldap.DerefAlways,
		0,
		0,
		false,
		"(objectClass=*)",
		[]string{"uid", "dn"},
		nil)

	results, err := ldapx.DeepSearch(ldapURL, ldapx.DeepSearchRequest{
		SearchRequest: req,
		Username:      "cn=admin,dc=test,dc=com",
		Password:      "admin",
	})

	if !assert.NoError(t, err) {
		return
	}

	assert.NotEmpty(t, results.Entries)
}