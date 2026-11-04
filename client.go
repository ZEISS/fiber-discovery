package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"reflect"
	"time"
)

// DefaultTimeout is the default timeout for the Perplexity API.
const DefaultTimeout = 30 * time.Second

// DefaultClient is the default HTTP client for the Perplexity API.
var DefaultClient = &http.Client{
	Timeout: DefaultTimeout,
}

// Client is an interface for discovering services.
type Client interface {
	// Discover discovers the well-known configuration for the given URL.
	Discover(ctx context.Context, url string, response any) error
}

type discoveryClient struct {
	client *http.Client
}

// Option is a function that configures the discovery client.
type Option func(*discoveryClient)

// WithClient returns an Option that sets the HTTP client for the discovery client.
func WithClient(client *http.Client) Option {
	return func(c *discoveryClient) {
		c.client = client
	}
}

// NewClient returns a new discovery client.
func NewClient(opts ...Option) Client {
	c := &discoveryClient{
		client: DefaultClient,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// Discover returns the well-known configuration for the given base URL.
func (d *discoveryClient) Discover(ctx context.Context, url string, value any) error {
	if value == nil {
		return fmt.Errorf("value must be a pointer")
	}

	if reflect.TypeOf(value).Kind() != reflect.Ptr {
		return fmt.Errorf("value must be a pointer to a struct")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("content-type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return err
	}

	switch mediaType {
	case "application/json":
		if err := json.NewDecoder(resp.Body).Decode(value); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unexpected content type: %s", contentType)
	}

	return nil
}
