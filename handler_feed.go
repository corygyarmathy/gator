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

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %v <name> <url>", cmd.Name)
	}

	name := cmd.Args[0]
	url := cmd.Args[1]

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
	printFeed(feed, user)

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" { // 23505: unique_violation
			return fmt.Errorf("feedFollow for '%v' already exists", name)
		}
		return fmt.Errorf("creating feedFollow in DB: %v", err)
	}
	fmt.Println("FeedFollow has been created in DB")
	printFeedFollow(feedFollow, user, feed)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %v", cmd.Name)
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("getting feeds table: %v", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}

	err = printFeeds(s, feeds)
	if err != nil {
		return fmt.Errorf("printing feeds: %v", err)
	}

	return nil
}

func printFeed(feed database.Feed, user database.User) {
	fmt.Printf(" * ID:          %v\n", feed.ID)
	fmt.Printf(" * Name:        %v\n", feed.Name)
	fmt.Printf(" * URL:         %v\n", feed.Url)
	fmt.Printf(" * UserID:      %v\n", feed.UserID)
	fmt.Printf(" * User(Name):  %v\n", user.Name)
}

func printFeeds(s *state, feeds []database.Feed) error {
	for _, feed := range feeds {
		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("getting current user: %v", err)
		}
		printFeed(feed, user)
		fmt.Println("=====================================")
	}
	return nil
}
