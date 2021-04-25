package tsdb

import (
	"context"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/tsdb"
	"github.com/thanos-io/thanos/pkg/errutil"
	"github.com/thanos-io/thanos/pkg/objstore"
	"github.com/thanos-io/thanos/pkg/receive"
	"github.com/thanos-io/thanos/pkg/runutil"
	"github.com/thanos-io/thanos/pkg/store/storepb"
)

const (
	DefaultTenantLabel = "tenant"
	DefaultTenantID    = "tsdb"
)

var DefaultTSDBOptions = &tsdb.Options{
	MinBlockDuration:  int64(2 * time.Hour / time.Millisecond),
	MaxBlockDuration:  int64(2 * time.Hour / time.Millisecond),
	RetentionDuration: int64(6 * time.Hour / time.Millisecond),
	NoLockfile:        true,
}

type Sample struct {
	Timestamp int64
	Value     float64
}

type TimeSeries struct {
	Labels  labels.Labels
	Samples []Sample
}

type Receiver struct {
	*receive.MultiTSDB
	dataDir string
}

func New(dataDir string, bucket objstore.Bucket, l labels.Labels) Receiver {
	logger := log.NewLogfmtLogger(os.Stderr)
	m := receive.NewMultiTSDB(
		dataDir, logger, prometheus.DefaultRegisterer, DefaultTSDBOptions,
		l, DefaultTenantLabel,
		bucket, false,
	)
	return Receiver{MultiTSDB: m, dataDir: dataDir}
}

func (r Receiver) Write(ctx context.Context, timeseries []TimeSeries) error {
	var (
		numOutOfOrder  = 0
		numDuplicates  = 0
		numOutOfBounds = 0
	)
	s, err := r.TenantAppendable(DefaultTenantID)
	if err != nil {
		return errors.Wrap(err, "get tenant appendable")
	}

	var app storage.Appender
	if err := runutil.Retry(1*time.Second, ctx.Done(), func() error {
		var err error
		app, err = s.Appender(ctx)
		return err
	}); err != nil {
		return errors.Wrap(err, "get appender")
	}
	var errs errutil.MultiError
	for _, t := range timeseries {
		lset := t.Labels
		for _, s := range t.Samples {
			_, err = app.Add(lset, s.Timestamp, s.Value)
			switch err {
			case nil:
				continue
			case storage.ErrOutOfOrderSample:
				numOutOfOrder++
			case storage.ErrDuplicateSampleForTimestamp:
				numDuplicates++
			case storage.ErrOutOfBounds:
				numOutOfBounds++
			}
		}
	}
	if numOutOfOrder > 0 {
		errs.Add(errors.Wrapf(storage.ErrOutOfOrderSample, "failed to non-fast add %d samples", numOutOfOrder))
	}
	if numDuplicates > 0 {
		errs.Add(errors.Wrapf(storage.ErrDuplicateSampleForTimestamp, "failed to non-fast add %d samples", numDuplicates))
	}
	if numOutOfBounds > 0 {
		errs.Add(errors.Wrapf(storage.ErrOutOfBounds, "failed to non-fast add %d samples", numOutOfBounds))
	}
	if err := app.Commit(); err != nil {
		errs.Add(errors.Wrap(err, "commit samples"))
	}
	return errs.Err()
}

func (r Receiver) Stores() []storepb.StoreServer {
	stores := r.TSDBStores()
	servers := make([]storepb.StoreServer, 0, len(stores))
	for _, s := range stores {
		servers = append(servers, s)
	}
	return servers
}
