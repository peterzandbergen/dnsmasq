package main

import (
	"net/http"

	"github.com/peterzandbergen/dnsmasq"
)

func main() {
	// Start the http server.
	http.HandleFunc("/leases", dnsmasq.LeasesServer)
	http.ListenAndServe(":8967", nil)
}
