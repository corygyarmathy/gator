package rssgator

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) doGET(url string, ctx context.Context) ([]byte, error) {
	data, ok := c.cache.Get(url)
	if !ok || data == nil { // URL not in cache or no cached data, retrieve via HTTP GET
		var req *http.Request
		var err error
		if ctx != nil {
			req, err = http.NewRequestWithContext(ctx, "GET", url, nil)
			if err != nil {
				return nil, fmt.Errorf("creating request: %w", err)
			}
		} else {
			req, err = http.NewRequest("GET", url, nil)
			if err != nil {
				return nil, fmt.Errorf("creating request: %w", err)
			}
		}

		req.Header.Set("User-Agent", "gator")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("making request: %w", err)
		}
		defer func() {
			if cerr := resp.Body.Close(); cerr != nil && err == nil {
				err = fmt.Errorf("closing response body: %w", cerr)
			}
		}()
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("reading data: %w", err)
		}

		c.cache.Add(url, data)
	}
	return data, nil
}
