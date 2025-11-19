package main

import "fmt"

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
