package go_twitter_oauth

import (
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type Config struct {
	ClientID          string
	RedirectURI       string
	ClientSecret      string
	Scopes            []string
	AuthorizeEndpoint string
	TokenEndPoint     string
}

type OAuth2Param interface {
	createClient() *Client
	createTokenRequest(code string, codeVerifier string) *http.Request
	buildAuthURL(session *PKCESession) string
}

var _ OAuth2Param = (*Config)(nil)

func (config *Config) buildAuthURL(session *PKCESession) string {
	scopesString := strings.Join(config.Scopes, " ")
	u, err := url.Parse(config.AuthorizeEndpoint)
	if err != nil {
		log.Fatal(err)
	}

	q := u.Query()
	q.Set("response_type", "code")
	q.Set("client_id", config.ClientID)
	q.Set("redirect_uri", config.RedirectURI)
	q.Set("scope", scopesString)
	q.Set("state", session.State)
	q.Set("code_challenge", session.CodeChallenge)
	q.Set("code_challenge_method", session.CodeChallengeMethod)

	u.RawQuery = q.Encode()
	escapedURL := regexp.MustCompile(`([^%])(\+)`).ReplaceAllString(u.String(), "$1%20")

	return escapedURL
}

func (config *Config) createClient() *Client {
	ch := make(chan string)
	return &Client{
		config:  config,
		session: nil,
		ch:      ch,
	}
}

func (config *Config) createTokenRequest(code string, codeVerifier string) *http.Request {
	values := url.Values{}
	values.Set("code", code)
	values.Add("grant_type", "authorization_code")
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
