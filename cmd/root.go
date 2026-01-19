// Copyright (C) 2024-2026 Bryce Thuilot
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the FSF, either version 3 of the License, or (at your option) any later version.
// See the LICENSE file in the root of this repository for full license text or
// visit: <https://www.gnu.org/licenses/gpl-3.0.html>.

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	// logLevel is the log level for the application
	logLevel string
	// logFormat is the log format for the application
	logFormat string

	version   string
	buildDate string
	gitCommit string
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "git-lost-and-found",
	Short: "git-lost-and-found will find all dangling commits in a git repository and create refs for them.",
	Long: `git-lost-and-found will find all dangling commits in a git repository and create refs for them.
This allows for scanners that use 'git log' to search blob data to not miss any changes.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		config := zap.NewProductionConfig()

		encoding := strings.ToLower(logFormat)
		switch encoding {
		case "json", "console":
			config.Encoding = encoding
		default:
			return fmt.Errorf("unsupported log format: %s", logFormat)
		}

		// Set the log level
		var level zap.AtomicLevel
		if err := level.UnmarshalText([]byte(logLevel)); err != nil {
			return fmt.Errorf("invalid log level: %s", logLevel)
		}
		config.Level = level

		logger, err := config.Build()
		if err != nil {
			return fmt.Errorf("failed to build zap logger: %w", err)
		}
		baseFields := []zap.Field{
			zap.String("version", version),
			zap.String("buildDate", buildDate),
			zap.String("gitCommit", gitCommit),
		}

		zap.ReplaceGlobals(logger.With(baseFields...))
		return nil
	},
}

func init() {
	rootCmd.SetErrPrefix("ERROR: ")
	rootCmd.PersistentFlags().
		StringVarP(&logLevel, "log-level", "l", "info", "log level (debug, info, warn, error, fatal, panic)")
	rootCmd.PersistentFlags().
		StringVar(&logFormat, "log-format", "console", "log format (console, json)")
}
