/*
Copyright © 2025

Licensed under the MIT License.
*/

package cmd

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	gracefulshutdown "github.com/daffahilmyf/personal-go-template/internal/graceful_shutdown"
)

// gracefulShutdownCmd represents the gracefulShutdown command
var gracefulShutdownCmd = &cobra.Command{
	Use:   "graceful_shutdown_example",
	Short: "Demonstrates a graceful shutdown using ShutdownManager",
	Long: `Demonstrates a production-grade graceful shutdown pattern.
This example registers multiple shutdown steps and listens for
termination signals (e.g., SIGINT, SIGTERM). Upon receiving a signal,
it executes all shutdown steps concurrently within a configured timeout.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

		log.Info().Msg("starting graceful shutdown demo. Press Ctrl+C to trigger shutdown.")

		shutdown := gracefulshutdown.NewShutdownManager(log.Logger, 10*time.Second)

		// Register shutdown steps (example only)
		shutdown.AddStep("cleanup-cache", func(ctx context.Context) error {
			log.Info().Msg("cleaning up cache...")
			time.Sleep(1 * time.Second)
			log.Info().Msg("cache cleaned.")
			return nil
		})

		shutdown.AddStep("flush-metrics", func(ctx context.Context) error {
			log.Info().Msg("flushing metrics...")
			time.Sleep(2 * time.Second)
			log.Info().Msg("metrics flushed.")
			return nil
		})

		shutdown.AddStep("close-temp-files", func(ctx context.Context) error {
			log.Info().Msg("closing temporary files...")
			time.Sleep(1 * time.Second)
			log.Info().Msg("temporary files closed.")
			return nil
		})

		<-shutdown.WaitForSignal(context.Background())

		log.Info().Msg("graceful shutdown complete.")
	},
}

func init() {
	rootCmd.AddCommand(gracefulShutdownCmd)
}
