package etcd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/douyu/juno-agent/util"
	"github.com/douyu/jupiter/pkg/xlog"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
)

func (d *DataSource) watchPrometheus(path string) {
	// etcd的key用作配置数据读取
	hostKey := strings.Join([]string{"/prometheus", "job"}, "/")
	// init watch
	watch, err := d.etcdClient.NewWatch(hostKey)

	if err != nil {
		panic("watch err: " + err.Error())
	}
	go func() {
		for {
			select {
			case event := <-watch.C():
				switch event.Type {
				case mvccpb.DELETE:
					key, value := string(event.Kv.Key), string(event.Kv.Value)
					fmt.Println("key", key, "value", value)

					keyArr := strings.Split(key, "/")
					if len(keyArr) != 5 {
						fmt.Println("key", key, "value", value)
						break
					}
					filename := keyArr[3] + "_" + value
					_ = os.Remove(path + "/" + filename + ".yml")
				case mvccpb.PUT:
					key, value := string(event.Kv.Key), string(event.Kv.Value)
					keyArr := strings.Split(key, "/")
					if len(keyArr) != 5 {
						fmt.Println("key", key, "value", value)
						break
					}
					filename := keyArr[3] + "_" + value
					content := `
- targets:

    - "` + value + `"
  labels:
    instance: ` + keyArr[4] + `
    job: ` + filename
					_ = util.WriteFile(path+"/"+filename+".yml", content)
				}
			}
		}
	}()
}

// PrometheusConfigScanner ..
func (d *DataSource) PrometheusConfigScanner(path string) {

	// etcd的key用作配置数据读取
	hostKey := strings.Join([]string{"/prometheus", "job"}, "/")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	resp, err := d.etcdClient.Get(ctx, hostKey, clientv3.WithPrefix())
	if err != nil {
		xlog.Error("init get hostKey error", xlog.String("plugin", "confgo"), xlog.String("msg", err.Error()), xlog.String("hostKey", hostKey))
		return
	}
	for _, kv := range resp.Kvs {
		key, value := string(kv.Key), string(kv.Value)
		keyArr := strings.Split(key, "/")
		if len(keyArr) != 5 {
			fmt.Println("key", key, "value", value)
			break
		}
		content := `
- targets:

    - "` + value + `"
  labels:
    instance: ` + keyArr[4] + `
    job: ` + keyArr[3]
		util.WriteFile(path+"/"+keyArr[3]+".yml", content)
	}
	return
}
