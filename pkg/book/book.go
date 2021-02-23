package book

import (
	"context"

	"github.com/adshao/go-binance/v2"
)

func WatchOrderBook(ctx context.Context, symbol string, eventCh chan *binance.WsDepthEvent) error {
	errCh := make(chan error, 1)
	eventHandler := func(event *binance.WsDepthEvent) {
		eventCh <- event
	}
	errHandler := func(err error) {
		errCh <- err
	}
	doneC, stopC, err := binance.WsDepthServe(symbol, eventHandler, errHandler)
	if err != nil {
		errCh <- err
	}
	defer func() {
		stopC <- struct{}{}
		<-doneC
	}()
	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}
