package dnsmasq

import (
	// 	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	TestJwtTokenValid           = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJraWQiOjAsInNjb3BlcyI6WyJmb28iLCJiYXIiLCJiYXoiXSwidXNlcmlkIjoicGV6YSJ9.bIBFt2RQeEyhC89UHsAN3fBAWQo-vHfidExiJlhYO9A`
	TestJwtTokenBadSignature    = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJraWQiOjAsInNjb3BlcyI6WyJmb28iLCJiYXIiLCJiYXoiXSwidXNlcmlkIjoicGV6YSJ9.bIBFt2RQeEyhC89UHsAN3fBAWQo-vHfidExiJlhYO9a`
	TestJwtTokenTamperedHeader  = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpxVCJ9.eyJraWQiOjAsInNjb3BlcyI6WyJmb28iLCJiYXIiLCJiYXoiXSwidXNlcmlkIjoicGV6YSJ9.bIBFt2RQeEyhC89UHsAN3fBAWQo-vHfidExiJlhYO9A`
	TestJwtTokenTamperedPayload = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJraWQiOjAsInNjb3BlcyI6WyJmb28iLCJiYXIiLcJiYXoiXNlcmlkIjoicGV6YSJ9.bIBFt2RQeEyhC89UHsAN3fBAWQo-vHfidExiJlhYO9A`
)

var (
	TestSecret = []byte(`geheimgeheim`)
)

func TestNewJwt(t *testing.T) {
	// User is valid. Create a jwt response.
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["kid"] = 0
	token.Claims["userid"] = "peza"
	token.Claims["scopes"] = []string{
		"foo",
		"bar",
		"baz",
	}
	ts, err := token.SignedString([]byte(TestSecret))
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	_ = ts
	// t.Logf("%s", ts)
}

func TestNewJwt1(t *testing.T) {
	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)
	// Set some claims
	token.Claims["foo"] = "bar"
	token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(TestSecret)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	_ = tokenString
	// t.Logf("%s", tokenString)
}

func TestParseValidJwt(t *testing.T) {
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		return TestSecret, nil
	}
	token, err := jwt.Parse(TestJwtTokenValid, keyFunc)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	if !token.Valid {
		t.Errorf("Token invalid.")
	}
	// t.Logf("%#v", token)
}

func TestParseJwtBadSignature(t *testing.T) {
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		return TestSecret, nil
	}
	token, err := jwt.Parse(TestJwtTokenBadSignature, keyFunc)
	if err == nil {
		t.Errorf("Token should be invalid.")
	}
	if token.Valid {
		t.Errorf("Token should be invalid.")
	}
}

func TestParseJwtTamperedHeader(t *testing.T) {
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		return TestSecret, nil
	}
	token, err := jwt.Parse(TestJwtTokenTamperedHeader, keyFunc)
	if err == nil {
		t.Errorf("Token should be invalid.")
	}
	if token.Valid {
		t.Errorf("Token should be invalid.")
	}
}

func TestParseJwtTamperedPayload(t *testing.T) {
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		return TestSecret, nil
	}
	token, err := jwt.Parse(TestJwtTokenTamperedPayload, keyFunc)
	if err == nil {
		t.Errorf("Token should be invalid.")
	}
	if token.Valid {
		t.Errorf("Token should be invalid.")
	}
}

func TestAuthRequest(t *testing.T) {
}
