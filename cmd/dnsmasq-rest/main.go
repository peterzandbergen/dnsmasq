package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/justinas/alice"
	"github.com/peterzandbergen/dnsmasq"
)

const (
	TlsCertFile = "../../../../../../certs/selfsigned.myhops.com.crt"
	TlsKeyFile  = "../../../../../../certs/selfsigned.myhops.com.key"
)

type fakeDirectory struct {
}

func (f *fakeDirectory) CheckCredentials(uid, pwd string) (*dnsmasq.Profile, error) {
	if uid != "peza" || pwd != "secret" {
		fmt.Printf("Bad uid/pwd: %s/%s\n", uid, pwd)
		return nil, errors.New("bad credentials")
	}
	return &dnsmasq.Profile{
		UserName: "peter.zandbergen",
		Scopes: []string{
			"Trala",
		},
	}, nil
}

func main() {
	// Start the http server.
	http.Handle("/leases", alice.New(dnsmasq.JwtAuthMiddleWare()).ThenFunc(dnsmasq.LeasesHandler))

	fd := &fakeDirectory{}
	authHandler := dnsmasq.NewAuthenticator(fd, []byte(dnsmasq.JwtSecret))
	http.Handle("/authenticate", authHandler)
	http.ListenAndServeTLS(":8967", TlsCertFile, TlsKeyFile, nil)
}
