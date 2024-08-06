package shutdown

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// shutdow.Handler sets up signal handling for clean shutdown and executes the provided shutdown function
func Handler(ctx context.Context, shutdownFunc func(context.Context) error, timeout time.Duration, cancelFunc context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		shutdownCtx, shutdownCancel := context.WithTimeout(ctx, timeout)
		defer shutdownCancel()
		if err := shutdownFunc(shutdownCtx); err != nil {
			fmt.Println("Failed to shutdown:", err)
		}
		cancelFunc()
	}()
}
