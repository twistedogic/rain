package archive

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

const (
	symbolsURL  = "https://api.binance.com/api/v3/exchangeInfo"
	dataBaseURL = "https://data.binance.vision/"
	candleURL   = dataBaseURL + "data/spot/klines/"
	tradeURL    = dataBaseURL + "data/spot/trades/"

	timeBucket        = time.Minute
	requestsPerBucket = 1200
)

type Client struct {
	*http.Client
	*rate.Limiter
}

func New() Client {
	return Client{
		Client:  http.DefaultClient,
		Limiter: rate.NewLimiter(rate.Every(timeBucket), requestsPerBucket),
	}
}

func (c Client) Do(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	if err := c.Wait(ctx); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c Client) getResponse(req *http.Request, out io.Writer) error {
	res, err := c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	_, err = io.Copy(out, res.Body)
	return err
}

func (c Client) getJSON(req *http.Request, i interface{}) error {
	buf := new(bytes.Buffer)
	if err := c.getResponse(req, buf); err != nil {
		return err
	}
	return json.Unmarshal(buf.Bytes(), i)
}

func (c Client) getCSV(req *http.Request) ([][]string, error) {
	records := make([][]string, 0)
	buf := new(bytes.Buffer)
	if err := c.getResponse(req, buf); err != nil {
		return nil, err
	}
	br := bytes.NewReader(buf.Bytes())
	zr, err := zip.NewReader(br, int64(buf.Len()))
	if err != nil {
		return nil, err
	}
	for _, f := range zr.File {
		r, err := f.Open()
		if err != nil {
			return nil, err
		}
		cr := csv.NewReader(r)
		rows, err := cr.ReadAll()
		if err != nil {
			return nil, err
		}
		records = append(records, rows...)
	}
	return records, nil
}
