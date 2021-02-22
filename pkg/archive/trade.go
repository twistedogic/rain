package archive

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Trade struct {
	Symbol                   string
	Id                       int
	Price, BaseQty, QuoteQty float64
	Time                     time.Time
}

func parseTrade(symbol string, row []string, t *Trade) error {
	var err error
	t.Symbol = symbol
	t.Id, err = parseIntField(err, row[0])
	t.Price, err = parseFloatField(err, row[1])
	t.BaseQty, err = parseFloatField(err, row[2])
	t.QuoteQty, err = parseFloatField(err, row[3])
	t.Time, err = parseTimeField(err, row[4])
	return err
}

func getSymbolTradesURLs(symbol string, start, end time.Time) []string {
	s := strings.ToUpper(symbol)
	ranges := getMonthlyDateRange(start, end)
	urls := make([]string, len(ranges))
	for i, r := range ranges {
		filename := fmt.Sprintf("%s-trades-%s.zip", s, r)
		urls[i] = fmt.Sprintf("%s%s/%s", tradeURL, s, filename)
	}
	return urls
}

func (c Client) GetSymbolTrades(ctx context.Context, s Symbol, start, end time.Time) ([]Trade, error) {
	trades := make([]Trade, 0)
	urls := getSymbolTradesURLs(s.Name, start, end)
	for _, u := range urls {
		req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
		if err != nil {
			return nil, err
		}
		rows, err := c.getCSV(req)
		if err != nil {
			return nil, err
		}
		for _, row := range rows {
			var t Trade
			if err := parseTrade(s.Name, row, &t); err != nil {
				return nil, err
			}
			trades = append(trades, t)
		}
	}
	return trades, nil
}
