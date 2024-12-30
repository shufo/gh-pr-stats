package utils

import (
	"time"

	"github.com/briandowns/spinner"
)

var (
	spin *spinner.Spinner
)

func StartSpinner(suffix string) {
	if !debug {
		spin = spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		spin.Suffix = suffix
		spin.Start()
	}
}

func StopSpinner() {
	if !debug {
		spin.Stop()
	}
}

func UpdateSpinnerSuffix(suffix string) {
	if !debug {
		spin.Suffix = suffix
	}
}
