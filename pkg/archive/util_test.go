package archive

import (
	"sort"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func Test_getMonthlyDateRange(t *testing.T) {
	cases := map[string]struct {
		start, end time.Time
		want       []string
	}{
		"base": {
			start: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			end:   time.Date(2021, 4, 1, 0, 0, 0, 0, time.UTC),
			want:  []string{"2021-01", "2021-02", "2021-03", "2021-04"},
		},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			got := getMonthlyDateRange(tc.start, tc.end)
			sort.Strings(got)
			sort.Strings(tc.want)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func Test_parseTimeField(t *testing.T) {
	cases := []struct {
		input string
		want  time.Time
	}{
		{"1601510340000", time.Date(2020, 9, 30, 23, 59, 0, 0, time.UTC)},
	}
	for i := range cases {
		tc := cases[i]
		t.Run(tc.input, func(t *testing.T) {
			got, err := parseTimeField(nil, tc.input)
			if err != nil {
				t.Fatal(err)
			}
			if !got.Equal(tc.want) {
				t.Fatalf("want: %s, got: %s", tc.want, got)
			}
		})
	}
}
