package main

import "fmt"

type command struct {
	Name string
	Args []string
}

type commands struct {
	cliCommands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	cliCommand, exists := c.cliCommands[cmd.Name]
	if !exists {
		return fmt.Errorf("%v is not a registered command", cmd.Name)
	}

	if err := cliCommand(s, cmd); err != nil {
		return fmt.Errorf("executing command %v: %v", cmd.Name, err)
	}

	return nil
}
func (c *commands) register(name string, f func(*state, command) error) error {
	c.cliCommands[name] = f
	return nil
}
