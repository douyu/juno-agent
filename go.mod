module github.com/douyu/juno-agent

go 1.14

require (
	github.com/HdrHistogram/hdrhistogram-go v1.1.0 // indirect
	github.com/apache/rocketmq-client-go/v2 v2.0.0
	github.com/armon/circbuf v0.0.0-20150827004946-bbbad097214e
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/coreos/etcd v3.3.22+incompatible
	github.com/douyu/jupiter v0.2.7
	github.com/fsnotify/fsnotify v1.4.9
	github.com/garyburd/redigo v1.6.0
	github.com/go-resty/resty/v2 v2.3.0
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/google/btree v1.0.1-0.20191016161528-479b5e81b0a9 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/jinzhu/gorm v1.9.12
	github.com/json-iterator/go v1.1.10
	github.com/labstack/echo/v4 v4.1.16
	github.com/labstack/gommon v0.3.0
	github.com/mitchellh/mapstructure v1.4.0 // indirect
	github.com/nats-io/nats.go v1.9.2
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/robfig/cron/v3 v3.0.1
	github.com/sony/sonyflake v1.0.0
	github.com/stretchr/testify v1.7.0
	github.com/tidwall/gjson v1.6.0 // indirect
	github.com/uber-go/atomic v1.4.0
	github.com/uber/jaeger-client-go v2.25.0+incompatible // indirect
	github.com/uber/jaeger-lib v2.4.0+incompatible // indirect
	github.com/yangchenxing/go-nginx-conf-parser v0.0.0-20190110023421-0d59f1b7a3f6
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0
	golang.org/x/net v0.0.0-20201216054612-986b41b23924 // indirect
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a
	golang.org/x/sys v0.0.0-20201218084310-7d0127a74742 // indirect
	google.golang.org/genproto v0.0.0-20201214200347-8c77b98c765d // indirect
	google.golang.org/grpc v1.34.0
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/ini.v1 v1.56.0
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
