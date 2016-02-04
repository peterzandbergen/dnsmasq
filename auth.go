package dnsmasq

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	JwtSecret = `blablablabla`
)

// JWT Authenticator handler for JWT authentication.
type Authenticator struct {
	dir    UserDirectory
	secret []byte
}

type Profile struct {
	UserName string
	Scopes   []string
}

func NewAuthenticator(ud UserDirectory, secret []byte) *Authenticator {
	return &Authenticator{
		dir:    ud,
		secret: secret,
	}
}

type UserDirectory interface {
	CheckCredentials(uid, pwd string) (*Profile, error)
}

func (a *Authenticator) parseJwt(tokenString string) (*jwt.Token, error) {
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		return a.secret, nil
	}
	token, err := jwt.Parse(tokenString, keyFunc)
	return token, err
}

func (a *Authenticator) generateJwt(p *Profile) ([]byte, error) {
	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)
	// Set some claims
	token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	// User is valid. Create a jwt response.
	token.Claims["kid"] = 0
	token.Claims["userid"] = p.UserName
	token.Claims["scopes"] = p.Scopes
	ts, err := token.SignedString(a.secret)
	// Sign and get the complete encoded token as a string

	ts, err = token.SignedString(a.secret)
	if err != nil {
		return nil, err
	}
	return []byte(ts), nil
}

func (a *Authenticator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check the username and password against the directory.
	uid, pwd := r.FormValue("uid"), r.FormValue("pwd")
	p, err := a.dir.CheckCredentials(uid, pwd)
	if err != nil {
		http.Error(w, "Incorrect credentials.", http.StatusUnauthorized)
		return
	}
	ts, err := a.generateJwt(p)
	if err != nil {
		http.Error(w, "Incorrect credentials.", http.StatusUnauthorized)
		return
	}
	w.Header().Add("Content-Type", "text/plain")
	w.Write(ts)
}
