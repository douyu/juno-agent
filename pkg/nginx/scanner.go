// Copyright 2020 Douyu
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package nginx

import (
	"github.com/douyu/juno-agent/pkg/structs"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	parser "github.com/yangchenxing/go-nginx-conf-parser"
)

// Config nginx config
type ConfigBlock = parser.NginxConfigureBlock

// ConfScanner ...
type ConfScanner struct {
	enable     bool
	confDir    string
	stop       chan struct{}
	chanConfig chan *structs.NginxConfExt
	watchPath  []string
}

// NewScanner  ...
func NewScanner(confDir string, enable bool) *ConfScanner {
	return &ConfScanner{
		enable:     enable,
		stop:       make(chan struct{}),
		confDir:    confDir,
		chanConfig: make(chan *structs.NginxConfExt, 128),
		watchPath:  make([]string, 0),
	}
}

// Start ...
func (c *ConfScanner) Start() error {
	if c.enable {
		return c.watch()
	}
	return nil
}

// Close ...
func (c *ConfScanner) Close() error {
	c.stop <- struct{}{}
	close(c.chanConfig)
	return nil
}

// Scan ...
func (c *ConfScanner) Scan() ([]*structs.NginxConfExt, error) {
	var configs = make([]*structs.NginxConfExt, 0)
	if c.enable {
		fileList := make([]string, 0)
		watchPath := make([]string, 0)
		walkFunc := func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				if filepath.Ext(info.Name()) == ".conf" {
					fileList = append(fileList, path)
				}
				return nil
			}
			watchPath = append(watchPath, path)
			return nil
		}
		if err := filepath.Walk(c.confDir, walkFunc); err != nil {
			return nil, err
		}

		for _, info := range fileList {
			config, err := c.parseFile(info)
			if err != nil {
				xlog.Error("nginx.parse", xlog.String("err", err.Error()), xlog.String("path", info))
				continue
			}
			configs = append(configs, &structs.NginxConfExt{
				Name:   info,
				Status: "list",
				Block:  config,
			})
		}
		c.withWatchPath(watchPath...)
		return configs, nil
	}
	return configs, nil
}

// C consume the chan program
func (c *ConfScanner) C() <-chan *structs.NginxConfExt {
	return c.chanConfig
}

// watch ...
func (c *ConfScanner) watch() error {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		xlog.Error("nginx fsnotify new watcher err", xlog.String("msg", err.Error()))
		return err
	}
	for _, dir := range c.watchPath {
		if err := w.Add(dir); err != nil {
			xlog.Error("nginx fsnotify add dir err", xlog.String("msg", err.Error()), xlog.String("path", dir))
			return nil
		}
	}
	go func() {
		for {
			select {
			case ev, ok := <-w.Events:
				if !ok {
					return
				}
				if !strings.HasSuffix(ev.Name, ".conf") {
					continue
				}
				xlog.Debug("nginxConfigChange", xlog.String("events", ev.String()))
				switch ev.Op {
				case fsnotify.Create:
					block, err := c.parseFile(ev.Name)
					if err != nil {
						xlog.Error("set status err", xlog.String("msg", err.Error()))
						continue
					}
					c.chanConfig <- &structs.NginxConfExt{
						Status: "create",
						Name:   ev.Name,
						Block:  block,
					}
				case fsnotify.Write:
					block, err := c.parseFile(ev.Name)
					if err != nil {
						xlog.Error("set status err", xlog.String("msg", err.Error()))
						continue
					}
					c.chanConfig <- &structs.NginxConfExt{
						Status: "update",
						Name:   ev.Name,
						Block:  block,
					}
				case fsnotify.Remove:
					c.chanConfig <- &structs.NginxConfExt{
						Status: "delete",
						Name:   ev.Name,
						Block:  ConfigBlock{},
					}
				}
			case err, ok := <-w.Errors:
				xlog.Error("watch err", xlog.String("msg", err.Error()))
				if !ok {
					return
				}
			case <-c.stop:
				return
			}
		}
	}()
	return nil
}

// parseFile ...
func (c *ConfScanner) parseFile(filename string) (ConfigBlock, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return c.parse(content)
}

// parse ...
func (c *ConfScanner) parse(content []byte) (ConfigBlock, error) {
	return parser.Parse(content)
}

func (c *ConfScanner) withWatchPath(path ...string) {
	for _, info := range path {
		c.watchPath = append(c.watchPath, info)
	}
}
