package archive

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Candle struct {
	Symbol                                           string
	OpenTime, CloseTime                              time.Time
	Open, High, Low, Close, Volume, QuoteAssetVolume float64
	NumberOfTrades                                   int
	TakerBaseAssetVolume, TakerQuoteAssetVolume      float64
}

func parseCandle(symbol string, row []string, c *Candle) error {
	var err error
	c.Symbol = symbol
	c.OpenTime, err = parseTimeField(err, row[0])
	c.Open, err = parseFloatField(err, row[1])
	c.High, err = parseFloatField(err, row[2])
	c.Low, err = parseFloatField(err, row[3])
	c.Close, err = parseFloatField(err, row[4])
	c.Volume, err = parseFloatField(err, row[5])
	c.CloseTime, err = parseTimeField(err, row[6])
	c.QuoteAssetVolume, err = parseFloatField(err, row[7])
	c.NumberOfTrades, err = parseIntField(err, row[8])
	c.TakerBaseAssetVolume, err = parseFloatField(err, row[9])
	c.TakerQuoteAssetVolume, err = parseFloatField(err, row[10])
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

func (c Client) GetSymbolCandles(ctx context.Context, s Symbol, start, end time.Time, window time.Duration) ([]Candle, error) {
	candles := make([]Candle, 0)
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
			var c Candle
			if err := parseCandle(s.Name, row, &c); err != nil {
				return nil, err
			}
			candles = append(candles, c)
		}
	}
	return candles, nil
}
