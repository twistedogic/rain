package subscribe

import (
	"context"
	"time"

	"github.com/adshao/go-binance/v2"

	"github.com/twistedogic/rain/pkg/event"
)

func handleSubscriptionWithContext(
	ctx context.Context, errCh chan error,
	doneCh, stopCh chan struct{}, err error) error {
	if err != nil {
		errCh <- err
	}
	defer func() {
		stopCh <- struct{}{}
		<-doneCh
	}()
	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}

func ParseOrderBook(e *binance.WsDepthEvent, entry *event.OrderBook) error {
	var err error
	entry.Symbol = e.Symbol
	entry.Time = time.Unix(0, e.Time*int64(time.Millisecond))
	entry.UpdateID, entry.FirstUpdateID = e.UpdateID, e.FirstUpdateID
	asks := make([]event.Ask, len(e.Asks))
	for i := range e.Asks {
		asks[i].Price, err = event.ParseFloatField(err, e.Asks[i].Price)
		asks[i].Quantity, err = event.ParseFloatField(err, e.Asks[i].Quantity)
	}
	bids := make([]event.Bid, len(e.Bids))
	for i := range e.Bids {
		bids[i].Price, err = event.ParseFloatField(err, e.Bids[i].Price)
		bids[i].Quantity, err = event.ParseFloatField(err, e.Bids[i].Quantity)
	}
	entry.Asks, entry.Bids = asks, bids
	return err
}

func SubscribeSymbolOrderBook(ctx context.Context, symbol string, eventCh chan event.OrderBook) error {
	errCh := make(chan error, 1)
	eventHandler := func(e *binance.WsDepthEvent) {
		var entry event.OrderBook
		if err := ParseOrderBook(e, &entry); err != nil {
			errCh <- err
			return
		}
		eventCh <- entry
	}
	errHandler := func(err error) {
		errCh <- err
	}
	doneCh, stopCh, err := binance.WsDepthServe(symbol, eventHandler, errHandler)
	return handleSubscriptionWithContext(ctx, errCh, doneCh, stopCh, err)
}

func ParseTrade(e *binance.WsTradeEvent, entry *event.Trade) error {
	var err error
	entry.Symbol = e.Symbol
	entry.Time = time.Unix(0, e.TradeTime*int64(time.Millisecond))
	entry.Id = int(e.TradeID)
	entry.Price, err = event.ParseFloatField(err, e.Price)
	entry.BaseQty, err = event.ParseFloatField(err, e.Quantity)
	return err
}

func SubscribeSymbolTrades(ctx context.Context, symbol string, eventCh chan event.Trade) error {
	errCh := make(chan error, 1)
	eventHandler := func(e *binance.WsTradeEvent) {
		var entry event.Trade
		if err := ParseTrade(e, &entry); err != nil {
			errCh <- err
			return
		}
		eventCh <- entry
	}
	errHandler := func(err error) {
		errCh <- err
	}
	doneCh, stopCh, err := binance.WsTradeServe(symbol, eventHandler, errHandler)
	return handleSubscriptionWithContext(ctx, errCh, doneCh, stopCh, err)
}

/*
func ParseCandle(e *binance.WsTradeEvent, entry *event.Candle) error {
	var err error
	return err
}

func SubscribeSymbolCandles(ctx context.Context, symbol, interval string, eventCh chan event.Candle) error {
	errCh := make(chan error, 1)
	eventHandler := func(event *binance.WsKlineEvent) {
		var entry event.Candle
		if err := ParseCandle(e, &entry); err != nil {
			errCh <- err
			return
		}
		eventCh <- entry
	}
	errHandler := func(err error) {
		errCh <- err
	}
	doneCh, stopCh, err := binance.WsKlineServe(symbol, interval, eventHandler, errHandler)
	return handleSubscriptionWithContext(ctx, errCh, doneCh, stopCh, err)
}

func ParseMiniTicker() error {}
*/
