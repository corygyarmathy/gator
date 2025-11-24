package main

import (
	"context"
	"fmt"
	"time"

	"github.com/corygyarmathy/gator/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %v <url>", cmd.Name)
	}

	url := cmd.Args[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("getting feed by URL: %v", err)
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("creating record in feed_follow table: %v", err)
	}

	fmt.Println("feed_follow record has been created in db.")
	printFeedFollow(feedFollow, user, feed)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %v", cmd.Name)
	}

	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("getting feed follows for user '%v': %v", user.Name, err)
	}

	if len(feedFollows) == 0 {
		fmt.Printf("No feed follows found for user '%v'\n", user.Name)
		return nil
	}

	fmt.Printf("Feed follows for user '%v':\n", user.Name)
	for _, feedFollow := range feedFollows {
		fmt.Printf(" * Feed(Name):  %v\n", feedFollow.FeedName)
	}

	return nil
}

func printFeedFollow(feedFollow database.CreateFeedFollowRow, user database.User, feed database.Feed) {
	fmt.Printf(" * ID:          %v\n", feedFollow.ID)
	fmt.Printf(" * User(Name):  %v\n", user.Name)
	fmt.Printf(" * Feed(Name):  %v\n", feed.Name)
}
