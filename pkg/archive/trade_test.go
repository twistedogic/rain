package archive

import (
	"context"
	"testing"
	"time"

	"github.com/twistedogic/rain/pkg/event"
)

func Test_Client_GetSymbolTrades(t *testing.T) {
	t.Skip()
	symbol := event.Symbol{Name: "ETHBTC"}
	start, end := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC)
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()
	got, err := New().GetSymbolTrades(ctx, symbol, start, end)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) == 0 {
		t.Fatal(got)
	}
}
