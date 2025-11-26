package github

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/shurcooL/githubv4"
)

type CreateCommitInput struct {
	Owner           string
	Repository      string
	Branch          string
	Message         string
	ExpectedHeadOid string
	Additions       []*FileAdditionInput
}

type FileAdditionInput struct {
	Path    string
	Content string
}

func (g *GitHub) CreateCommit(ctx context.Context, input *CreateCommitInput) error {
	parts := strings.SplitN(input.Message, "\n", 2)
	headline := parts[0]
	body := ""
	if len(parts) > 1 {
		body = parts[1]
	}

	additions := make([]githubv4.FileAddition, 0, len(input.Additions))
	for _, f := range input.Additions {
		enc := base64.StdEncoding.EncodeToString([]byte(f.Content))
		additions = append(additions, githubv4.FileAddition{
			Path:     githubv4.String(f.Path),
			Contents: githubv4.Base64String(enc),
		})
	}

	var m struct {
		CreateCommitOnBranch struct {
			Commit struct {
				URL string
			}
		} `graphql:"createCommitOnBranch(input:$input)"`
	}

	mutationInput := githubv4.CreateCommitOnBranchInput{
		Branch: githubv4.CommittableBranch{
			RepositoryNameWithOwner: githubv4.NewString(githubv4.String(input.Owner + "/" + input.Repository)),
			BranchName:              githubv4.NewString(githubv4.String(input.Branch)),
		},
		Message: githubv4.CommitMessage{
			Headline: githubv4.String(headline),
			Body:     githubv4.NewString(githubv4.String(body)),
		},
		FileChanges: &githubv4.FileChanges{
			Additions: &additions,
		},
		ExpectedHeadOid: githubv4.GitObjectID(input.ExpectedHeadOid),
	}

	if err := g.clientV4.Mutate(ctx, &m, mutationInput, nil); err != nil {
		return fmt.Errorf("failed to create commit: %w", err)
	}

	return nil
}
