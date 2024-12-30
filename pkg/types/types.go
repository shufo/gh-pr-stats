package types

import (
	"time"
)

// PullRequest represents a GitHub pr with relevant fields
type PullRequest struct {
	Number      int        `json:"number"`
	Title       string     `json:"title"`
	State       string     `json:"state"`
	Labels      []Label    `json:"labels"`
	PullRequest *struct{}  `json:"pull_request,omitempty"`
	CreatedAt   *time.Time `json:"created_at"`
	ClosedAt    *time.Time `json:"closed_at"`
}

// Label represents a GitHub pr label
type Label struct {
	Name string `json:"name"`
}

const UnlabeledLabel = "*unlabeled*"

// LabelStat stores statistics for a specific label
type LabelStat struct {
	Name              string  `json:"name"`
	Open              int     `json:"open"`
	Closed            int     `json:"closed"`
	Total             int     `json:"total"`
	OpenPercentage    float64 `json:"openPercentage"`
	AvgDaysToClose    float64 `json:"AvgDaysToClose"`
	MedianDaysToClose float64 `json:"MedianDaysToClose"`
}

// OverallStats stores the overall pr statistics
type OverallStats struct {
	Total             int     `json:"total"`
	Open              int     `json:"open"`
	Closed            int     `json:"closed"`
	OpenPercentage    float64 `json:"openPercentage"`
	AvgDaysToClose    float64 `json:"AvgDaysToClose"`
	MedianDaysToClose float64 `json:"MedianDaysToClose"`
}

// Statistics combines both label and overall statistics
type Statistics struct {
	LabelStats   []LabelStat  `json:"labelStats"`
	OverallStats OverallStats `json:"overallStats"`
}
