package gbfs

import (
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// Client represents abstraction on top of all GBFS operations
type Client struct {
	client    *http.Client
	userAgent string
}

// NewClient creates new client
func NewClient(userAgent string, timeout time.Duration) *Client {
	client := &http.Client{
		Timeout: timeout,
	}

	return &Client{
		client:    client,
		userAgent: userAgent,
	}
}

func (c *Client) sendRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "create new request")
	}

	req.Header.Add("User-Agent", c.userAgent)

	return c.client.Do(req)
}
