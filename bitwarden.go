package bitwarden

import (
	"context"
	"fmt"
	"os/exec"

	"golang.org/x/exp/slog"
)

// StartBWServe starts up the local API endpoint in the background and
// supervises it. It also maintains the lifecycle.
func StartBWServe(ctx context.Context) {
	cmd := exec.Command("bw", "serve")
	if err := cmd.Start(); err != nil {
		slog.Error("failed to run 'bw serve'", err, slog.String("component", "bitwarden"))
	}
	<-ctx.Done()
	cmd.Process.Kill()
}

// DataFromBWItem has the goal to convert data from a Bitwarden item to be ready
// to be inserted into Vault
func DataFromBWItem(resp map[string]interface{}) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	if resp["login"] != nil {
		login := resp["login"].(map[string]interface{})
		if login["username"] != nil {
			data["username"] = login["username"]
		} else {
			return nil, fmt.Errorf("bw item has no username")
		}
		if login["password"] != nil {
			data["password"] = login["password"]
		} else {
			return nil, fmt.Errorf("bw item has no password")
		}
	}
	return data, nil
}
