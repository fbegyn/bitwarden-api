package bitwarden

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

// StartBWServe starts up the local API endpoint in the background and
// supervises it. It also maintains the lifecycle.
func StartBWServe(ctx context.Context) (*exec.Cmd, error) {
	cmd := exec.Command("bw", "serve")
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to run 'bw serve': %w", err)
	}
	go func() {
		<-ctx.Done()
		cmd.Process.Kill()
	}()

	time.Sleep(2 * time.Second)
	return cmd, nil
}
