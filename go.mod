module github.com/twistedogic/rain

go 1.14

require (
	github.com/Azure/azure-sdk-for-go v48.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest v0.11.11 // indirect
	github.com/HdrHistogram/hdrhistogram-go v0.9.0 // indirect
	github.com/adshao/go-binance/v2 v2.2.1
	github.com/alicebob/gopher-json v0.0.0-20200520072559-a9ecdc9d1d3a // indirect
	github.com/alicebob/miniredis v2.5.0+incompatible // indirect
	github.com/aws/aws-sdk-go v1.35.31 // indirect
	github.com/dgryski/go-sip13 v0.0.0-20200911182023-62edffca9245 // indirect
	github.com/digitalocean/godo v1.52.0 // indirect
	github.com/felixge/fgprof v0.9.1 // indirect
	github.com/go-kit/kit v0.10.0
	github.com/go-openapi/validate v0.19.14 // indirect
	github.com/go-redis/redis/v8 v8.2.3 // indirect
	github.com/golang/snappy v0.0.2 // indirect
	github.com/google/go-cmp v0.5.4
	github.com/google/pprof v0.0.0-20201117184057-ae444373da19 // indirect
	github.com/gophercloud/gophercloud v0.14.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/hashicorp/consul/api v1.7.0 // indirect
	github.com/hetznercloud/hcloud-go v1.23.1 // indirect
	github.com/influxdata/influxdb v1.8.3 // indirect
	github.com/klauspost/cpuid v1.3.1 // indirect
	github.com/miekg/dns v1.1.35 // indirect
	github.com/minio/minio-go/v7 v7.0.2 // indirect
	github.com/moby/term v0.0.0-20200611042045-63b9a826fb74 // indirect
	github.com/ncw/swift v1.0.52 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/alertmanager v0.21.1-0.20200911160112-1fdff6b3f939 // indirect
	github.com/prometheus/client_golang v1.9.0
	github.com/prometheus/client_model v0.2.0
	github.com/prometheus/common v0.15.0
	github.com/prometheus/prometheus v1.8.2-0.20201119142752-3ad25a6dc3d9
	github.com/samuel/go-zookeeper v0.0.0-20200724154423-2164a8ac840e // indirect
	github.com/shurcooL/vfsgen v0.0.0-20200824052919-0d455de96546 // indirect
	github.com/thanos-io/thanos v0.18.0
	github.com/uber/jaeger-client-go v2.25.0+incompatible // indirect
	github.com/uber/jaeger-lib v2.4.0+incompatible // indirect
	github.com/weaveworks/common v0.0.0-20200914083218-61ffdd448099 // indirect
	github.com/yuin/gopher-lua v0.0.0-20200816102855-ee81675732da // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/goleak v1.1.10 // indirect
	golang.org/x/exp v0.0.0-20200821190819-94841d0725da // indirect
	golang.org/x/oauth2 v0.0.0-20201109201403-9fd604954f58 // indirect
	golang.org/x/time v0.0.0-20210220033141-f8bda1e9f3ba
	golang.org/x/tools v0.0.0-20201119054027-25dc3e1ccc3c // indirect
	google.golang.org/api v0.35.0 // indirect
	k8s.io/api v0.19.4 // indirect
	k8s.io/klog/v2 v2.4.0 // indirect
	k8s.io/utils v0.0.0-20200729134348-d5654de09c73 // indirect
)

replace (
	google.golang.org/grpc => google.golang.org/grpc v1.29.1
	// From Prometheus.
	k8s.io/klog => github.com/simonpasquier/klog-gokit v0.3.0
	k8s.io/klog/v2 => github.com/simonpasquier/klog-gokit/v2 v2.0.1
)
