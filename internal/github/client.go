package github

import (
	"fmt"

	"github.com/cli/go-gh/v2"
	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/shufo/gh-pr-stats/internal/utils"
	"github.com/shufo/gh-pr-stats/pkg/types"
)

var debug bool

func SetDebug(d bool) {
	debug = d
}

func GetRepoInfo() (string, error) {
	stdOut, stdErr, err := gh.Exec("repo", "view", "--json", "nameWithOwner", "-q", ".nameWithOwner")
	if err != nil {
		return "", fmt.Errorf("%v", stdErr.String())
	}
	return stdOut.String()[:stdOut.Len()-1], nil
}

// FetchPullRequestsFunc is a function type for fetching prs
type FetchPullRequestsFunc func(string) ([]types.PullRequest, error)

// DefaultFetchPullRequests is the actual implementation
func DefaultFetchPullRequests(repository string) ([]types.PullRequest, error) {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create GitHub client: %v", err)
	}

	if repository == "" {
		currentRepo, err := GetRepoInfo()
		if err != nil {
			return nil, fmt.Errorf("failed to get current repository: %w", err)
		}
		repository = currentRepo
	}

	// First, get the total count of prs to calculate pages
	var totalCount int
	// path := fmt.Sprintf("repos/%s/prs?state=all&per_page=1", repo)
	path := fmt.Sprintf("search/issues?q=repo:%s", repository)
	response := struct {
		TotalCount int `json:"total_count"`
	}{}

	err = client.Get(path, &response)
	if err == nil {
		totalCount = response.TotalCount
	}

	// Create and start spinner
	perPage := 100
	totalPages := (totalCount + perPage - 1) / perPage

	if totalPages > 0 {
		utils.DebugPrintf("Total prs (including PRs): %d", totalCount)
		utils.DebugPrintf("Total pages: %d\n", totalPages)
	}

	utils.StartSpinner(" Fetching prs...")

	var allPullRequests []types.PullRequest

	utils.DebugPrintf("starting to fetch pull requests")

	for page := 1; page <= totalPages; page++ {
		if debug {
			utils.DebugPrintf("fetching pull requests (%d/%d)", page, totalPages)
		} else {
			utils.UpdateSpinnerSuffix(fmt.Sprintf(" Fetching pull requests... (%d/%d)", page, totalPages))
		}

		var pagePullRequests []types.PullRequest
		path := fmt.Sprintf("repos/%s/issues?state=all&per_page=%d&page=%d", repository, perPage, page)
		err := client.Get(path, &pagePullRequests)
		if err != nil {
			utils.StopSpinner()
			return nil, fmt.Errorf("failed to fetch prs: %v", err)
		}

		if len(pagePullRequests) == 0 {
			break
		}

		// Filter out pull requests and count prs
		prsCount := 0
		for _, pr := range pagePullRequests {
			if pr.PullRequest != nil {
				allPullRequests = append(allPullRequests, pr)
				prsCount++
			}
		}
		utils.DebugPrintf("fetched %d: found %d prs (total so far: %d)",
			page, prsCount, len(allPullRequests))
	}

	// Stop spinner and clear the line
	if !debug {
		utils.StopSpinner()
	}

	utils.DebugPrintf("finished fetching prs (total: %d)", len(allPullRequests))
	return allPullRequests, nil
}

// fetchPullRequests is the package variable that can be swapped in tests
var fetchPullRequests FetchPullRequestsFunc = DefaultFetchPullRequests

// FetchPullRequests is the public function that uses the variable
func FetchPullRequests(repository string) ([]types.PullRequest, error) {
	return fetchPullRequests(repository)
}

// SetFetchPullRequestsFunc allows replacing the fetch function for testing
func SetFetchPullRequestsFunc(f FetchPullRequestsFunc) FetchPullRequestsFunc {
	old := fetchPullRequests
	fetchPullRequests = f
	return old
}
