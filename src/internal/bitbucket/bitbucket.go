package bitbucket

import (
	"fmt"

	bb "github.com/rwilgaard/bitbucket-go-api"
)

type Client struct {
	API *bb.API
}

func NewClient(url, username, token string) (*Client, error) {
	api, err := bb.NewAPI(url, username, token)
	if err != nil {
		return nil, fmt.Errorf("failed to create API client: %w", err)
	}
	return &Client{API: api}, nil
}

func (c *Client) TestAuthentication() error {
	_, resp, err := c.API.GetInboxPullRequestCount()
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) GetAllRepositories() ([]*bb.RepositoryList, error) {
	query := bb.RepositoriesQuery{Limit: 1000}
	repos, resp, err := c.API.GetRepositories(query)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get repositories. StatusCode: %d", resp.StatusCode)
	}
	var results []*bb.RepositoryList
	results = append(results, repos)
	for !repos.IsLastPage {
		query.Start = uint(repos.NextPageStart)
		repos, resp, err = c.API.GetRepositories(query)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("failed to get repositories. StatusCode: %d", resp.StatusCode)
		}
		results = append(results, repos)
	}
	return results, nil
}

func (c *Client) GetAllProjects() (*bb.ProjectList, error) {
	query := bb.ProjectsQuery{Limit: 1000}
	projects, resp, err := c.API.GetProjects(query)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get projects. StatusCode: %d", resp.StatusCode)
	}
	result := projects
	for !projects.IsLastPage {
		query.Start = projects.NextPageStart
		projects, resp, err = c.API.GetProjects(query)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("failed to get projects. StatusCode: %d", resp.StatusCode)
		}
		result.Values = append(result.Values, projects.Values...)
	}
	result.IsLastPage = true
	return result, nil
}

func (c *Client) GetCommits(projectKey, repoSlug string) (*bb.CommitList, error) {
	query := bb.CommitsQuery{ProjectKey: projectKey, RepositorySlug: repoSlug}
	commits, resp, err := c.API.GetCommits(query)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get commits. StatusCode: %d", resp.StatusCode)
	}
	return commits, nil
}

func (c *Client) GetTags(projectKey, repoSlug string) (*bb.TagList, error) {
	query := bb.TagsQuery{ProjectKey: projectKey, RepositorySlug: repoSlug, OrderBy: "MODIFICATION"}
	tags, resp, err := c.API.GetTags(query)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get tags. StatusCode: %d", resp.StatusCode)
	}
	return tags, nil
}

func (c *Client) GetBranches(projectKey, repoSlug string) (*bb.BranchList, error) {
	query := bb.BranchesQuery{ProjectKey: projectKey, RepositorySlug: repoSlug}
	br, resp, err := c.API.GetBranches(query)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get branches. StatusCode: %d", resp.StatusCode)
	}
	return br, nil
}

func (c *Client) GetPullRequests(projectKey, repoSlug string) (*bb.PullRequestList, error) {
	query := bb.PullRequestsQuery{ProjectKey: projectKey, RepositorySlug: repoSlug}
	pr, resp, err := c.API.GetPullRequests(query)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get pull requests. StatusCode: %d", resp.StatusCode)
	}
	return pr, nil
}
