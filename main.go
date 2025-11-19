package main

import (
	"log"
	"os"

	"github.com/corygyarmathy/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Config file read error: %v\n", err)
	}
	pgrmState := &state{cfg: &cfg}
	cmds := commands{cliCommands: map[string]func(*state, command) error{}}

	if err != nil {
		log.Fatalf("Config file set user error: %v\n", err)
	}

	args := os.Args
	if len(args) < 2 {
		log.Fatalf("Too few args. Usage: cli <command> [args...]")
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]
	cmd := command{Name: cmdName, Args: cmdArgs}
	err = cmds.run(pgrmState, cmd)
	if err != nil {
		log.Fatalf("Command run error: %v\n", err)
	}
}
