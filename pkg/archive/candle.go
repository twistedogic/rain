package archive

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/twistedogic/rain/pkg/event"
)

func parseCandle(symbol string, row []string, c *event.Candle) error {
	var err error
	s := event.SymbolTicker{}
	s.Symbol = symbol
	c.OpenTime, err = event.ParseTimeField(err, row[0])
	s.Open, err = event.ParseFloatField(err, row[1])
	s.High, err = event.ParseFloatField(err, row[2])
	s.Low, err = event.ParseFloatField(err, row[3])
	s.Close, err = event.ParseFloatField(err, row[4])
	s.BaseVolume, err = event.ParseFloatField(err, row[5])
	c.CloseTime, err = event.ParseTimeField(err, row[6])
	s.QuoteVolume, err = event.ParseFloatField(err, row[7])
	c.NumberOfTrades, err = event.ParseIntField(err, row[8])
	c.TakerBaseVolume, err = event.ParseFloatField(err, row[9])
	c.TakerQuoteVolume, err = event.ParseFloatField(err, row[10])
	c.SymbolTicker = s
	return err
}

func getSymbolCandlesURLs(symbol string, start, end time.Time, window time.Duration) []string {
	s := strings.ToUpper(symbol)
	ranges := getMonthlyDateRange(start, end)
	interval := getWindow(window)
	urls := make([]string, len(ranges))
	for i, r := range ranges {
		filename := fmt.Sprintf("%s-%s-%s.zip", s, interval, r)
		urls[i] = fmt.Sprintf("%s%s/%s/%s", candleURL, s, interval, filename)
	}
	return urls
}

func (c Client) GetSymbolCandles(ctx context.Context, s event.Symbol, start, end time.Time, window time.Duration) ([]event.Candle, error) {
	candles := make([]event.Candle, 0)
	urls := getSymbolCandlesURLs(s.Name, start, end, window)
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
			var c event.Candle
			if err := parseCandle(s.Name, row, &c); err != nil {
				return nil, err
			}
			candles = append(candles, c)
		}
	}
	return candles, nil
}
