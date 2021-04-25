package tsdb

import (
	"context"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/thanos-io/thanos/pkg/objstore"
	"github.com/thanos-io/thanos/pkg/objstore/filesystem"
	"github.com/thanos-io/thanos/pkg/store"
)

func setupTempDir(t *testing.T, prefix ...string) (string, func()) {
	dir, err := ioutil.TempDir("", strings.Join(prefix, "_"))
	if err != nil {
		t.Fatal(err)
	}
	return dir, func() {
		os.RemoveAll(dir)
	}
}

func setupBucket(t *testing.T) (objstore.Bucket, func()) {
	dir, cleanup := setupTempDir(t, "bucket")
	bkt, err := filesystem.NewBucket(dir)
	if err != nil {
		t.Fatal(err)
	}
	return bkt, cleanup
}

func setupReceiver(t *testing.T, bkt objstore.Bucket, ls ...string) (Receiver, func()) {
	dir, cleanup := setupTempDir(t, "receiver")
	r := New(dir, bkt, labels.FromStrings(ls...))
	return r, cleanup
}

func setupStore(t *testing.T, bkt objstore.Bucket) (*store.BucketStore, func()) {
	dir, cleanup := setupTempDir(t, "store")
	bs, err := NewBucketStore(dir, bkt)
	if err != nil {
		t.Fatal(err)
	}
	return bs, cleanup
}

func generateTimeSeries(start, end time.Time, step time.Duration, ls ...string) TimeSeries {
	lset := labels.FromStrings(ls...)
	steps := int(end.Sub(start) / step)
	samples := make([]Sample, steps)
	for i := 0; i < steps; i++ {
		samples[i] = Sample{
			Timestamp: start.Add(time.Duration(i) * step).Unix(),
			Value:     rand.Float64(),
		}
	}
	return TimeSeries{
		Labels:  lset,
		Samples: samples,
	}
}

func Test_Receiver(t *testing.T) {
	logger = log.NewNopLogger()
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC)
	input := []TimeSeries{
		generateTimeSeries(start, end, 30*time.Second, "__name__", "open", "symbol", "a"),
		generateTimeSeries(start, end, 15*time.Second, "__name__", "open", "symbol", "b"),
	}
	bkt, cleanupBucket := setupBucket(t)
	r, cleanupReceiver := setupReceiver(t, bkt)
	s, cleanupStore := setupStore(t, bkt)
	defer func() {
		cleanupStore()
		cleanupReceiver()
		cleanupBucket()
	}()
	ctx := context.TODO()
	if err := r.Write(ctx, input); err != nil {
		t.Fatal(err)
	}
	if err := r.Flush(); err != nil {
		t.Fatal(err)
	}
	if _, err := r.Sync(ctx); err != nil {
		t.Fatal(err)
	}
	if err := s.SyncBlocks(ctx); err != nil {
		t.Fatal(err)
	}
	mint, _ := s.TimeRange()
	if mint != start.Unix() {
		t.Fatalf("mint want: %d, got: %d", start.Unix(), mint)
	}
	stores := append(r.Stores(), s)
	engine := NewQueryEngine(stores...)
	res, err := engine.GetRange(ctx, "{symbol=\"a\"}", start, end, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.String()) == 0 {
		t.Fatal(res)
	}
}
