package etcd

// PluginRegProxyPrometheus ETCD dataSource config
type PluginRegProxyPrometheus struct {
	Enable        bool // 是否开启用该数据源
	Path          string
	EnableZone    bool     // 是否开启zone过滤
	Zones         []string // 指定zone
	EnableCleanup bool     // 是否清理无效的yml文件
	TimeInterval  uint32   // 清理间隔多久的历史yml文件，单位s
	Prefixs       []string // 添加多个前缀
}
