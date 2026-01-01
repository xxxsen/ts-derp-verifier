package tailscale

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Device struct {
	NodeKey    string `json:"nodeKey"`
	Authorized bool   `json:"authorized"`
}

type listDevicesResponse struct {
	Devices []*Device `json:"devices"`
}

type Client struct {
	baseURL    string
	tailnet    string
	apiKey     string
	httpClient *http.Client
}

func NewClient(tailnet, apiKey string) *Client {
	return &Client{
		baseURL: "https://api.tailscale.com",
		tailnet: tailnet,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 20 * time.Second,
		},
	}
}

func (c *Client) ListDevices(ctx context.Context) ([]*Device, error) {
	path := fmt.Sprintf("/api/v2/tailnet/%s/devices", url.PathEscape(c.tailnet))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
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
