// lsleases.go project main.go
package main

import (
	"github.com/peterzandbergen/dnsmasq"

	"errors"
	"fmt"
	"io"
	"os"
)

var (
	locations = []string{
		"/var/db/dnsmasq.leases",
		"/var/lib/misc/dnsmasq.leases",
		"./dnsmasq.leases",
	}
)

func findFile() (io.ReadCloser, error) {
	for _, fn := range locations {
		if f, err := os.Open(fn); err == nil {
			return f, nil
		}
	}
	return nil, errors.New("file not found or not accessible.")
}

func loadLeases() ([]dnsmasq.Lease, error) {
	r, err := findFile()
	if err != nil {
		return nil, err
	}
	defer r.Close()
	ls, err := dnsmasq.ParseLeases(r)
	if err != nil {
		return nil, err
	}
	return ls, nil
}

func listResult(ls []dnsmasq.Lease) {
	fmt.Printf("[%d]Lease\n", len(ls))
	for _, l := range ls {
		fmt.Println(l.String())
	}
}

func main() {
	ls, err := loadLeases()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	} else {
		listResult(ls)
	}
}
