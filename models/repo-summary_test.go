package models_test

import (
	"github.com/gh-utils/summary/models"
	"github.com/google/go-github/v33/github"
	"testing"
)

func TestNewRepoSummaryNilArg(t *testing.T) {
	_, summaryError := models.NewRepoSummary(nil)
	if summaryError == nil {
		t.Errorf("NewRepoSummary should return error when provided an invalid arguement")
	}
}

func TestNewRepoSummary(t *testing.T) {
	repository := github.Repository{}
	summary, summaryError := models.NewRepoSummary(&repository)
	if summaryError != nil {
		t.Errorf("NewRepoSummary should not return error when provided a valid arguement")
	}
	if summary == nil {
		t.Errorf("NewRepoSummary should return a summary when provided a valid arguement")
	}
}
