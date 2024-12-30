package utils

import (
	"fmt"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

var logger *slog.Logger

func SetupLogger(cmd *cobra.Command, debug bool) {
	opts := &slog.HandlerOptions{}
	if debug {
		opts.Level = slog.LevelDebug
	} else {
		opts.Level = slog.LevelInfo
	}
	handler := slog.NewJSONHandler(cmd.OutOrStdout(), opts)
	logger = slog.New(handler)
}

func DebugPrintf(format string, a ...interface{}) {
	if debug {
		logger.Debug(fmt.Sprintf(format, a...))

	}
}
