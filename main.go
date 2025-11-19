package main

import (
	"fmt"
	"log"

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
	if err != nil {
		log.Fatalf("Config file set user error: %v\n", err)
	}

	if err != nil {
		log.Fatalf("Config file read error: %v\n", err)
	}
	fmt.Printf("Read config: %+v\n", cfg)
}
