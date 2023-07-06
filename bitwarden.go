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

// DataFromBWItem has the goal to convert data from a Bitwarden item to be ready
// to be inserted into Vault
func DataFromBWItem(item Item) (map[string]interface{}, error) {
	data := make(map[string]interface{})

	if item.Login.Password == "" {
		return nil, fmt.Errorf("bw item has no password")
	}
	data["password"] = item.Login.Password
	if item.Login.Username == "" {
		return nil, fmt.Errorf("bw item has no username")
	}
	data["username"] = item.Login.Username

	for _, field := range item.Fields {
		data[field.Name] = field.Value
	}
	return data, nil
}
