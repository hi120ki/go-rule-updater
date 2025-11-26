package main

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/hi120ki/go-rule-updater/env"
	ghclient "github.com/hi120ki/go-rule-updater/github"
	"github.com/hi120ki/go-rule-updater/rule"
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

	new := uuid.New().String()

	if err := gh.CreateBranch(ctx, cfg.Owner, cfg.Repository, new, cfg.BaseBranch); err != nil {
		log.Fatalf("Failed to create branch: %v", err)
	}

	content, err := gh.GetFile(ctx, cfg.Owner, cfg.Repository, cfg.RulePath, cfg.BaseBranch)
	if err != nil {
		log.Fatalf("Failed to get file: %v", err)
	}

	newContent, err := rule.Add(content, new)
	if err != nil {
		log.Fatalf("Failed to add rule: %v", err)
	}

	sha, err := gh.GetLatestCommitSHA(ctx, cfg.Owner, cfg.Repository, cfg.BaseBranch)
	if err != nil {
		log.Fatalf("Failed to get latest commit SHA: %v", err)
	}

	if err := gh.CreateCommit(ctx, &ghclient.CreateCommitInput{
		Owner:           cfg.Owner,
		Repository:      cfg.Repository,
		Branch:          new,
		Message:         "Add new rule",
		Additions:       []*ghclient.FileAdditionInput{{Path: cfg.RulePath, Content: newContent}},
		ExpectedHeadOid: sha,
	}); err != nil {
		log.Fatalf("Failed to create commit: %v", err)
	}

	pr, err := gh.CreatePullRequest(ctx, cfg.Owner, cfg.Repository, "Add new rule", new, cfg.BaseBranch, "This PR adds a new rule.")
	if err != nil {
		log.Fatalf("Failed to create pull request: %v", err)
	}

	if err := gh.CreatePullRequestComment(ctx, cfg.Owner, cfg.Repository, pr.GetNumber(), "Automated PR created to add a new rule."); err != nil {
		log.Fatalf("Failed to create pull request comment: %v", err)
	}

	if err := gh.MergePullRequest(ctx, cfg.Owner, cfg.Repository, pr.GetNumber()); err != nil {
		log.Fatalf("Failed to merge pull request: %v", err)
	}
}
