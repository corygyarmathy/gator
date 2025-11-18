package main

import (
	"fmt"
	"log"

	"github.com/corygyarmathy/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Config file read error: %v\n", err)
	}
	fmt.Printf("Read config: %+v\n", cfg)
	err = cfg.SetUser("coryg")
	if err != nil {
		log.Fatalf("Config file set user error: %v\n", err)
	}

	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("Config file read error: %v\n", err)
	}
	fmt.Printf("Read config: %+v\n", cfg)
}
