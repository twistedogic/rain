package tsdb

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/thanos-io/thanos/pkg/block"
	"github.com/thanos-io/thanos/pkg/gate"
	"github.com/thanos-io/thanos/pkg/model"
	"github.com/thanos-io/thanos/pkg/objstore"
	"github.com/thanos-io/thanos/pkg/store"
	storecache "github.com/thanos-io/thanos/pkg/store/cache"
)

const (
	metaFetchConcurrency       = 20
	maxConcurrency             = 20
	blockSyncConcurrency       = 20
	consistencyDelay           = 0
	ignoreDeletionMarksDelay   = 24 * time.Hour
	IndexCacheSize             = model.Bytes(250 * 1000 * 1000)
	chunkPoolSizeBytes         = model.Bytes(2048 * 1000 * 1000)
	lazyIndexReaderIdleTimeout = 5 * time.Minute
)

func NewBucketStore(dataDir string, bkt objstore.Bucket) (*store.BucketStore, error) {
	reg := prometheus.DefaultRegisterer
	indexCacheConfig := storecache.InMemoryIndexCacheConfig{
		MaxSize:     IndexCacheSize,
		MaxItemSize: storecache.DefaultInMemoryIndexCacheConfig.MaxItemSize,
	}
	indexCache, err := storecache.NewInMemoryIndexCacheWithConfig(
		logger, reg, indexCacheConfig,
	)
	if err != nil {
		return nil, errors.Wrap(err, "create index cache")
	}
	instrumentedBkt := objstore.NewTracingBucket(bkt)

	ignoreDeletionMarkFilter := block.NewIgnoreDeletionMarkFilter(
		logger, instrumentedBkt,
		ignoreDeletionMarksDelay, metaFetchConcurrency,
	)
	metaFetcher, err := block.NewMetaFetcher(
		logger, metaFetchConcurrency, instrumentedBkt, dataDir, reg,
		[]block.MetadataFilter{
			block.NewConsistencyDelayMetaFilter(logger, consistencyDelay, reg),
			ignoreDeletionMarkFilter,
			block.NewDeduplicateFilter(),
		}, nil)
	if err != nil {
		return nil, errors.Wrap(err, "meta fetcher")
	}
	queriesGate := gate.New(reg, maxConcurrency)
	if err != nil {
		return nil, errors.Wrap(err, "create chunk pool")
	}
	bs, err := store.NewBucketStore(
		logger, reg, instrumentedBkt, metaFetcher, dataDir,
		indexCache, queriesGate,
		uint64(chunkPoolSizeBytes),
		store.NewChunksLimiterFactory(0),
		store.NewSeriesLimiterFactory(0),
		true,
		blockSyncConcurrency,
		new(store.FilterConfig), // filter config for minTime and maxTime
		false,
		store.DefaultPostingOffsetInMemorySampling,
		false,
		false, // Enable lazy index reader
		lazyIndexReaderIdleTimeout,
	)
	if err != nil {
		return nil, errors.Wrap(err, "create object storage store")
	}
	if err := bs.InitialSync(context.Background()); err != nil {
		return nil, errors.Wrap(err, "bucket store initial sync")
	}
	return bs, nil
}
