package etcd

// PluginRegProxyPrometheus ETCD dataSource config
type PluginRegProxyPrometheus struct {
	Enable       bool // 是否开启用该数据源
	Path         string
	EnableRegion bool     // 是否开启region过滤
	Region       []string // 指定region
}
