module github.com/douyu/juno-agent

go 1.14

require (
	github.com/apache/rocketmq-client-go/v2 v2.0.0-rc2
	github.com/armon/circbuf v0.0.0-20150827004946-bbbad097214e
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/coreos/etcd v3.3.22+incompatible
	github.com/douyu/jupiter v0.2.4
	github.com/fsnotify/fsnotify v1.4.9
	github.com/garyburd/redigo v1.6.0
	github.com/go-resty/resty/v2 v2.2.0
	github.com/google/btree v1.0.1-0.20191016161528-479b5e81b0a9 // indirect
	github.com/jinzhu/gorm v1.9.12
	github.com/json-iterator/go v1.1.10
	github.com/labstack/echo/v4 v4.1.16
	github.com/labstack/gommon v0.3.0
	github.com/nats-io/nats-server/v2 v2.1.6 // indirect
	github.com/nats-io/nats.go v1.9.2
	github.com/robfig/cron/v3 v3.0.1
	github.com/stretchr/testify v1.6.1
	github.com/uber-go/atomic v1.4.0
	github.com/yangchenxing/go-nginx-conf-parser v0.0.0-20190110023421-0d59f1b7a3f6
	go.etcd.io/bbolt v1.3.4 // indirect
	go.uber.org/zap v1.15.0
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
	google.golang.org/grpc v1.29.0
	gopkg.in/ini.v1 v1.56.0
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
