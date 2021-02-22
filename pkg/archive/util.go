package archive

import (
	"strconv"
	"time"
)

const (
	dateFormat = "2006-01"

	day  = 24 * time.Hour
	week = 7 * day
)

func parseIntField(err error, s string) (int, error) {
	if err != nil {
		return 0, err
	}
	i, err := strconv.ParseInt(s, 10, 64)
	return int(i), err
}

func parseFloatField(err error, s string) (float64, error) {
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(s, 64)
}

func parseTimeField(err error, s string) (time.Time, error) {
	var ms int
	if err != nil {
		return time.Time{}, err
	}
	ms, err = parseIntField(err, s)
	return time.Unix(0, int64(ms)*int64(time.Millisecond)), err
}

func getMonthlyDateRange(start, end time.Time) []string {
	n := end.Month() - start.Month()
	dateRange := make([]string, n+1)
	for i := 0; i <= int(n); i++ {
		dateRange[i] = start.AddDate(0, i, 0).UTC().Format(dateFormat)
	}
	return dateRange
}

func getWindow(w time.Duration) string {
	window := "1m"
	switch {
	case w > week:
		window = "1mo"
	case w/week > 0:
		window = "1w"
	case w/(3*day) > 0:
		window = "3d"
	case w/day > 0:
		window = "1d"
	case w/(12*time.Hour) > 0:
		window = "12h"
	case w/(8*time.Hour) > 0:
		window = "8h"
	case w/(6*time.Hour) > 0:
		window = "6h"
	case w/(4*time.Hour) > 0:
		window = "4h"
	case w/(2*time.Hour) > 0:
		window = "2h"
	case w/time.Hour > 0:
		window = "1h"
	case w/(30*time.Minute) > 0:
		window = "4h"
	case w/(15*time.Minute) > 0:
		window = "15m"
	case w/(5*time.Minute) > 0:
		window = "5m"
	case w/(3*time.Minute) > 0:
		window = "3m"
	}
	return window
}
