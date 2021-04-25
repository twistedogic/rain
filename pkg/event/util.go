package event

import (
	"strconv"
	"time"
)

func ParseIntField(err error, s string) (int, error) {
	if err != nil {
		return 0, err
	}
	i, err := strconv.ParseInt(s, 10, 64)
	return int(i), err
}

func ParseFloatField(err error, s string) (float64, error) {
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(s, 64)
}

func ParseTimeField(err error, s string) (time.Time, error) {
	var ms int
	if err != nil {
		return time.Time{}, err
	}
	ms, err = ParseIntField(err, s)
	return time.Unix(0, int64(ms)*int64(time.Millisecond)), err
}
