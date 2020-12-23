package summary_test

import (
	"github.com/gh-utils/summary"
	"os"
	"testing"
	"time"
)

func TestGenerateRepoSummary(t *testing.T) {
	token := os.Getenv("GH_SUMMARY_GITHUB_PERSONAL_ACCESS_TOKEN")
	owner := "datastax"
	repo := "adelphi"
	startDate := time.Date(2020, 12, 15, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2020, 12, 18, 0, 0, 0, 0, time.UTC)

	summary, err := summary.GenerateRepoSummary(token, owner, repo, startDate, endDate)
	if err != nil {
		t.Errorf("GenerateRepoSummary failed - resulted in error: %v", err)
		return
	}
	if summary == nil {
		t.Errorf("GenerateRepoSummary failed - resulted in nil return")
		return
	}
	if summary.Repository == nil {
		t.Errorf("GenerateRepoSummary should produce a summary with Repository field")
	}
	if *summary.Repository.ID != 280185818 {
		t.Errorf(
			"GenerateRepoSummarys hould produce a a summary with Repository ID 280185818, it produced %v",
			*summary.Repository.ID,
		)
	}
	if summary.NewIssues == nil {
		t.Errorf("GenerateRepoSummary should produce a summary with NewIssues field")
	}
	if len(*summary.NewIssues) != 4 {
		t.Errorf("GenerateRepoSummary should produce a summary with 4 NewIssues, it produced %v", len(*summary.NewIssues))
	}
	if summary.ClosedIssues == nil {
		t.Errorf("GenerateRepoSummary should produce a summary with ClosedIssues field")
	}
	if len(*summary.ClosedIssues) != 5 {
		t.Errorf("GenerateRepoSummary should produce a summary with 5 ClosedIssues, it produced %v", len(*summary.ClosedIssues))
	}
}
