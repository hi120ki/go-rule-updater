package service

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/v79/github"
	"github.com/hi120ki/go-rule-updater/env"
	ghclient "github.com/hi120ki/go-rule-updater/github"
	"github.com/hi120ki/go-rule-updater/rule"
)

type Service struct {
	cfg *env.Env
	gh  *ghclient.GitHub
}

func NewService(cfg *env.Env, gh *ghclient.GitHub) *Service {
	return &Service{
		cfg: cfg,
		gh:  gh,
	}
}

func (s *Service) Add(ctx context.Context, id string) (*github.PullRequest, error) {
	if err := s.gh.CreateBranch(ctx, s.cfg.Owner, s.cfg.Repository, id, s.cfg.BaseBranch); err != nil {
		return nil, fmt.Errorf("failed to create branch: %w", err)
	}

	content, err := s.gh.GetFile(ctx, s.cfg.Owner, s.cfg.Repository, s.cfg.RulePath, s.cfg.BaseBranch)
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	newContent, err := rule.Add(content, id)
	if err != nil {
		return nil, fmt.Errorf("failed to add rule: %w", err)
	}

	sha, err := s.gh.GetLatestCommitSHA(ctx, s.cfg.Owner, s.cfg.Repository, s.cfg.BaseBranch)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest commit SHA: %w", err)
	}

	if err := s.gh.CreateCommit(ctx, &ghclient.CreateCommitInput{
		Owner:           s.cfg.Owner,
		Repository:      s.cfg.Repository,
		Branch:          id,
		Message:         "Add new rule",
		Additions:       []*ghclient.FileAdditionInput{{Path: s.cfg.RulePath, Content: newContent}},
		ExpectedHeadOid: sha,
	}); err != nil {
		return nil, fmt.Errorf("failed to create commit: %w", err)
	}

	pr, err := s.gh.CreatePullRequest(ctx, s.cfg.Owner, s.cfg.Repository, "Add new rule", id, s.cfg.BaseBranch, "This PR adds a new rule.")
	if err != nil {
		return nil, fmt.Errorf("failed to create pull request: %w", err)
	}

	if err := s.gh.CreatePullRequestComment(ctx, s.cfg.Owner, s.cfg.Repository, pr.GetNumber(), "Automated PR created to add a new rule."); err != nil {
		return nil, fmt.Errorf("failed to create pull request comment: %w", err)
	}

	return pr, nil
}

func (s *Service) Merge(ctx context.Context, prNumber int) error {
	if err := s.gh.MergePullRequest(ctx, s.cfg.Owner, s.cfg.Repository, prNumber); err != nil {
		return fmt.Errorf("failed to merge pull request: %w", err)
	}
	return nil
}

func (s *Service) UpdateConflictingPRs(ctx context.Context) error {
	prs, err := s.gh.ListOpenPullRequests(ctx, s.cfg.Owner, s.cfg.Repository)
	if err != nil {
		return err
	}

	for _, pr := range prs {
		fullPR, err := s.gh.GetPullRequest(ctx, s.cfg.Owner, s.cfg.Repository, pr.GetNumber())
		if err != nil {
			log.Printf("Failed to get PR #%d: %v", pr.GetNumber(), err)
			continue
		}

		if s.gh.IsConflicting(fullPR) {
			log.Printf("Updating conflicting PR #%d: %s", fullPR.GetNumber(), fullPR.GetTitle())
			if err := s.gh.UpdatePullRequestBranch(ctx, s.cfg.Owner, s.cfg.Repository, fullPR.GetNumber()); err != nil {
				log.Printf("Failed to update PR #%d: %v", fullPR.GetNumber(), err)
				continue
			}
			log.Printf("Successfully updated PR #%d", fullPR.GetNumber())
		}
	}

	return nil
}
