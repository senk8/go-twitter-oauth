package go_twitter_oauth

import (
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)


const (
	responseType = "code"
	grantType = "authorization_code"
)

type Config struct {
	ClientID          string
	RedirectURI       string
	ClientSecret      string
	Scopes            []string
	AuthzEndpoint 	  string
	TokenEndPoint     string
}

type OAuth2Config interface {
	buildAuthzURL(session *Session) string
	buildTokenRequest(code string, codeVerifier string) *http.Request
}

var _ OAuth2Config= (*Config)(nil)

func (config *Config) buildAuthzURL(session *Session) string {
	scopesString := strings.Join(config.Scopes, " ")
	u, err := url.Parse(config.AuthzEndpoint)
	if err != nil {
		log.Fatal(err)
	}

	q := u.Query()
	q.Set("response_type", responseType)
	q.Set("client_id", config.ClientID)
	q.Set("redirect_uri", config.RedirectURI)
	q.Set("scope", scopesString)
	q.Set("state", session.State)
	q.Set("code_challenge", session.CodeChallenge)
	q.Set("code_challenge_method", session.CodeChallengeMethod)

	u.RawQuery = q.Encode()

	r := regexp.MustCompile(`([^%])(\+)`)
	escapedURL := r.ReplaceAllString(u.String(), "$1%20")

	return escapedURL
}

func (config *Config) buildTokenRequest(code string, codeVerifier string) *http.Request {
	values := url.Values{}
	values.Set("code", code)
	values.Add("grant_type", grantType)
	values.Add("redirect_uri", config.RedirectURI)
	values.Add("code_verifier", codeVerifier)

	req, err := http.NewRequest(http.MethodPost, config.TokenEndPoint, strings.NewReader(values.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(config.ClientID, config.ClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}

