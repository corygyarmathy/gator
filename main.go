package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/corygyarmathy/gator/internal/cachegator"
	"github.com/corygyarmathy/gator/internal/config"
	"github.com/corygyarmathy/gator/internal/database"
	"github.com/corygyarmathy/gator/internal/rssgator"

	_ "github.com/lib/pq"
)

type state struct {
	cfg       *config.Config
	db        *database.Queries
	apiClient *rssgator.Client
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Config file read error: %v\n", err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("DB SQL open error: %v\n", err)
	}
	defer func() {
		if cerr := db.Close(); cerr != nil && err == nil {
			log.Fatalf("DB error: %v", cerr)
		}
	}()
	dbQueries := database.New(db)

	c := cachegator.NewCache(5 * time.Minute)
	defer c.Close()

	client := rssgator.NewClient(5*time.Second, c)

	pgrmState := &state{cfg: &cfg, db: dbQueries, apiClient: client}
	cmds := commands{cliCommands: map[string]func(*state, command) error{}}

	err = cmds.register("login", handlerLogin)
	if err != nil {
		log.Fatalf("Failed to register command: %v\n", err)
	}
	err = cmds.register("register", handlerRegister)
	if err != nil {
		log.Fatalf("Failed to register command: %v\n", err)
	}
	err = cmds.register("reset", handlerReset)
	if err != nil {
		log.Fatalf("Failed to register command: %v\n", err)
	}
	err = cmds.register("users", handlerListUsers)
	if err != nil {
		log.Fatalf("Failed to register command: %v\n", err)
	}
	err = cmds.register("agg", handlerAggregator)
	if err != nil {
		log.Fatalf("Failed to register command: %v\n", err)
	}
	err = cmds.register("addfeed", handlerAddFeed)
	if err != nil {
		log.Fatalf("Failed to register command: %v\n", err)
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
