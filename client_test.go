package bitwarden_test

import (
	"context"
	"os"
	"os/signal"
	"testing"

	"github.com/fbegyn/bitwarden-api"
)

func TestNewBitwardenClient(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	_, err := bitwarden.StartBWServe(ctx)
	if err != nil {
		t.Fatal(err)
	}

	conf := bitwarden.Config{
		Addr: "http://localhost:8087",
	}

	client, err := bitwarden.NewBitwardenClient(conf)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.ListItems("")
	if err != nil {
		t.Fatal(err)
	}
}
