package bitwarden_test

import (
	"context"
	"net"
	"os"
	"os/signal"
	"testing"

	"github.com/fbegyn/bitwarden-api"
)

func TestStartBWServe(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)

	cmd, err := bitwarden.StartBWServe(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("process for bw serve: %v", cmd.Process)

	conn, err := net.Dial("tcp", "localhost:8087")
	if err != nil {
		cancel()
		t.Fatal(err)
	}
	_ = conn

	cancel()
	cmd.Wait()
	if cmd.ProcessState.String() != "signal: killed" {
		t.Fatalf("The process was not killed: %s", cmd.ProcessState.String())
	}
}
