package etcd

// PluginRegProxyPrometheus ETCD dataSource config
type PluginRegProxyPrometheus struct {
	Enable     bool // 是否开启用该数据源
	Path       string
	EnableZone bool     // 是否开启zone过滤
	Zones      []string // 指定zone
}
