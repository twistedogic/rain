package tsdb

import (
	"context"
	"math"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/storage"
	"github.com/thanos-io/thanos/pkg/component"
	"github.com/thanos-io/thanos/pkg/query"
	"github.com/thanos-io/thanos/pkg/store"
	"github.com/thanos-io/thanos/pkg/store/storepb"
)

const (
	maxSourceResolution = 0
	enableDedup         = true
)

type inProcessClient struct{ storepb.StoreClient }

func (i inProcessClient) TimeRange() (int64, int64) {
	return 0, time.Now().UnixNano()
}
func (i inProcessClient) LabelSets() []labels.Labels {
	return []labels.Labels{labels.FromStrings(DefaultTenantLabel, DefaultTenantID)}
}
func (i inProcessClient) Addr() string   { return "localhost" }
func (i inProcessClient) String() string { return "localhost" }

func GetStoreClients(servers ...storepb.StoreServer) func() []store.Client {
	clients := make([]store.Client, len(servers))
	for i, s := range servers {
		clients[i] = inProcessClient{storepb.ServerAsClient(s, 0)}
	}
	return func() []store.Client {
		return clients
	}
}

func NewQuerier(timeout time.Duration, storeServers ...storepb.StoreServer) storage.Queryable {
	reg := prometheus.DefaultRegisterer
	selectorLset := labels.FromStrings(DefaultTenantLabel, DefaultTenantID)
	proxy := store.NewProxyStore(
		logger, reg,
		GetStoreClients(storeServers...),
		component.Query,
		selectorLset, timeout,
	)
	replicaLabels := []string{DefaultTenantLabel}
	return query.NewQueryableCreator(
		logger,
		reg,
		proxy,
		len(storeServers),
		time.Duration(len(storeServers))*timeout,
	)(
		enableDedup, replicaLabels, nil, maxSourceResolution, true, false,
	)
}

func NewEngine(timeout time.Duration) *promql.Engine {
	reg := prometheus.DefaultRegisterer
	engineOpts := promql.EngineOpts{
		Logger:     logger,
		Reg:        reg,
		MaxSamples: math.MaxInt32,
		Timeout:    timeout,
	}
	return promql.NewEngine(engineOpts)
}

type QueryEngine struct {
	*promql.Engine
	query storage.Queryable
}

func NewQueryEngine(servers ...storepb.StoreServer) QueryEngine {
	timeout := 5 * time.Minute
	return QueryEngine{
		Engine: NewEngine(time.Duration(len(servers)) * timeout),
		query:  NewQuerier(timeout, servers...),
	}
}

func (q QueryEngine) GetInstant(ctx context.Context, qs string, ts time.Time) (*promql.Result, error) {
	query, err := q.NewInstantQuery(q.query, qs, ts)
	if err != nil {
		return nil, err
	}
	return query.Exec(ctx), nil
}

func (q QueryEngine) GetRange(ctx context.Context, qs string, start, end time.Time, step time.Duration) (*promql.Result, error) {
	query, err := q.NewRangeQuery(q.query, qs, start, end, step)
	if err != nil {
		return nil, err
	}
	return query.Exec(ctx), nil
}
