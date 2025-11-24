package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/corygyarmathy/gator/internal/database"
	"github.com/corygyarmathy/gator/internal/rssgator"
	"github.com/google/uuid"
	"github.com/lib/pq"
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

	fmt.Printf("\nSaving posts from: %v\n\n", feed.Channel.Title)
	for _, item := range feed.Channel.Item {
		pubDate, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			return fmt.Errorf("parsing post '%v' date '%v': %v", item.Title, item.PubDate, err)
		}
		post, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: pubDate,
			FeedID:      nextFeed.ID,
		})
		if err != nil {
			var pqErr *pq.Error
			if errors.As(err, &pqErr) && pqErr.Code == "23505" { // 23505: unique_violation
				continue
			}
			return fmt.Errorf("creating user in DB: %v", err)
		}
		fmt.Printf("Stored post to db: %v \n", post.Title)
	}

	return nil
}

func printFeedItems(feed *rssgator.RSSFeed) {
	fmt.Printf("\n * Feed:        %v\n\n", feed.Channel.Title)
	for _, item := range feed.Channel.Item {
		fmt.Printf(" * Title:       %v\n", item.Title)
	}
}
