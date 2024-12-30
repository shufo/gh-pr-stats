package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/shufo/gh-pr-stats/internal/github"
	"github.com/shufo/gh-pr-stats/internal/stats"
	"github.com/shufo/gh-pr-stats/internal/utils"
	"github.com/spf13/cobra"
)

var (
	outputFile string
	statsFile  string
	format     string
	debug      bool

	Version = "dev"
)

func Exec() {
	rootCmd := &cobra.Command{
		Use:   "gh pr-stats [repository]",
		Short: "Generate GitHub pr statistics",
		Long: `A GitHub CLI extension to analyze repository prs and generate statistics.
Provides detailed information about prs grouped by labels and overall statistics.

Examples:
  # Current repository
  gh pr-stats

  # Specific repository
  gh pr-stats owner/repo

  # With output format
  gh pr-stats owner/repo --format json`,
		Args:          cobra.MaximumNArgs(1),
		RunE:          runCommand,
		SilenceErrors: true,
		SilenceUsage:  true,
		Version:       getVersion(),
	}

	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file for raw prs data (optional)")
	rootCmd.Flags().StringVarP(&statsFile, "stats", "s", "", "Output file for statistics data (optional)")
	rootCmd.Flags().StringVarP(&format, "format", "f", "", "Output format: table (default), json, csv, or tsv")
	rootCmd.Flags().BoolVarP(&debug, "debug", "v", false, "Enable verbose debug output")

	// Customize version template
	rootCmd.SetVersionTemplate(`gh-pr-stats {{printf "version: %s" .Version}}
`)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runCommand(cmd *cobra.Command, args []string) error {
	utils.SetupLogger(cmd, debug)
	utils.SetDebug(debug)
	github.SetDebug(debug)

	var repository string
	if len(args) > 0 {
		repository = args[0]
		// Validate repository format
		if !isValidRepositoryFormat(repository) {
			return fmt.Errorf("invalid repository format. Expected format: owner/repo")
		}
	}
	// Fetch prs
	prs, err := github.FetchPullRequests(repository)
	if err != nil {
		return err
	}

	// Save prs if output file is specified
	if outputFile != "" {
		if err := utils.SaveToFile(prs, outputFile); err != nil {
			return err
		}
	}

	// Calculate statistics
	stats := stats.CalculateStatistics(prs)

	// Save statistics if stats file is specified
	if statsFile != "" {
		if err := utils.SaveToFile(stats, statsFile); err != nil {
			return err
		}
	}

	// Output based on format
	switch strings.ToLower(format) {
	case "json":
		encoder := json.NewEncoder(cmd.OutOrStdout())
		encoder.SetIndent("", "  ")
		return encoder.Encode(stats)
	case "csv":
		return utils.WriteDelimitedOutput(cmd, stats, ',')
	case "tsv":
		return utils.WriteDelimitedOutput(cmd, stats, '\t')
	default:
		utils.PrintStatistics(cmd, stats)
	}

	return nil
}

// isValidRepositoryFormat validates the repository argument format
func isValidRepositoryFormat(repo string) bool {
	parts := strings.Split(repo, "/")
	return len(parts) == 2 && parts[0] != "" && parts[1] != ""
}

// getVersion returns the version string
func getVersion() string {
	return Version
}
