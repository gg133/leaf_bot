package pocket

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	host                 = "https://getpocket.com/v3"
	authorizeUrl         = "https://getpocket.com/auth/authorize?request_token=%s&redirect_uri=%s"
	endpointAdd          = "/add"
	endpointRequestToken = "/oauth/request"
	endpointAuth         = "/oauth/authorize"
	xErrorHeader         = "X-Error"
	defaultTimeout       = 5 * time.Second
)

type (
	requestTokenRequest struct {
		ConsumerKey string `json:"consumer_key"`
		RedirectURI string `json:"redirect_uri"`
	}

	authRequest struct {
		ConsumerKey string `json:"consumer_key"`
		Code        string `json:"code"`
	}

	AuthResponse struct {
		AccessToken string `json:"access_token"`
		Username    string `json:"user_agent"`
	}

	addRequest struct {
		URL         string `json:"url"`
		Title       string `json:"title,omitempty"`
		Tags        string `json:"tags.omitempty"`
		AccessToken string `json:"access_token"`
		ConsumerKey string `json:"consumer_key"`
	}

	AddInput struct {
		URL         string
		Title       string
		Tags        []string
		AccessToken string
	}
)

type Client struct {
	client      *http.Client
	consumerKey string
}

func (c *Client) doHTTP(ctx context.Context, endpoint string, body interface{}) (url.Values, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return url.Values{}, errors.WithMessage(err, "failed to marshal input body")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, host+endpoint, bytes.NewBuffer(b))
	if err != nil {
		return url.Values{}, errors.WithMessage(err, "failed to create new request")
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF8")

	resp, err := c.client.Do(req)
	if err != nil {
		return url.Values{}, errors.WithMessage(err, "failed to send http request")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Sprintf("API Error: %s", resp.Header.Get(xErrorHeader))
		return url.Values{}, errors.New(err)
	}

	respB, err := io.ReadAll(resp.Body)
	if err != nil {
		return url.Values{}, errors.WithMessage(err, "failed to read request body")
	}

	values, err := url.ParseQuery(string(respB))
	if err != nil {
		return url.Values{}, errors.WithMessage(err, "failed to parse response body")
	}

	return values, nil
}

func (a AddInput) validate() error {
	if a.URL == "" {
		return errors.New("required URL values is empty")
	}

	if a.AccessToken == "" {
		return errors.New("access token is empty")
	}

	return nil
}

func (a AddInput) generateRequest(consumerKey string) addRequest {
	return addRequest{
		URL:         a.URL,
		Tags:        strings.Join(a.Tags, ","),
		Title:       a.Title,
		AccessToken: a.AccessToken,
		ConsumerKey: consumerKey,
	}
}

func NewClient(consumerKey string) (*Client, error) {
	if consumerKey == "" {
		return nil, errors.New("consumer key is empty")
	}

	return &Client{
		client: &http.Client{
			Timeout: defaultTimeout,
		},
		consumerKey: consumerKey,
	}, nil
}

func (c *Client) GetRequestToken(ctx context.Context, redirectUrl string) (string, error) {
	inp := &requestTokenRequest{
		ConsumerKey: c.consumerKey,
		RedirectURI: redirectUrl,
	}

	values, err := c.doHTTP(ctx, endpointRequestToken, inp)
	if err != nil {
		return "", err
	}

	if values.Get("code") == "" {
		return "", errors.New("empty request token in API response")
	}

	return values.Get("code"), nil
}

func (c *Client) GetAuthURL(requestToken, redirectUrl string) (string, error) {
	if requestToken == "" || redirectUrl == "" {
		return "", errors.New("empty params")
	}

	return fmt.Sprintf(authorizeUrl, requestToken, redirectUrl), nil
}

func (c *Client) Auth(ctx context.Context, requestToken string) (*AuthResponse, error) {
	if requestToken == "" {
		return nil, errors.New("empty request token")
	}

	inp := &authRequest{
		Code:        requestToken,
		ConsumerKey: c.consumerKey,
	}

	values, err := c.doHTTP(ctx, endpointAuth, inp)
	if err != nil {
		return nil, err
	}

	accessToken, username := values.Get("access_token"), values.Get("username")
	if accessToken == "" {
		return nil, errors.New("empty access token in API response")
	}

	return &AuthResponse{
		AccessToken: accessToken,
		Username:    username,
	}, nil
}

func (c *Client) Add(ctx context.Context, input AddInput) error {
	if err := input.validate(); err != nil {
		return err
	}

	req := input.generateRequest(c.consumerKey)
	_, err := c.doHTTP(ctx, endpointAdd, req)

	return err
}
