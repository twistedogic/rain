package archive

import (
	"context"
	"net/http"

	"github.com/twistedogic/rain/pkg/event"
)

type exchangeInfo struct {
	Symbols []event.Symbol `json:"symbols"`
}

func (c Client) GetAllSymbols(ctx context.Context) ([]event.Symbol, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", symbolsURL, nil)
	if err != nil {
		return nil, err
	}
	var info exchangeInfo
	if err := c.getJSON(req, &info); err != nil {
		return nil, err
	}
	return info.Symbols, nil
}
