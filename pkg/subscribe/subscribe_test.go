package subscribe

import (
	"context"
	"testing"
	"time"

	"github.com/twistedogic/rain/pkg/event"
)

func Test_SubscribeSymbolOrderBook(t *testing.T) {
	t.Skip()
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	symbol := "ETHBTC"
	eventCh := make(chan event.OrderBook)
	go func() {
		defer close(eventCh)
		SubscribeSymbolOrderBook(ctx, symbol, eventCh)
	}()
	for event := range eventCh {
		if event.Symbol != symbol {
			t.Fatal(event)
		}
	}
}

func Test_SubscribeSymbolTrades(t *testing.T) {
	t.Skip()
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	symbol := "ETHBTC"
	eventCh := make(chan event.Trade)
	go func() {
		defer close(eventCh)
		SubscribeSymbolTrades(ctx, symbol, eventCh)
	}()
	for event := range eventCh {
		if event.Symbol != symbol {
			t.Fatal(event)
		}
	}
}

/*
func Test_SubscribeSymbolCandles(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	symbol := "ETHBTC"
	eventCh := make(chan event.Candle)
	go func() {
		defer close(eventCh)
		SubscribeSymbolCandles(ctx, symbol, eventCh)
	}()
	for event := range eventCh {
		if event.Symbol != symbol {
			t.Fatal(event)
		}
	}
}
*/
