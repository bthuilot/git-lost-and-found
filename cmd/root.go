package cmd

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var (
	// logLevel is the log level for the application
	logLevel string
	// logFormat is the log format for the application
	logFormat string
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
		level, err := logrus.ParseLevel(logLevel)
		if err != nil {
			return err
		}
		logrus.SetLevel(level)
		switch logFormat {
		case "text":
			logrus.SetFormatter(&logrus.TextFormatter{})
		case "json":
			logrus.SetFormatter(&logrus.JSONFormatter{})
		default:
			return fmt.Errorf("invalid log format: %s", logFormat)
		}

		return nil
	},
}

func init() {
	rootCmd.SetErrPrefix("ERROR: ")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", "log level (debug, info, warn, error, fatal, panic)")
	rootCmd.PersistentFlags().StringVar(&logFormat, "log-format", "text", "log format (text, json)")
}
