package main

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/hi120ki/go-rule-updater/env"
	ghclient "github.com/hi120ki/go-rule-updater/github"
	"github.com/hi120ki/go-rule-updater/service"
)

func main() {
	ctx := context.Background()

	cfg, err := env.Load()
	if err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}

	gh, err := ghclient.NewClient(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to initialize GitHub: %v", err)
	}

	svc := service.NewService(cfg, gh)

	pr, err := svc.Add(ctx, uuid.New().String())
	if err != nil {
		log.Fatalf("Failed to add new rule and create PR: %v", err)
	}

	time.Sleep(10 * time.Second)

	if err := svc.Merge(ctx, pr.GetNumber()); err != nil {
		log.Fatalf("Failed to merge pull request: %v", err)
	}

	time.Sleep(10 * time.Second)

	if err := svc.UpdatePRs(ctx); err != nil {
		log.Printf("Warning: Failed to update conflicting PRs: %v", err)
	}
}
