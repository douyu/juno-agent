package etcd

// ConfDataSourceEtcd ETCD dataSource config
type ConfDataSourceEtcd struct {
	Enable                        bool // 是否开启用该数据源
	Secure                        bool
	IsWatchPrometheusTargetConfig bool
	PrometheusTargetConfigPath    string
	EndPoints                     []string `json:"endpoints"` // 注册中心etcd节点信息
}
