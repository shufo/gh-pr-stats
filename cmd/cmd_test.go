package cmd

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/shufo/gh-pr-stats/internal/github"
	"github.com/shufo/gh-pr-stats/pkg/types"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// Helper function to create test prs
func createTestPullRequests() []types.PullRequest {
	now := time.Now()
	createdAt := now.Add(-24 * time.Hour) // 1 day ago
	closedAt := now

	return []types.PullRequest{
		{
			Title:     "Test PullRequest 1",
			State:     "open",
			CreatedAt: &createdAt,
			Labels: []types.Label{
				{Name: "test_bug"},
			},
		},
		{
			Title:     "Test PullRequest 2",
			State:     "closed",
			CreatedAt: &createdAt,
			ClosedAt:  &closedAt,
			Labels: []types.Label{
				{Name: "test_enhancement"},
			},
		},
	}
}

// Helper to setup test command with mocked FetchPullRequests
func setupTestCommand() (*cobra.Command, *bytes.Buffer) {
	buf := new(bytes.Buffer)
	cmd := &cobra.Command{
		Use:          "gh pr-stats [repository]",
		RunE:         runCommand,
		SilenceUsage: true,
	}
	cmd.SetOutput(buf)

	// Add flags
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "")
	cmd.Flags().StringVarP(&statsFile, "stats", "s", "", "")
	cmd.Flags().StringVarP(&format, "format", "f", "", "")
	cmd.Flags().BoolVarP(&debug, "debug", "v", false, "")

	return cmd, buf
}

func TestRunCommandWithMock(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		format         string
		mockFetch      github.FetchPullRequestsFunc
		expectError    bool
		validateOutput func(*testing.T, []byte)
	}{
		{
			name:   "Successfully fetch prs with JSON output",
			args:   []string{"owner/repo"},
			format: "json",
			mockFetch: func(repo string) ([]types.PullRequest, error) {
				assert.Equal(t, "owner/repo", repo)
				return createTestPullRequests(), nil
			},
			validateOutput: func(t *testing.T, output []byte) {
				var stats types.Statistics
				err := json.Unmarshal(output, &stats)
				assert.NoError(t, err, "Failed to parse JSON output")

				// Validate overall stats
				assert.Equal(t, 2, stats.OverallStats.Total, "Total prs should be 2")
				assert.Equal(t, 1, stats.OverallStats.Open, "Open prs should be 1")
				assert.Equal(t, 1, stats.OverallStats.Closed, "Closed prs should be 1")
				assert.Equal(t, 50.0, stats.OverallStats.OpenPercentage, "Open percentage should be 50%")

				// Validate label stats
				assert.Equal(t, 2, len(stats.LabelStats), "Should have 2 labels")

				// Find and validate specific labels
				for _, labelStat := range stats.LabelStats {
					switch labelStat.Name {
					case "test_bug":
						assert.Equal(t, 1, labelStat.Open, "Bug label should have 1 open pr")
						assert.Equal(t, 0, labelStat.Closed, "Bug label should have 0 closed prs")
						assert.Equal(t, 100.0, labelStat.OpenPercentage, "Bug label should be 100% open")
					case "test_enhancement":
						assert.Equal(t, 0, labelStat.Open, "Enhancement label should have 0 open prs")
						assert.Equal(t, 1, labelStat.Closed, "Enhancement label should have 1 closed pr")
						assert.Equal(t, 0.0, labelStat.OpenPercentage, "Enhancement label should be 0% open")
					}
				}
			},
		},
		{
			name:   "Invalid repository format",
			args:   []string{"invalid-repo"},
			format: "json",
			mockFetch: func(repo string) ([]types.PullRequest, error) {
				t.Fatal("FetchPullRequests should not be called")
				return nil, nil
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalFetch := github.SetFetchPullRequestsFunc(tt.mockFetch)
			defer github.SetFetchPullRequestsFunc(originalFetch)

			cmd, buf := setupTestCommand()
			format = tt.format

			cmd.SetArgs(tt.args)
			err := cmd.Execute()

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			if tt.validateOutput != nil {
				tt.validateOutput(t, buf.Bytes())
			}
		})
	}
}
