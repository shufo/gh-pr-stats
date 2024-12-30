package utils

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/briandowns/spinner"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/shufo/gh-pr-stats/pkg/types"
	"github.com/spf13/cobra"
)

var header = []string{"Label", "Open", "Closed", "Total", "Open %", "Average Time to close (days)", "Median Time to close (days)"}

func PrintStatistics(cmd *cobra.Command, stats types.Statistics) {
	t := table.NewWriter()
	t.SetOutputMirror(cmd.OutOrStdout())
	t.SetStyle(table.StyleRounded)

	// Configure table style
	t.Style().Format.Header = text.FormatTitle
	t.Style().Options.DrawBorder = true
	t.Style().Options.SeparateHeader = true
	t.Style().Options.SeparateRows = false

	// Set header
	row := make(table.Row, len(header))
	for i, v := range header {
		row[i] = v
	}
	t.AppendHeader(row)

	// Add label statistics rows
	for _, stat := range stats.LabelStats {
		t.AppendRow(table.Row{
			stat.Name,
			stat.Open,
			stat.Closed,
			stat.Total,
			fmt.Sprintf("%.2f%%", stat.OpenPercentage),
			fmt.Sprintf("%.0f", stat.AvgDaysToClose),
			fmt.Sprintf("%.0f", stat.MedianDaysToClose),
		})
	}

	// Add separator and total row
	t.AppendSeparator()
	t.AppendRow(table.Row{
		"Total",
		stats.OverallStats.Open,
		stats.OverallStats.Closed,
		stats.OverallStats.Total,
		fmt.Sprintf("%.2f%%", stats.OverallStats.OpenPercentage),
		fmt.Sprintf("%.0f", stats.OverallStats.AvgDaysToClose),
		fmt.Sprintf("%.0f", stats.OverallStats.MedianDaysToClose),
	})

	// Render the table
	t.Render()
}

func SaveToFile(data interface{}, filename string) error {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = fmt.Sprintf(" Saving to %s...", filename)
	if !debug {
		s.Start()
	}

	StartSpinner(fmt.Sprintf(" Saving to %s...", filename))

	file, err := os.Create(filename)
	if err != nil {
		if !debug {
			StopSpinner()
		}
		return fmt.Errorf("failed to create file %s: %v", filename, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		if !debug {
			s.Stop()
		}
		return fmt.Errorf("failed to write to file %s: %v", filename, err)
	}

	if !debug {
		StopSpinner()
	}

	DebugPrintf("Data saved to %s", filename)

	return nil
}

func WriteDelimitedOutput(cmd *cobra.Command, stats types.Statistics, delimiter rune) error {
	writer := csv.NewWriter(cmd.OutOrStdout())
	writer.Comma = delimiter

	// Write header
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("error writing header: %v", err)
	}

	// Write label statistics
	for _, stat := range stats.LabelStats {
		row := []string{
			stat.Name,
			strconv.Itoa(stat.Open),
			strconv.Itoa(stat.Closed),
			strconv.Itoa(stat.Total),
			fmt.Sprintf("%.2f", stat.OpenPercentage),
			fmt.Sprintf("%.0f", stat.AvgDaysToClose),
			fmt.Sprintf("%.0f", stat.MedianDaysToClose),
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("error writing row: %v", err)
		}
	}

	// Write total row
	totalRow := []string{
		"Total",
		strconv.Itoa(stats.OverallStats.Open),
		strconv.Itoa(stats.OverallStats.Closed),
		strconv.Itoa(stats.OverallStats.Total),
		fmt.Sprintf("%.2f%%", stats.OverallStats.OpenPercentage),
		fmt.Sprintf("%.0f", stats.OverallStats.AvgDaysToClose),
		fmt.Sprintf("%.0f", stats.OverallStats.MedianDaysToClose),
	}
	if err := writer.Write(totalRow); err != nil {
		return fmt.Errorf("error writing total row: %v", err)
	}

	writer.Flush()
	return writer.Error()
}
