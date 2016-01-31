// dnsmasq project dnsmasq.go
package dnsmasq

import (
	"bufio"
	"errors"
	"io"
	"net"
	"strconv"
	"strings"
)

// Type MacAddr is used to allow custom String() and Json Marshaling.
type MacAddr net.HardwareAddr

// Type IPAddr is used to allow custom String() and Json Marshaling.
type IPAddr net.IP

// String returns a readable string representation of a hardware address.
func (a MacAddr) String() string {
	return net.HardwareAddr(a).String()
}

func ParseMacAddr(s string) (MacAddr, error) {
	if s == "*" {
		return MacAddr{}, nil
	}
	if hw, err := net.ParseMAC(s); err != nil {
		return nil, err
	} else {
		return MacAddr(hw), nil
	}
}

// MarshalJSON marshals the hardware address to a valid JSON value,
// using the value returned by String.
func (a MacAddr) MarshalJSON() ([]byte, error) {
	return []byte("\"" + a.String() + "\""), nil
}

// Unmarshal takes the representation and sets itself.
func (a *MacAddr) UnmarshalJSON(b []byte) error {
	s := string(b[1 : len(b)-1])
	aa, _ := net.ParseMAC(s)
	a1 := MacAddr(aa)
	*a = a1
	return nil
}

// String returns the string representation of the ip address.
func (ip IPAddr) String() string {
	return net.IP(ip).String()
}

func ParseIpAddr(s string) IPAddr {
	return IPAddr(net.ParseIP(s))
}

// MarshalJSON is implemented because we need custom marshalling.
func (ip IPAddr) MarshalJSON() ([]byte, error) {
	return []byte("\"" + ip.String() + "\""), nil
}

// UnmarshalJSON implements the Unmarshaler interface.
func (ip *IPAddr) UnmarshalJSON(b []byte) error {
	// Strip the quotes.
	s := string(b[1 : len(b)-1])
	// Parse the ip address and set myself.
	i := net.ParseIP(s)
	*ip = IPAddr(i)

	return nil
}

// Lease contains the information of a dnsmasq lease.
type Lease struct {
	ExpireTime int64
	ClientMac  MacAddr
	ClientIP   IPAddr
	HostName   string
	ClientId   string
}

func (l Lease) String() string {
	return strconv.FormatInt(l.ExpireTime, 10) + "|" +
		l.ClientMac.String() + "|" +
		l.ClientIP.String() + "|" +
		l.HostName + "|" +
		l.ClientId
}

func ParseLease(s string) (Lease, error) {
	var err error
	var l Lease

	parts := strings.Split(s, " ")
	// Check the number of parts.
	if len(parts) != 5 {
		return l, errors.New("ParseLease error, too few parameters.")
	}
	// Parse the individual elements.
	// 0: ExpireTime
	if l.ExpireTime, err = strconv.ParseInt(parts[0], 10, 64); err != nil {
		return l, errors.New("ParseLease error, expiretime failed.")
	}
	// 1: ClientMac
	if l.ClientMac, err = ParseMacAddr(parts[1]); err != nil {
		return l, errors.New("ParseLease error, ClientMac failed.")
	}
	// 2: ClientIP
	l.ClientIP = IPAddr(net.ParseIP(parts[2]))
	// 3: HostName
	l.HostName = parts[3]
	// 4: ClientId
	l.ClientId = parts[4]
	return l, nil
}

/*
ParseLease parses the lines from the dnsmasq.leases files and returns a
slice with the found leases, or an error.

First return value is nil if an error occured.
*/
func ParseLeases(r io.Reader) ([]Lease, error) {
	var ls []Lease

	// Get a line scanner.
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			if l, err := ParseLease(line); err != nil {
				return nil, errors.New("ParseLeaeses failed: " + err.Error())
			} else {
				ls = append(ls, l)
			}
		}
	}

	return ls, nil
}
