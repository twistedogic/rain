package archive

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/twistedogic/rain/pkg/event"
)

func parseTrade(symbol string, row []string, t *event.Trade) error {
	var err error
	t.Symbol = symbol
	t.Id, err = event.ParseIntField(err, row[0])
	t.Price, err = event.ParseFloatField(err, row[1])
	t.BaseQty, err = event.ParseFloatField(err, row[2])
	t.QuoteQty, err = event.ParseFloatField(err, row[3])
	t.Time, err = event.ParseTimeField(err, row[4])
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

func (c Client) GetSymbolTrades(ctx context.Context, s event.Symbol, start, end time.Time) ([]event.Trade, error) {
	trades := make([]event.Trade, 0)
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
			var t event.Trade
			if err := parseTrade(s.Name, row, &t); err != nil {
				return nil, err
			}
			trades = append(trades, t)
		}
	}
	return trades, nil
}
