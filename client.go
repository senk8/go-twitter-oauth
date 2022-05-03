package go_twitter_oauth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

const (
	listenAddr = "127.0.0.1:3000"
)

type Client struct {
	client *http.Client
	config  *Config
	session *Session
	ch      chan *TokenResponse
}

type OAuth2Client interface {
	ExecFlow(ctx context.Context) (*TokenResponse, error)
	authHandler(w http.ResponseWriter, r *http.Request)
	callbackHandler(w http.ResponseWriter, r *http.Request)
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

var _ OAuth2Client = (*Client)(nil)

func New(config *Config) *Client {
	ch := make(chan *TokenResponse)
	return &Client{
		client:  http.DefaultClient,
		config:  config,
		session: nil,
		ch:      ch,
	}
}

func (c *Client) ExecFlow(ctx context.Context) (*TokenResponse, error) {
	// ローカルにサーバーを立てて、リダイレクトを待機します
	mux := http.NewServeMux()
	mux.HandleFunc("/auth", c.authHandler)
	mux.HandleFunc("/callback", c.callbackHandler)
	srv := &http.Server{
		Addr:    listenAddr,
		Handler: mux,
	}

	// サーバーを任意のタイミングで終了できるようにgoroutineで実行します
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 認可フローの終了を待ちます
	tokenResponse := <-c.ch

	// 終了したらサーバーを停止します
	if err := srv.Shutdown(ctx); err != nil {
		return nil, err
	}

	return tokenResponse, nil
}

func (c *Client) authHandler(w http.ResponseWriter, r *http.Request) {
	c.session = newSession()
	authzURL := c.config.buildAuthzURL(c.session)
	w.Header().Set("Location", authzURL)
	w.WriteHeader(http.StatusFound)
	return
}

func (c *Client) callbackHandler(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if c.session.State != state {
		log.Println("invalid state")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	code := r.URL.Query().Get("code")
	req := c.config.buildTokenRequest(code, c.session.CodeVerifier)

	res, err := c.client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	var tokenResponse TokenResponse
	if err := json.NewDecoder(res.Body).Decode(&tokenResponse); err != nil {
		log.Fatal(err)
	}
	c.ch <- &tokenResponse
	return
}
