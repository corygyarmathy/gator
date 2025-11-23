package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/corygyarmathy/gator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
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

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %v <name> <url>", cmd.Name)
	}

	name := cmd.Args[0]
	url := cmd.Args[1]
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("getting current user: %v", err)
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" { // 23505: unique_violation
			return fmt.Errorf("feed '%v' already exists", name)
		}
		return fmt.Errorf("creating feed in DB: %v", err)
	}

	fmt.Println("Feed has been created in DB")
	printFeed(feed)

	return nil
}

func printFeed(feed database.Feed) {
	fmt.Printf(" * ID:      %v\n", feed.ID)
	fmt.Printf(" * Name:    %v\n", feed.Name)
	fmt.Printf(" * URL:     %v\n", feed.Url)
	fmt.Printf(" * UserID:  %v\n", feed.UserID)
}
