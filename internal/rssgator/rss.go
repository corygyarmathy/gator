// Package rssgator provides functions for interacting with RSS APIs.
package rssgator

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
)

func (c *Client) FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	data, err := c.doGET(feedURL, ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching RSS feed: %w", err)
	}

	var rssFeed RSSFeed
	if err := xml.Unmarshal(data, &rssFeed); err != nil {
		return nil, fmt.Errorf("unmarshaling XML RSS feed: %w", err)
	}

	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
	for i, item := range rssFeed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		rssFeed.Channel.Item[i] = item
	}

	return &rssFeed, nil
}
