package main

import (
	"context"
	"fmt"
)

func handlerAggregator(s *state, cmd command) error {
	// if len(cmd.Args) != 1 {
	// return fmt.Errorf("usage: %v <url>", cmd.Name)
	// }

	// url := cmd.Args[0]
	url := "https://www.wagslane.dev/index.xml"
	feed, err := s.apiClient.FetchFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("fetching feed: %v", err)
	}

	fmt.Printf("Feed: %+v\n", feed)

	return nil
}
