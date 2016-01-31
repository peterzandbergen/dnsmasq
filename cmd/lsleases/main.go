// lsleases.go project main.go
package main

import (
	"dnsmasq"
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	locations = []string{
		"",
		"",
	}
)

func findFile() (io.ReadCloser, error) {
	for _, fn := range locations {
		if f, err := os.Open(fn); err != nil {
			return f, nil
		}
	}
	return nil, errors.New("file not found or not accessible.")
}

func loadLeases() []dnsmasq.Lease {
	r, err := findFile()
	if err != nil {
		return nil
	}
	defer r.Close()
	ls, err := dnsmasq.ParseLeases(r)
	if err != nil {
		return nil
	}
	return ls
}

func listResult(ls []dnsmasq.Lease) {
	fmt.Println("[%d]Lease\n", len(ls))
	for _, l := range ls {
		fmt.Println(l.String())
	}
}

func main() {
	ls := loadLeases()
	listResult(ls)
}
