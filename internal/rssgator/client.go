package rssgator

import (
	"net/http"
	"time"

	"github.com/corygyarmathy/gator/internal/cachegator"
)

type Doer interface { // Enables mocking
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	httpClient Doer
	cache      *cachegator.Cache
}

func NewClient(timeout time.Duration, cache *cachegator.Cache) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: timeout},
		cache:      cache,
	}
}
