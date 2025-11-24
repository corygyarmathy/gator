package main

import (
	"context"
	"fmt"
	"time"

	"github.com/corygyarmathy/gator/internal/rssgator"
)

func handlerAggregator(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %v <time_between_reps>", cmd.Name)
	}

	timeArg := cmd.Args[0]
	timeBetweenRequests, err := time.ParseDuration(timeArg)
	if err != nil {
		return fmt.Errorf("parsing time arg '%v': %v", timeArg, err)
	}

	fmt.Printf("Collecting feeds every %v\n", timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			return fmt.Errorf("scraping feeds: %v", err)
		}
	}
}

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("getting next feed to fetch: %v", err)
	}
	feed, err := s.apiClient.FetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return fmt.Errorf("fetching feed: %v", err)
	}

	err = s.db.MarkFeedFetched(context.Background(), nextFeed.ID)
	if err != nil {
		return fmt.Errorf("marking feed as fetched: %v", err)
	}

	printFeedItems(feed)

	return nil
}

func printFeedItems(feed *rssgator.RSSFeed) {
	fmt.Printf("\n * Feed:        %v\n\n", feed.Channel.Title)
	for _, item := range feed.Channel.Item {
		fmt.Printf(" * Title:       %v\n", item.Title)
	}
}
