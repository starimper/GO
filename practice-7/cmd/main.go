package main

import (
	"log"
	"practice-7/config"
	"practice-7/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("main: config: %v", err)
	}
	app.Run(cfg)
}
