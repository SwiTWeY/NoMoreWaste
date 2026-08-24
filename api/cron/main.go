package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"

	"github.com/SwiTWeY/NoMoreWaste/api/config"
	"github.com/SwiTWeY/NoMoreWaste/api/rappel"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	db, err := config.Connect(cfg)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()

	n, err := rappel.EnvoyerRappels(context.Background(), db, cfg)
	if err != nil {
		log.Fatalf("rappels: %v", err)
	}
	log.Printf("rappels envoyes: %d", n)

}
