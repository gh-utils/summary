package models

import (
	"errors"
	"github.com/google/go-github/v33/github"
)

type RepoSummary struct {
	Repository          *github.Repository      `json:"repository,omitempty"`
	NewIssues           *map[int64]github.Issue `json:"new_issues,omitempty"`
	UpdatedIssues       *map[int64]github.Issue `json:"updated_issues,omitempty"`
	ClosedIssues        *map[int64]github.Issue `json:"closed_issues,omitempty"`
	NewPullRequests     *map[int64]github.Issue `json:"new_prs,omitempty"`
	UpdatedPullRequests *map[int64]github.Issue `json:"updated_prs,omitempty"`
	ClosedPullRequests  *map[int64]github.Issue `json:"closed_prs,omitempty"`
}

func NewRepoSummary(repository *github.Repository) (summary *RepoSummary, err error) {
	if repository == nil {
		return nil, errors.New("repository provided is invalid")
	}

	summary = new(RepoSummary)
	summary.Repository = repository
	return summary, nil
}
