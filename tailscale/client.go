package tailscale

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type Device struct {
	Addresses                 []string `json:"addresses"`
	ID                        string   `json:"id"`
	NodeID                    string   `json:"nodeId"`
	User                      string   `json:"user"`
	Name                      string   `json:"name"`
	Hostname                  string   `json:"hostname"`
	ClientVersion             string   `json:"clientVersion"`
	UpdateAvailable           bool     `json:"updateAvailable"`
	OS                        string   `json:"os"`
	Created                   string   `json:"created"`
	ConnectedToControl        bool     `json:"connectedToControl"`
	LastSeen                  string   `json:"lastSeen"`
	Expires                   string   `json:"expires"`
	KeyExpiryDisabled         bool     `json:"keyExpiryDisabled"`
	Authorized                bool     `json:"authorized"`
	IsExternal                bool     `json:"isExternal"`
	MachineKey                string   `json:"machineKey"`
	NodeKey                   string   `json:"nodeKey"`
	TailnetLockKey            string   `json:"tailnetLockKey"`
	BlocksIncomingConnections bool     `json:"blocksIncomingConnections"`
	TailnetLockError          string   `json:"tailnetLockError"`
}

type listDevicesResponse struct {
	Devices []*Device `json:"devices"`
}

type Client struct {
	baseURL    string
	tailnet    string
	clientID   string
	secret     string
	httpClient *http.Client
	mu         sync.Mutex
	token      string
	expiresAt  time.Time
}

func NewClient(tailnet, clientID, clientSecret string) *Client {
	return &Client{
		baseURL:  "https://api.tailscale.com",
		tailnet:  tailnet,
		clientID: clientID,
		secret:   clientSecret,
		httpClient: &http.Client{
			Timeout: 20 * time.Second,
		},
	}
}

func (c *Client) ListDevices(ctx context.Context) ([]*Device, error) {
	token, err := c.getToken(ctx)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/api/v2/tailnet/%s/devices", url.PathEscape(c.tailnet))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var body struct {
			Message string `json:"message"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&body)
		if body.Message != "" {
			return nil, fmt.Errorf("tailscale api error: %s", body.Message)
		}
		return nil, fmt.Errorf("tailscale api error: status %s", resp.Status)
	}
	var payload listDevicesResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	return payload.Devices, nil
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

func (c *Client) getToken(ctx context.Context) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.token != "" && time.Now().Before(c.expiresAt) {
		return c.token, nil
	}
	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", c.clientID)
	form.Set("client_secret", c.secret)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/v2/oauth/token", strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("oauth token request failed: status %s", resp.Status)
	}
	var payload tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", err
	}
	if payload.AccessToken == "" {
		return "", fmt.Errorf("oauth token response missing access_token")
	}
	if payload.ExpiresIn <= 0 {
		payload.ExpiresIn = 3600
	}
	payload.ExpiresIn = payload.ExpiresIn * 3 / 4
	c.token = payload.AccessToken
	c.expiresAt = time.Now().Add(time.Duration(payload.ExpiresIn) * time.Second)
	return c.token, nil
}
