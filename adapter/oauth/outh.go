package oauth

import (
	"errors"
	"time"
)

var ErrNotImplemented = errors.New("this method is not yet implemented by the invoking object")

type Config struct {
	ClientID    string
	SecretKey   string
	Scopes      []string
	RedirectURL string
	Endpoint    string
}

type AuthURL struct {
	URL          string
	State        string
	CodeVerifier string
}

type Provider interface {
	Load(*Config)
	PKCEAuthCodeURL() AuthURL
	PKCEExchange(code, codeVerifier string) (*Token, error)
	PKCEGenerator() (codeVerifier string, codeChallenge string)
}

type Token struct {
	Access  string
	Refresh string
	Expiry  time.Time
}

/*
Struct implements all the methods that are defined on the parent struct
It is useful to embed this struct in the struct that implements the parent struct
to act as a default method in case you want to implement only some of the methods and not all
*/

type UnImplemented struct{}

var _ Provider = (*UnImplemented)(nil)

func (unImpl *UnImplemented) Load(*Config) {

}

func (unImpl *UnImplemented) PKCEAuthCodeURL() AuthURL {
	return AuthURL{URL: ErrNotImplemented.Error()}
}

func (unImpl *UnImplemented) PKCEExchange(string, string) (*Token, error) {
	return &Token{}, ErrNotImplemented
}
func (unImpl *UnImplemented) PKCEGenerator() (string, string) {
	return "", ""
}

// ****************** //
