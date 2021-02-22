package archive

import (
	"context"
	"net/http"
)

type Symbol struct {
	Name  string `json:"symbol"`
	Base  string `json:"baseAsset"`
	Quote string `json:"quoteAsset"`
}

type exchangeInfo struct {
	Symbols []Symbol `json:"symbols"`
}

func (c Client) GetAllSymbols(ctx context.Context) ([]Symbol, error) {
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
