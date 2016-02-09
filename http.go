package dnsmasq

import (
	"encoding/json"
	"net/http"
)

// Handler for the leases.
func LeasesHandler(w http.ResponseWriter, r *http.Request) {
	l, err := LoadLeases()

	if err != nil {
		http.Error(w, "Internal error.", http.StatusInternalServerError)
		return
	}

	// Set the output format.
	w.Header().Set("Content-Type", "application/json")

	je := json.NewEncoder(w)
	je.Encode(l)
}
