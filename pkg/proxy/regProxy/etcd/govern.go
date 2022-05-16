package etcd

import (
	"context"
	"encoding/json"
	"io/fs"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/douyu/juno-agent/util"
	"github.com/douyu/jupiter/pkg/xlog"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type governValue struct {
	Addr     string `json:"addr"`
	Region   string `json:"region"`
	Zone     string `json:"zone"`
	Env      string `json:"env"`
	AppName  string `json:"-"`
	Hostname string `json:"-"`
}

func (d *DataSource) parseGovernKey(key string) *governValue {
	govern := new(governValue)

	keyArr := strings.Split(key, "/")
	if len(keyArr) != 4 {
		xlog.Error("parseGovernKey invalid key", xlog.String("key", key))
		return nil
	}

	govern.AppName = strings.Split(keyArr[2], ":")[0]
	govern.Hostname = keyArr[3]

	return govern
}

func (d *DataSource) parseGovern(key, value string) *governValue {
	govern := d.parseGovernKey(key)
	if govern == nil {
		xlog.Error("parseGovernKey failed", xlog.String("key", key), xlog.String("value", value))
		return nil
	}

	err := json.Unmarshal([]byte(value), govern)
	if err != nil {
		xlog.Error("json.Unmarshal failed", xlog.String("key", key), xlog.String("value", value))
		return nil
	}

	return govern
}

func (d *DataSource) filter(govern *governValue) *governValue {
	if govern == nil {
		return nil
	}

	// 指处理指定zone的实例
	for _, zone := range d.zones {
		if zone == govern.Zone {
			return govern
		}
	}

	return nil
}

func (d *DataSource) writeFile(path string, govern *governValue) error {
	filename := govern.AppName + "_" + govern.Hostname
	content := `
- targets:
    - "` + govern.Addr + `"
  labels:
    instance: ` + govern.Hostname + `
    job: ` + govern.AppName
	contentForRoscope := `
- application: ` + govern.AppName + `
  targets:
    - "` + govern.Addr + `"
  labels:
    instance: ` + govern.Hostname + `
    job: ` + govern.AppName
	util.WriteFile(path+"/"+"pyroscope/"+filename+".yml", contentForRoscope)
	return util.WriteFile(path+"/"+filename+".yml", content)
}

func (d *DataSource) watchGovern(path string) {
	// etcd的key用作配置数据读取
	hostKey := strings.Join([]string{"/govern"}, "/")
	// init watch
	watch, err := d.etcdClient.WatchPrefix(context.Background(), hostKey)

	if err != nil {
		panic("watch err: " + err.Error())
	}
	go func() {
		for {
			select {
			case event := <-watch.C():
				switch event.Type {
				case mvccpb.DELETE:
					key := string(event.Kv.Key)
					govern := d.parseGovernKey(key)

					if govern == nil {
						continue
					}

					filename := govern.AppName + "_" + govern.Hostname
					_ = os.Remove(path + "/" + filename + ".yml")
					_ = os.Remove(path + "/" + "pyroscope/" + filename + ".yml")
				case mvccpb.PUT:
					key, value := string(event.Kv.Key), string(event.Kv.Value)
					govern := d.parseGovern(key, value)

					if d.filter(govern) == nil {
						continue
					}

					err = d.writeFile(path, govern)
					if err != nil {
						xlog.Error("writeFile error", xlog.FieldErr(err))
						continue
					}
				}
			}
		}
	}()
}

// GovernConfigScanner ..
func (d *DataSource) GovernConfigScanner(path string) {
	// etcd的key用作配置数据读取
	hostKey := strings.Join([]string{"/govern"}, "/")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	resp, err := d.etcdClient.Get(ctx, hostKey, clientv3.WithPrefix())
	if err != nil {
		xlog.Error("etcdClient.Get error", xlog.String("plugin", "confgo"), xlog.String("msg", err.Error()), xlog.String("hostKey", hostKey))
		return
	}

	xlog.Info("GovernConfigScanner begin")

	for _, kv := range resp.Kvs {
		key, value := string(kv.Key), string(kv.Value)
		govern := d.parseGovern(key, value)

		if d.filter(govern) == nil {
			continue
		}

		err = d.writeFile(path, govern)
		if err != nil {
			xlog.Error("writeFile error", xlog.FieldErr(err))
			continue
		}
	}
	return
}

// cleanup clean invalid prometheus yml
func (d *DataSource) cleanup(prometheusDir string, timeInterval uint32) {
	var (
		once sync.Once
	)
	fileSystem := os.DirFS(prometheusDir)
	cleanFunc := func() {
		err := fs.WalkDir(fileSystem, ".", func(p string, d fs.DirEntry, err error) error {
			depth := strings.Count(prometheusDir, "/") - strings.Count(p, "/")
			if depth != strings.Count(prometheusDir, "/") {
				return nil
			}

			fileInfo, _ := d.Info()
			if fileInfo.IsDir() {
				return nil
			}

			if time.Now().Unix()-fileInfo.ModTime().Unix() > int64(timeInterval) {
				if path.Ext(path.Base(p)) != ".yml" {
					return nil
				}
				return os.Remove(path.Join(prometheusDir, p))
			}
			return nil
		})
		if err != nil {
			xlog.Error("remove file error", xlog.FieldErr(err))
		}
	}
	once.Do(func() {
		cleanFunc()
	})
}
