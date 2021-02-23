package book

import (
	"context"
	"testing"
	"time"

	"github.com/adshao/go-binance/v2"
)

func Test_WatchOrderBook(t *testing.T) {
	t.Skip()
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	symbol := "ETHBTC"
	eventCh := make(chan *binance.WsDepthEvent)
	go func() {
		defer close(eventCh)
		WatchOrderBook(ctx, symbol, eventCh)
	}()
	for event := range eventCh {
		if event.Symbol != symbol {
			t.Fatal(event)
		}
	}
}
