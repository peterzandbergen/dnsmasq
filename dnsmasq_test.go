package dnsmasq

import (
	"bytes"
	"encoding/json"
	"net"
	// "time"

	"net/http/httptest"
	"testing"
)

const dnsmasqLeases = `1454174694 0c:8b:fd:ac:16:74 192.168.2.102 CT00428 01:0c:8b:fd:ac:16:74
1454174265 c8:d1:0b:51:3c:1d 192.168.2.145 Windows-Phone 01:c8:d1:0b:51:3c:1d
1454174431 60:92:17:79:b0:6e 192.168.2.87 iPhone 01:60:92:17:79:b0:6e
1454188541 20:64:32:5d:c7:79 192.168.2.65 android-634c956e97330182 01:20:64:32:5d:c7:79
1454178611 00:21:5c:69:2b:03 192.168.2.114 zandbergen-nb 01:00:21:5c:69:2b:03
1454200072 80:19:34:85:cc:26 192.168.2.86 peza-nb-hp 01:80:19:34:85:cc:26
1454194272 00:f7:6f:ec:d1:2b 192.168.2.133 iPadvaneEijsden 01:00:f7:6f:ec:d1:2b
1454186150 cc:c3:ea:5f:6c:67 192.168.2.132 android-7b3624c64a34b21c 01:cc:c3:ea:5f:6c:67
1454188602 6c:fa:a7:63:ae:37 192.168.2.158 DESKTOP-JPBEM49 01:6c:fa:a7:63:ae:37
1454189265 00:90:a9:42:41:53 192.168.2.67 MediaServer1 *
1454201763 70:77:81:7c:85:85 192.168.2.79 android-4cb72593965678da *
1454187972 00:22:61:cc:e6:a2 192.168.2.66 * *
1454198372 ec:88:92:72:c4:2d 192.168.2.73 android-d7e4398573ab8818 01:ec:88:92:72:c4:2d
1454200693 c4:12:f5:2c:ed:15 192.168.2.97 ap-woonkamer-1 01:c4:12:f5:2c:ed:15
1454200887 f8:e9:03:e6:1e:93 192.168.2.91 ap-router-1 01:f8:e9:03:e6:1e:93
1454196518 f8:e9:03:e6:20:64 192.168.2.146 ap-kantoor-1 01:f8:e9:03:e6:20:64
1454195503 8c:dc:d4:e6:43:8c 192.168.2.98 HPC91022 01:8c:dc:d4:e6:43:8c
1454184751 68:94:23:d1:b6:c5 192.168.2.81 GabyTablet 01:68:94:23:d1:b6:c5
1454188028 54:60:09:f8:8f:48 192.168.2.137 Chromecast *
`

const jsonLeaseIndent = `{
		    "ExpireTime": 1454086301,
		    "ClientMac": "01:23:45:67:89:ab",
		    "ClientIP": "192.168.1.2",
		    "HostName": "hostname1",
		    "ClientId": "01:0c:8b:fd:ac:16:74"
		}`

func MustParseMac(ha string) MacAddr {
	a, err := net.ParseMAC(ha)
	if err != nil {
		panic(err.Error())
	}
	return MacAddr(a)
}

var (
	testExpireTime int64   = 1454086301
	testClientMac  MacAddr = MustParseMac("01:02:03:04:05:06")
	testIp         IPAddr  = IPAddr(net.ParseIP("192.168.1.2"))
	testClienId    string  = "01:0c:8b:fd:ac:16:74"
	testHostname   string  = "testhostname1"
)

const (
	jsonLease = `{"ExpireTime":1454086301,"ClientMac":"01:02:03:04:05:06","ClientIP":"192.168.1.2","HostName":"testhostname1","ClientId":"01:0c:8b:fd:ac:16:74"}`
)

func TestMarshal(t *testing.T) {
	l := &Lease{
		ExpireTime: testExpireTime,
		ClientMac:  testClientMac,
		ClientIP:   testIp,
		ClientId:   testClienId,
		HostName:   testHostname,
	}
	b, err := json.Marshal(l)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	if string(b) != jsonLease {
		t.Errorf("\nexpected: %s\nreceived: %s", jsonLease, string(b))
	}
}

func TestUnmarshalMacAddr(t *testing.T) {
	var adr MacAddr
	// Assign the mac address to a unmarshaller interface.
	var _ json.Unmarshaler = &adr

	err := json.Unmarshal([]byte("\"01:23:45:67:89:ab\""), &adr)
	if err != nil {
		t.Errorf("error: %s", err.Error())
	}

	//	t.Logf("%#v", adr.String())
}

// TestUnmarshalIP tests if unmarshaling works correctly.
func TestUnmarshalIP(t *testing.T) {
	var ipaddr IPAddr
	var _ json.Unmarshaler = &ipaddr

	err := json.Unmarshal([]byte("\"192.168.1.2\""), &ipaddr)
	if err != nil {
		t.Errorf("error: %s", err.Error())
	}

}

func TestUnmarshalLease(t *testing.T) {
	var l Lease

	err := json.Unmarshal([]byte(jsonLease), &l)
	if err != nil {
		t.Errorf("error: %s", err.Error())
	}
	//t.Logf("%s", l.String())
}

func TestParseLeases(t *testing.T) {
	var ls []Lease
	var err error
	r := bytes.NewBufferString(dnsmasqLeases)
	if ls, err = ParseLeases(r); err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	if len(ls) != 19 {
		t.Errorf("Expected %d lines, parsed %d lines.", 19, len(ls))
	}
}

func TestParseLeasesJson(t *testing.T) {
	var ls []Lease
	var err error
	r := bytes.NewBufferString(dnsmasqLeases)
	if ls, err = ParseLeases(r); err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	if len(ls) != 19 {
		t.Errorf("Expected %d lines, parsed %d lines.", 19, len(ls))
	}
	// b, err := json.MarshalIndent(&struct{ Leases []Lease }{ls}, "  ", "  ")
	b, err := json.MarshalIndent(ls, "  ", "  ")
	t.Log(string(b))
}

func TestLeasesServer(t *testing.T) {
	w := httptest.NewRecorder()
	LeasesServer(w, nil)
	t.Logf("%d - %s", w.Code, w.Body.String())
}
