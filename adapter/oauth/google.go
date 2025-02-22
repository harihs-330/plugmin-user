package oauth

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"time"
	"unsafe"

	"math/rand"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleProvider struct {
	// UnImplemented
	config *oauth2.Config
}

const (
	_letters       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	_letterIdxBits = 6                     // 6 bits to represent a letter index
	_letterIdxMask = 1<<_letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	_letterIdxMax  = 63 / _letterIdxBits   // # of letter indices fitting in 63 bits
)

const (
	UserProfileScope = "https://www.googleapis.com/auth/userinfo.profile"
	USerEmailScope   = "https://www.googleapis.com/auth/userinfo.email"
)

var src = rand.NewSource(time.Now().UnixNano())

var _ Provider = (*GoogleProvider)(nil)

func (gp *GoogleProvider) Load(cfg *Config) {
	gp.config = &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.SecretKey,
		RedirectURL:  cfg.RedirectURL,
		Scopes:       []string{UserProfileScope, USerEmailScope},
		Endpoint:     google.Endpoint,
	}
}

func (gp *GoogleProvider) PKCEGenerator() (string, string) {
	var (
		codeChallenge, codeVerifier string
	)
	codeVerifier = RandomString(50)
	sha2 := sha256.New()
	_, _ = io.WriteString(sha2, codeVerifier)
	// _ = err
	codeChallenge = base64.RawURLEncoding.EncodeToString(sha2.Sum(nil))

	return codeVerifier, codeChallenge
}

// thanks to
// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func RandomString(n int) string {
	byt := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), _letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), _letterIdxMax
		}
		if idx := int(cache & _letterIdxMask); idx < len(_letters) {
			byt[i] = _letters[idx]
			i--
		}
		cache >>= _letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&byt))
}

func (gp *GoogleProvider) PKCEAuthCodeURL() AuthURL {
	state := RandomString(24)
	codeVerifier, codeChallenge := gp.PKCEGenerator()
	URL := gp.config.AuthCodeURL(
		state,
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("code_challenge", codeChallenge))

	return AuthURL{
		URL:          URL,
		State:        state,
		CodeVerifier: codeVerifier,
	}
}

func (gp *GoogleProvider) PKCEExchange(code string, codeVerifier string) (*Token, error) {
	ctx := context.Background()
	token, err := gp.config.Exchange(
		ctx,
		code,
		oauth2.SetAuthURLParam("code_verifier", codeVerifier))
	if err != nil {
		return nil, err
	}

	return &Token{
		Access:  token.AccessToken,
		Refresh: token.RefreshToken,
		Expiry:  token.Expiry,
	}, nil
}

// refer - https://chrisguitarguy.com/2022/12/07/oauth-pkce-with-go/
