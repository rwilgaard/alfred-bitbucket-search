package main

import (
    "fmt"

    bb "github.com/rwilgaard/bitbucket-go-api"
)

func getAllRepositories(api *bb.API) ([]*bb.RepositoryList, error) {
    query := bb.RepositoriesQuery{
        Limit: 1000,
    }

    repos, resp, err := api.GetRepositories(query)
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
            Start: int(repos.NextPageStart),
        }
        repos, resp, err = api.GetRepositories(query)
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

func getCommits(api *bb.API, projectKey string, repoSlug string) (*bb.CommitList, error) {
    query := bb.CommitsQuery{
        ProjectKey:     projectKey,
        RepositorySlug: repoSlug,
    }

    commits, resp, err := api.GetCommits(query)
    if err != nil {
        return nil, err
    }

    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("Failed to get commits. StatusCode: %d", resp.StatusCode)
    }

    return commits, nil
}

func getTags(api *bb.API, projectKey string, repoSlug string) (*bb.TagList, error) {
    query := bb.TagsQuery{
        ProjectKey:     projectKey,
        RepositorySlug: repoSlug,
        OrderBy:        "MODIFICATION",
    }

    tags, resp, err := api.GetTags(query)
    if err != nil {
        return nil, err
    }

    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("Failed to get tags. StatusCode: %d", resp.StatusCode)
    }

    return tags, nil
}

func getPullRequests(api *bb.API, projectKey string, repoSlug string) (*bb.PullRequestList, error) {
    query := bb.PullRequestsQuery{
        ProjectKey:     projectKey,
        RepositorySlug: repoSlug,
    }

    pr, resp, err := api.GetPullRequests(query)
    if err != nil {
        return nil, err
    }

    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("Failed to get pullrequests. StatusCode: %d", resp.StatusCode)
    }

    return pr, nil
}

func getBranches(api *bb.API, projectKey string, repoSlug string) (*bb.BranchList, error) {
    query := bb.BranchesQuery{
        ProjectKey:     projectKey,
        RepositorySlug: repoSlug,
    }

    br, resp, err := api.GetBranches(query)
    if err != nil {
        return nil, err
    }

    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("Failed to get branches. StatusCode: %d", resp.StatusCode)
    }

    return br, nil
}
