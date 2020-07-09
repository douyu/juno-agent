module github.com/douyu/juno-agent

go 1.14

require (
	github.com/apache/rocketmq-client-go/v2 v2.0.0-rc2
	github.com/armon/circbuf v0.0.0-20150827004946-bbbad097214e
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/douyu/jupiter v0.0.0-00010101000000-000000000000
	github.com/fsnotify/fsnotify v1.4.9
	github.com/garyburd/redigo v1.6.0
	github.com/go-resty/resty/v2 v2.2.0
	github.com/jinzhu/gorm v1.9.12
	github.com/json-iterator/go v1.1.9
	github.com/labstack/echo/v4 v4.1.16
	github.com/labstack/gommon v0.3.0
	github.com/nats-io/nats-server/v2 v2.1.6 // indirect
	github.com/nats-io/nats.go v1.9.2
	github.com/stretchr/testify v1.6.1
	github.com/uber-go/atomic v1.4.0
	github.com/yangchenxing/go-nginx-conf-parser v0.0.0-20190110023421-0d59f1b7a3f6
	go.etcd.io/etcd v0.5.0-alpha.5.0.20200425165423-262c93980547
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e
	google.golang.org/grpc v1.29.0
	gopkg.in/ini.v1 v1.56.0
)

replace github.com/douyu/jupiter => ../jupiter
