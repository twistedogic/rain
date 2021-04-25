package archive

import (
	"time"
)

const (
	dateFormat = "2006-01"

	day  = 24 * time.Hour
	week = 7 * day
)

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
