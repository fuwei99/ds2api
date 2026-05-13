package account

import (
	"testing"
	"time"

	"ds2api/internal/config"
)

func TestPoolSkipsMutedAndInactiveAccounts(t *testing.T) {
	t.Setenv("DS2API_CONFIG_JSON", `{
		"accounts":[
			{"email":"muted@example.com","token":"t1","active":true,"muted":true,"mute_until":4102444800},
			{"email":"inactive@example.com","token":"t2","active":false},
			{"email":"ready@example.com","token":"t3","active":true}
		]
	}`)
	pool := NewPool(config.LoadStore())

	acc, ok := pool.Acquire("", nil)
	if !ok {
		t.Fatal("expected available account")
	}
	if got := acc.Identifier(); got != "ready@example.com" {
		t.Fatalf("expected ready account, got %q", got)
	}
}

func TestPoolClearsExpiredMute(t *testing.T) {
	t.Setenv("DS2API_CONFIG_JSON", `{
		"accounts":[
			{"email":"expired@example.com","token":"t1","active":true,"muted":true,"mute_until":1}
		]
	}`)
	store := config.LoadStore()
	pool := NewPool(store)

	acc, ok := pool.Acquire("", nil)
	if !ok {
		t.Fatal("expected expired mute account to be available")
	}
	if got := acc.Identifier(); got != "expired@example.com" {
		t.Fatalf("unexpected account %q", got)
	}
	refreshed, ok := store.FindAccount("expired@example.com")
	if !ok {
		t.Fatal("expected account in store")
	}
	if refreshed.IsMuted(time.Now()) {
		t.Fatalf("expected expired mute to be cleared, got %#v", refreshed)
	}
}
