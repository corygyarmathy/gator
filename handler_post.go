package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/corygyarmathy/gator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	if len(cmd.Args) > 1 {
		return fmt.Errorf("usage: %v <limit>", cmd.Name)
	}

	var limit int32
	if len(cmd.Args) == 0 {
		limit = 2
	} else {
		parsedLimit, err := strconv.ParseInt(cmd.Args[0], 10, 32)
		if err != nil {
			return fmt.Errorf("converting string '%v' to int: %v", cmd.Args[0], err)
		}
		limit = int32(parsedLimit)
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	})
	if err != nil {
		return fmt.Errorf("getting posts for user '%v': %v", user.Name, err)
	}
	printPosts(posts)

	return nil
}

func printPosts(posts []database.Post) {
	for _, post := range posts {
		fmt.Printf(" * Title:       %v\n", post.Title)
	}
}
