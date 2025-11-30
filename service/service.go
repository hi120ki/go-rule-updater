package service

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

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
	branch := fmt.Sprintf("add/%s", id)

	if err := s.gh.CreateBranch(ctx, s.cfg.Owner, s.cfg.Repository, branch, s.cfg.BaseBranch); err != nil {
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
		Branch:          branch,
		Message:         "Add new rule",
		Additions:       []*ghclient.FileAdditionInput{{Path: s.cfg.RulePath, Content: newContent}},
		ExpectedHeadOid: sha,
	}); err != nil {
		return nil, fmt.Errorf("failed to create commit: %w", err)
	}

	pr, err := s.gh.CreatePullRequest(ctx, s.cfg.Owner, s.cfg.Repository, "Add new rule", branch, s.cfg.BaseBranch, "This PR adds a new rule.")
	if err != nil {
		return nil, fmt.Errorf("failed to create pull request: %w", err)
	}

	if err := s.gh.CreatePullRequestComment(ctx, s.cfg.Owner, s.cfg.Repository, pr.GetNumber(), "Automated PR created to add a new rule."); err != nil {
		return nil, fmt.Errorf("failed to create pull request comment: %w", err)
	}

	return pr, nil
}

func (s *Service) Merge(ctx context.Context, prNumber int) error {
	retryDelay := time.Duration(s.cfg.MergeRetryDelaySeconds) * time.Second

	attempts := 0
	for {
		pr, err := s.gh.GetPullRequest(ctx, s.cfg.Owner, s.cfg.Repository, prNumber)
		if err != nil {
			return fmt.Errorf("failed to get pull request: %w", err)
		}

		if pr.GetMergeable() && pr.GetMergeableState() == "clean" {
			if err := s.gh.MergePullRequest(ctx, s.cfg.Owner, s.cfg.Repository, prNumber); err != nil {
				return fmt.Errorf("failed to merge pull request: %w", err)
			}
			return nil
		}

		attempts++
		if attempts >= s.cfg.MergeMaxRetries {
			return fmt.Errorf("PR #%d is not mergeable", prNumber)
		}

		log.Printf("PR #%d is not mergeable yet (mergeable: %v, state: %s)", prNumber, pr.GetMergeable(), pr.GetMergeableState())
		time.Sleep(retryDelay)
	}
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

		if s.gh.IsConflicting(fullPR) && strings.HasPrefix(fullPR.Head.GetRef(), "add/") {
			log.Printf("Resolving conflicts for PR #%d with rebase: %s", fullPR.GetNumber(), fullPR.GetTitle())
			if err := s.updateConflictingPR(ctx, fullPR); err != nil {
				log.Printf("Failed to resolve conflicts for PR #%d: %v", fullPR.GetNumber(), err)
				continue
			}
		}
	}

	return nil
}

func (s *Service) updateConflictingPR(ctx context.Context, input *github.PullRequest) error {
	branchName := input.Head.GetRef()

	// Get PR commits to find the original base
	commits, err := s.gh.ListPullRequestCommits(ctx, s.cfg.Owner, s.cfg.Repository, input.GetNumber())
	if err != nil {
		return fmt.Errorf("failed to get PR commits for PR: %w", err)
	}

	if len(commits) == 0 {
		return fmt.Errorf("no commits found for PR")
	}

	// Get the parent of the first commit (original base)
	firstCommit := commits[0]
	if len(firstCommit.Parents) == 0 {
		return fmt.Errorf("first commit has no parent for PR")
	}
	originalBaseSHA := firstCommit.Parents[0].GetSHA()

	// Get rule file content from original base
	originalBaseContent, err := s.gh.GetFile(ctx, s.cfg.Owner, s.cfg.Repository, s.cfg.RulePath, originalBaseSHA)
	if err != nil {
		return fmt.Errorf("failed to get original base rule file for PR: %w", err)
	}

	// Get rule file content from PR branch
	prContent, err := s.gh.GetFile(ctx, s.cfg.Owner, s.cfg.Repository, s.cfg.RulePath, branchName)
	if err != nil {
		return fmt.Errorf("failed to get PR rule file for PR: %w", err)
	}

	// Get added rules from the diff
	addedRules, err := rule.GetAddedRules(originalBaseContent, prContent)
	if err != nil {
		return fmt.Errorf("failed to get added rules for PR: %w", err)
	}

	if len(addedRules) == 0 {
		return fmt.Errorf("no added rules found for PR")
	}

	// Create temporary branch from base
	tempBranch := fmt.Sprintf("tmp/rebase/%d-%d", input.GetNumber(), time.Now().Unix())
	if err := s.gh.CreateBranch(ctx, s.cfg.Owner, s.cfg.Repository, tempBranch, s.cfg.BaseBranch); err != nil {
		return fmt.Errorf("failed to create temp branch for PR: %w", err)
	}

	// Cleanup temp branch on exit
	defer func(branch string) {
		if err := s.gh.DeleteBranch(ctx, s.cfg.Owner, s.cfg.Repository, branch); err != nil {
			log.Printf("Failed to delete temp branch %s: %v", branch, err)
		}
	}(tempBranch)

	// Get temp branch HEAD SHA (should match baseSHA, but verify for safety)
	tempBranchSHA, err := s.gh.GetLatestCommitSHA(ctx, s.cfg.Owner, s.cfg.Repository, tempBranch)
	if err != nil {
		return fmt.Errorf("failed to get temp branch SHA for PR: %w", err)
	}

	// Get current rule file content from base
	content, err := s.gh.GetFile(ctx, s.cfg.Owner, s.cfg.Repository, s.cfg.RulePath, s.cfg.BaseBranch)
	if err != nil {
		return fmt.Errorf("failed to get rule file for PR: %w", err)
	}

	// Re-apply all added rules
	for _, ruleID := range addedRules {
		content, err = rule.Add(content, ruleID)
		if err != nil {
			return fmt.Errorf("failed to re-apply rule %s for PR: %w", ruleID, err)
		}
	}

	// Create new commit on temp branch
	if err := s.gh.CreateCommit(ctx, &ghclient.CreateCommitInput{
		Owner:           s.cfg.Owner,
		Repository:      s.cfg.Repository,
		Branch:          tempBranch,
		Message:         "Rebase: Add new rule",
		Additions:       []*ghclient.FileAdditionInput{{Path: s.cfg.RulePath, Content: content}},
		ExpectedHeadOid: tempBranchSHA,
	}); err != nil {
		return fmt.Errorf("failed to create rebased commit on temp branch for PR: %w", err)
	}

	// Get new commit SHA from temp branch
	newSHA, err := s.gh.GetLatestCommitSHA(ctx, s.cfg.Owner, s.cfg.Repository, tempBranch)
	if err != nil {
		return fmt.Errorf("failed to get temp branch SHA for PR: %w", err)
	}

	// Atomically update PR branch to point to new commit (force push)
	if err := s.gh.UpdateBranchRef(ctx, s.cfg.Owner, s.cfg.Repository, branchName, newSHA, true); err != nil {
		return fmt.Errorf("failed to update PR branch for PR: %w", err)
	}

	return nil
}
