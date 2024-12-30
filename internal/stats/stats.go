package stats

import (
	"sort"

	"github.com/shufo/gh-pr-stats/pkg/types"
)

// calculateMedian calculates the median value from a slice of float64
func calculateMedian(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	sort.Float64s(values)
	middle := len(values) / 2

	if len(values)%2 == 0 {
		return (values[middle-1] + values[middle]) / 2
	}
	return values[middle]
}

func CalculateStatistics(prs []types.PullRequest) types.Statistics {
	labelStatsSlice := make([]types.LabelStat, 0)
	labelStats := make(map[string]*types.LabelStat)
	overallStats := types.OverallStats{}

	// Calculate average close time per label
	labelAvgDaysToClose := make(map[string]float64)
	labelCloseTimes := make(map[string][]float64)

	var totalCloseTime float64
	var allCloseTimes []float64
	var closedPullRequestsCount int

	for _, pr := range prs {
		// Update overall stats
		overallStats.Total++
		if pr.State == "open" {
			overallStats.Open++
		} else {
			overallStats.Closed++
		}

		// Update label stats
		if len(pr.Labels) == 0 {
			// Handle unlabeled prs
			label := types.UnlabeledLabel
			stat, exists := labelStats[label]
			if !exists {
				stat = &types.LabelStat{Name: label}
				labelStats[label] = stat
			}
			stat.Total++
			if pr.State == "open" {
				stat.Open++
			} else {
				stat.Closed++
			}
		} else {
			// Handle labeled prs
			for _, label := range pr.Labels {
				stat, exists := labelStats[label.Name]
				if !exists {
					stat = &types.LabelStat{Name: label.Name}
					labelStats[label.Name] = stat
				}
				stat.Total++
				if pr.State == "open" {
					stat.Open++
				} else {
					stat.Closed++
				}
			}
		}

		if pr.State == "closed" {
			if pr.ClosedAt != nil && pr.CreatedAt != nil {
				closeTime := pr.ClosedAt.Sub(*pr.CreatedAt)
				if closeTime >= 0 {
					closeTimeDays := closeTime.Hours() / 24 // Convert to days
					totalCloseTime += closeTimeDays         // Store in days
					allCloseTimes = append(allCloseTimes, closeTimeDays)
					closedPullRequestsCount++
					if len(pr.Labels) == 0 {
						labelAvgDaysToClose[types.UnlabeledLabel] += closeTimeDays
						labelCloseTimes[types.UnlabeledLabel] = append(
							labelCloseTimes[types.UnlabeledLabel],
							closeTimeDays,
						)
					} else {
						for _, label := range pr.Labels {
							labelAvgDaysToClose[label.Name] += closeTimeDays
							labelCloseTimes[label.Name] = append(
								labelCloseTimes[label.Name],
								closeTimeDays,
							)
						}
					}
				}
			}
		}
	}

	// Convert map to sorted slice
	for _, stat := range labelStats {
		labelStatsSlice = append(labelStatsSlice, *stat)
	}
	sort.Slice(labelStatsSlice, func(i, j int) bool {
		return labelStatsSlice[i].Total > labelStatsSlice[j].Total
	})

	// Calculate the average close time for each label (already in days)
	for i, stat := range labelStatsSlice {
		totalLabelCloseTime := labelAvgDaysToClose[stat.Name]

		if stat.Total > 0 {
			stat.OpenPercentage = float64(stat.Open) / float64(stat.Total) * 100
		}

		if stat.Closed > 0 {
			stat.AvgDaysToClose = totalLabelCloseTime / float64(stat.Closed)
			stat.MedianDaysToClose = calculateMedian(labelCloseTimes[stat.Name])
		}

		labelStatsSlice[i] = stat
	}

	// Calculate the overall average close time (already in days)
	var overallAvgDaysToClose, overallMedianDaysToClose float64
	if closedPullRequestsCount > 0 {
		overallAvgDaysToClose = totalCloseTime / float64(closedPullRequestsCount)
		overallMedianDaysToClose = calculateMedian(allCloseTimes)
	}

	return types.Statistics{
		LabelStats: labelStatsSlice,
		OverallStats: types.OverallStats{
			Total:             overallStats.Total,
			Open:              overallStats.Open,
			OpenPercentage:    float64(overallStats.Open) / float64(overallStats.Total) * 100,
			Closed:            overallStats.Closed,
			AvgDaysToClose:    overallAvgDaysToClose,
			MedianDaysToClose: overallMedianDaysToClose,
		},
	}
}
