package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/peterzandbergen/dnsmasq"
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
	http.HandleFunc("/leases", dnsmasq.LeasesServer)

	fd := &fakeDirectory{}
	authHandler := dnsmasq.NewAuthenticator(fd, []byte(dnsmasq.JwtSecret))
	http.Handle("/authenticate", authHandler)
	http.ListenAndServe(":8967", nil)
}
