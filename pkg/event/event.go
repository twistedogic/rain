package event

import (
	"time"
)

type Symbol struct {
	Name  string `json:"symbol"`
	Base  string `json:"baseAsset"`
	Quote string `json:"quoteAsset"`
}

type SymbolTicker struct {
	Symbol                  string
	Open, High, Low, Close  float64
	BaseVolume, QuoteVolume float64
}

type Candle struct {
	SymbolTicker
	OpenTime, CloseTime               time.Time
	NumberOfTrades                    int
	TakerBaseVolume, TakerQuoteVolume float64
}

type Trade struct {
	Symbol                   string
	Id                       int
	Price, BaseQty, QuoteQty float64
	Time                     time.Time
}

type Ask struct {
	Price, Quantity float64
}

type Bid struct {
	Price, Quantity float64
}

type OrderBook struct {
	Symbol                  string
	Time                    time.Time
	UpdateID, FirstUpdateID int64
	Bids                    []Bid
	Asks                    []Ask
}

type MiniTicker struct {
	SymbolTicker
	Event string
	Time  time.Time
}

type Ticker struct {
	MiniTicker
	OpenTime, CloseTime              time.Time
	PriceChange, PriceChangePercent  float64
	WeightedAvgPrice, PrevClosePrice float64
	CloseQty                         float64
	BidPrice, BidQty                 float64
	AskPrice, AskQty                 float64
	NumberOfTrades                   int
}
