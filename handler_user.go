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

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("login args empty")
	}

	if len(cmd.Args) > 1 {
		return fmt.Errorf("login expects 1 arg, received more")
	}

	s.cfg.CurrentUserName = cmd.Args[0]
	fmt.Printf("Username has been set: %v", s.cfg.CurrentUserName)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %v <name>", cmd.Name)
	}
	name := cmd.Args[0]

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
	})
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" { // 23505: unique_violation
			return fmt.Errorf("user '%v' already exists", name)
		}
		return fmt.Errorf("creating user in DB: %v", err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("setting username: %v", err)
	}
	fmt.Println("User has been created in DB")
	printUser(user)

	return nil
}

func printUser(user database.User) {
	fmt.Printf(" * ID:      %v\n", user.ID)
	fmt.Printf(" * Name:    %v\n", user.Name)
}
