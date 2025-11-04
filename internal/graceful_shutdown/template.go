package gracefulshutdown

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/rs/zerolog"
)

type ShutdownStep struct {
	Name string
	Func func(ctx context.Context) error
}

type ShutdownManager struct {
	Logger  zerolog.Logger
	Timeout time.Duration
	Steps   []ShutdownStep
}

func NewShutdownManager(logger zerolog.Logger, timeout time.Duration) *ShutdownManager {
	return &ShutdownManager{
		Logger:  logger,
		Timeout: timeout,
	}
}

func (sm *ShutdownManager) AddStep(name string, fn func(ctx context.Context) error) {
	sm.Steps = append(sm.Steps, ShutdownStep{
		Name: name,
		Func: fn,
	})
}

func (sm *ShutdownManager) WaitForSignal(ctx context.Context) <-chan struct{} {
	done := make(chan struct{})

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		defer signal.Stop(signals)

		<-signals
		sm.Logger.Info().Msg("shutdown signal received. Starting graceful shutdown...")

		timeoutCtx, cancel := context.WithTimeout(ctx, sm.Timeout)
		defer cancel()

		timeoutFunc := time.AfterFunc(sm.Timeout, func() {
			sm.Logger.Error().Msg("graceful shutdown timed out. Forcing exit.")
			os.Exit(1)
		})
		defer timeoutFunc.Stop()

		var wg sync.WaitGroup
		for _, step := range sm.Steps {
			wg.Add(1)

			go func(s ShutdownStep) {
				defer wg.Done()
				sm.Logger.Info().
					Str("step", s.Name).
					Msg("executing shutdown step")

				if err := s.Func(timeoutCtx); err != nil {
					sm.Logger.Error().
						Err(err).
						Str("step", s.Name).
						Msg("shutdown step failed")
					return
				}

				sm.Logger.Info().
					Str("step", s.Name).
					Msg("shutdown step completed")
			}(step)
		}

		wg.Wait()
		close(done)
	}()

	return done
}
