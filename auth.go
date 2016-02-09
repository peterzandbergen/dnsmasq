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

// Profile contains the name of the user and the scopes that the user can operate within.
type Profile struct {
	UserName string
	Scopes   []string
}

// NewAuthenticator returns an authenticator that uses the passed UserDirectory to check the
// credentials against.
func NewAuthenticator(ud UserDirectory, secret []byte) *Authenticator {
	return &Authenticator{
		dir:    ud,
		secret: secret,
	}
}

// UserDirectory interface needs to be implemented by UserDirectory providers.
type UserDirectory interface {
	CheckCredentials(uid, pwd string) (*Profile, error)
}

// parseJWT parses a JWT token, validates it and returns it if valild.
func (a *Authenticator) parseJwt(tokenString string) (*jwt.Token, error) {
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		return a.secret, nil
	}
	token, err := jwt.Parse(tokenString, keyFunc)
	return token, err
}

// generateJWT generates a new token for the profile. Token is returned as an array of bytes or nil if error.
func (a *Authenticator) generateJwt(p *Profile) ([]byte, error) {
	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)
	// Set some claims
	// TODO: complete the claims.
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

// ServeHTTP validates the query parameter uid and pwd and returns a JWT token in the
// response if valid. The token contains the user name and the scopes.
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
