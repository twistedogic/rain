package event

import (
	"testing"
	"time"
)

func Test_ParseTimeField(t *testing.T) {
	cases := []struct {
		input string
		want  time.Time
	}{
		{"1601510340000", time.Date(2020, 9, 30, 23, 59, 0, 0, time.UTC)},
	}
	for i := range cases {
		tc := cases[i]
		t.Run(tc.input, func(t *testing.T) {
			got, err := ParseTimeField(nil, tc.input)
			if err != nil {
				t.Fatal(err)
			}
			if !got.Equal(tc.want) {
				t.Fatalf("want: %s, got: %s", tc.want, got)
			}
		})
	}
}
