package summary

import (
	"context"
	"errors"
	"fmt"
	"github.com/gh-utils/summary/models"
	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
	"time"
)

func getNewIssuesInRepo(
	client *github.Client,
	context context.Context,
	token string,
	repoOwner string,
	repoName string,
	startDateInc time.Time,
	endDateExc time.Time,
) (issues map[int64]github.Issue, err error) {
	return getIssuesOrPullRequestsInRepo(
		client,
		context,
		token,
		repoOwner,
		repoName,
		false,
		"created",
		startDateInc,
		endDateExc,
	)
}

func getClosedIssuesInRepo(
	client *github.Client,
	context context.Context,
	token string,
	repoOwner string,
	repoName string,
	startDateInc time.Time,
	endDateExc time.Time,
) (issues map[int64]github.Issue, err error) {
	return getIssuesOrPullRequestsInRepo(
		client,
		context,
		token,
		repoOwner,
		repoName,
		false,
		"closed",
		startDateInc,
		endDateExc,
	)
}

func getUpdatedIssuesInRepo(
	client *github.Client,
	context context.Context,
	token string,
	repoOwner string,
	repoName string,
	startDateInc time.Time,
	endDateExc time.Time,
) (issues map[int64]github.Issue, err error) {
	return getIssuesOrPullRequestsInRepo(
		client,
		context,
		token,
		repoOwner,
		repoName,
		false,
		"updated",
		startDateInc,
		endDateExc,
	)
}

func getNewPullRequestsInRepo(
	client *github.Client,
	context context.Context,
	token string,
	repoOwner string,
	repoName string,
	startDateInc time.Time,
	endDateExc time.Time,
) (issues map[int64]github.Issue, err error) {
	return getIssuesOrPullRequestsInRepo(
		client,
		context,
		token,
		repoOwner,
		repoName,
		true,
		"created",
		startDateInc,
		endDateExc,
	)
}

func getClosedPullRequestsInRepo(
	client *github.Client,
	context context.Context,
	token string,
	repoOwner string,
	repoName string,
	startDateInc time.Time,
	endDateExc time.Time,
) (issues map[int64]github.Issue, err error) {
	return getIssuesOrPullRequestsInRepo(
		client,
		context,
		token,
		repoOwner,
		repoName,
		true,
		"closed",
		startDateInc,
		endDateExc,
	)
}

func getUpdatedPullRequestsInRepo(
	client *github.Client,
	context context.Context,
	token string,
	repoOwner string,
	repoName string,
	startDateInc time.Time,
	endDateExc time.Time,
) (issues map[int64]github.Issue, err error) {
	return getIssuesOrPullRequestsInRepo(
		client,
		context,
		token,
		repoOwner,
		repoName,
		true,
		"updated",
		startDateInc,
		endDateExc,
	)
}

func getIssuesOrPullRequestsInRepo(
	client *github.Client,
	context context.Context,
	token string,
	repoOwner string,
	repoName string,
	prs bool,
	dateContext string,
	startDateInc time.Time,
	endDateExc time.Time,
) (issues map[int64]github.Issue, err error) {
	if len(token) < 1 {
		return nil, errors.New("token provided is invalid")
	}
	if len(repoOwner) < 1 {
		return nil, errors.New("owner provided is invalid")
	}
	if len(repoName) < 1 {
		return nil, errors.New("repo provided is invalid")
	}

	// We want a full date range like
	// 2020-12-15..2020-12-16
	// <target date>..<target date +1>
	// This is inclusive..exclusive
	formattedStartDate := startDateInc.Format("2006-01-02")
	formattedEndDate := endDateExc.Format("2006-01-02")

	// Build a query
	issueTypeKeyword := "issue"
	if prs == true {
		issueTypeKeyword = "pull-request"
	}
	query := fmt.Sprintf("repo:%v/%v is:%v %v:%v..%v", repoOwner, repoName, issueTypeKeyword, dateContext, formattedStartDate, formattedEndDate)

	issuesSearchResult, response, err := client.Search.Issues(context, query, nil)
	if err != nil {
		return nil, err
	}
	if response == nil || response.Response.StatusCode != 200 {
		return nil, errors.New("failed response status code received from GitHub API: " + response.Response.Status)
	}

	issues = make(map[int64]github.Issue)
	for _, issue := range issuesSearchResult.Issues {
		if issue == nil || issue.ID == nil {
			continue
		}
		issues[*issue.ID] = *issue
	}
	return issues, nil
}

func GenerateRepoSummary(
	token string,
	repoOwner string,
	repoName string,
	startDateInc time.Time,
	endDateExc time.Time,
) (repoSummary *models.RepoSummary, e error) {
	if len(token) < 1 {
		return nil, errors.New("token provided is invalid")
	}
	if len(repoOwner) < 1 {
		return nil, errors.New("repoOwner provided is invalid")
	}
	if len(repoName) < 1 {
		return nil, errors.New("repoName provided is invalid")
	}

	// Build a client using the token provided
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Use the client to get the requested repository
	repository, response, err := client.Repositories.Get(ctx, repoOwner, repoName)
	// Check the response for validity
	if err != nil {
		return nil, err
	}
	if response == nil || response.Response.StatusCode != 200 {
		return nil, errors.New("failed response status code received from GitHub API: " + response.Response.Status)
	}

	summary, summaryError := models.NewRepoSummary(repository)
	if summaryError != nil {
		return nil, summaryError
	}

	newIssues, newIssuesError := getNewIssuesInRepo(client, ctx, token, repoOwner, repoName, startDateInc, endDateExc)
	if newIssuesError == nil {
		summary.NewIssues = &newIssues
	}
	updatedIssues, updatedIssuesError := getUpdatedIssuesInRepo(client, ctx, token, repoOwner, repoName, startDateInc, endDateExc)
	if updatedIssuesError == nil {
		summary.UpdatedIssues = &updatedIssues
	}
	closedIssues, closedIssuesError := getClosedIssuesInRepo(client, ctx, token, repoOwner, repoName, startDateInc, endDateExc)
	if closedIssuesError == nil {
		summary.ClosedIssues = &closedIssues
	}
	newPullRequests, newPullRequestsError := getNewPullRequestsInRepo(client, ctx, token, repoOwner, repoName, startDateInc, endDateExc)
	if newPullRequestsError == nil {
		summary.NewPullRequests = &newPullRequests
	}
	updatedPullRequests, updatedPullRequestsError := getUpdatedPullRequestsInRepo(client, ctx, token, repoOwner, repoName, startDateInc, endDateExc)
	if updatedPullRequestsError == nil {
		summary.UpdatedPullRequests = &updatedPullRequests
	}
	closedPullRequests, closedPullRequestsError := getClosedPullRequestsInRepo(client, ctx, token, repoOwner, repoName, startDateInc, endDateExc)
	if closedPullRequestsError == nil {
		summary.ClosedPullRequests = &closedPullRequests
	}

	return summary, nil
}
