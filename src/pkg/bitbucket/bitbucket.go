package bitbucket

import (
    "fmt"

    bb "github.com/rwilgaard/bitbucket-go-api"
)

type BitbucketService struct {
    API *bb.API
}

func NewBitbucketService(url string, username string, apiToken string) (*BitbucketService, error) {
    api, err := bb.NewAPI(url, username, apiToken)
    if err != nil {
        return nil, err
    }
    return &BitbucketService{API: api}, nil
}

func (b *BitbucketService) TestAuthentication() (statusCode int, err error) {
    _, resp, err := b.API.GetInboxPullRequestCount()
    if err != nil {
        return resp.StatusCode, err
    }
    return resp.StatusCode, nil
}

func (b *BitbucketService) GetAllRepositories() ([]*bb.RepositoryList, error) {
    query := bb.RepositoriesQuery{
        Limit: 1000,
    }

    repos, resp, err := b.API.GetRepositories(query)
    if err != nil {
        return nil, err
    }
    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("Failed to get repositories. StatusCode: %d", resp.StatusCode)
    }

    var results []*bb.RepositoryList
    results = append(results, repos)

    for !repos.IsLastPage {
        query := bb.RepositoriesQuery{
            Limit: 1000,
            Start: uint(repos.NextPageStart),
        }
        repos, resp, err = b.API.GetRepositories(query)
        if err != nil {
            return nil, err
        }
        if resp.StatusCode != 200 {
            return nil, fmt.Errorf("Failed to get repositories. StatusCode: %d", resp.StatusCode)
        }
        results = append(results, repos)
    }

    return results, nil
}

func (b *BitbucketService) GetAllProjects() (*bb.ProjectList, error) {
    query := bb.ProjectsQuery{
        Limit: 1000,
    }

    projects, resp, err := b.API.GetProjects(query)
    if err != nil {
        return nil, err
    }
    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("Failed to get projects. StatusCode: %d", resp.StatusCode)
    }

    return projects, nil
}

func (b *BitbucketService) GetCommits(projectKey string, repoSlug string) (*bb.CommitList, error) {
    query := bb.CommitsQuery{
        ProjectKey:     projectKey,
        RepositorySlug: repoSlug,
    }

    commits, resp, err := b.API.GetCommits(query)
    if err != nil {
        return nil, err
    }

    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("Failed to get commits. StatusCode: %d", resp.StatusCode)
    }

    return commits, nil
}

func (b *BitbucketService) GetTags(projectKey string, repoSlug string) (*bb.TagList, error) {
    query := bb.TagsQuery{
        ProjectKey:     projectKey,
        RepositorySlug: repoSlug,
        OrderBy:        "MODIFICATION",
    }

    tags, resp, err := b.API.GetTags(query)
    if err != nil {
        return nil, err
    }

    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("Failed to get tags. StatusCode: %d", resp.StatusCode)
    }

    return tags, nil
}

func (b *BitbucketService) GetPullRequests(projectKey string, repoSlug string) (*bb.PullRequestList, error) {
    query := bb.PullRequestsQuery{
        ProjectKey:     projectKey,
        RepositorySlug: repoSlug,
    }

    pr, resp, err := b.API.GetPullRequests(query)
    if err != nil {
        return nil, err
    }

    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("Failed to get pullrequests. StatusCode: %d", resp.StatusCode)
    }

    return pr, nil
}

func (b *BitbucketService) GetBranches(projectKey string, repoSlug string) (*bb.BranchList, error) {
    query := bb.BranchesQuery{
        ProjectKey:     projectKey,
        RepositorySlug: repoSlug,
    }

    br, resp, err := b.API.GetBranches(query)
    if err != nil {
        return nil, err
    }

    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("Failed to get branches. StatusCode: %d", resp.StatusCode)
    }

    return br, nil
}
