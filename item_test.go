package bitwarden_test

import (
	"context"
	"os"
	"os/signal"
	"testing"

	"github.com/fbegyn/bitwarden-api"
)

func TestItemCRUD(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	_, err := bitwarden.StartBWServe(ctx)
	if err != nil {
		t.Fatal(err)
	}

	client, err := bitwarden.NewBitwardenClient(bitwarden.Config{
		Addr: "http://localhost:8087",
	})
	if err != nil {
		t.Fatal(err)
	}

	test := bitwarden.Item{
		Type: 1,
		Name: "bitwarden-api/test/path",
		Login: bitwarden.Login{
			Username: "hello",
			Password: "world",
		},
	}

	err = client.CreateItem(test)
	if err != nil {
		t.Fatal(err)
	}

	items, err := client.ListItems(test.Name)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) == 0 {
		t.Fatalf("no items found, at least 1 should be found")
	}

	get, err := client.GetItem(items[0])
	if err != nil {
		t.Fatal(err)
	}
	if get.Name != items[0].Name || get.ID != items[0].ID {
		t.Fatalf("the fetched item does not match the created one")
	}

	items[0].Login.Username = "foo"
	items[0].Login.Password = "bar"
	err = client.UpdateItem(items[0])
	if err != nil {
		t.Fatal(err)
	}

	for _, it := range items {
		err := client.DeleteItem(it)
		if err != nil {
			t.Fatal(err)
		}
	}
	items, err = client.ListItems(test.Name)
	if err != nil {
		t.Fatal(err)
	}

	if len(items) != 0 {
		t.Log(items)
		t.Fatalf("expected 0 items, got %d", len(items))
	}
}
