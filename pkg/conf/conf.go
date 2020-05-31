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

package conf

import (
	"log"
	"time"

	"github.com/douyu/jupiter/pkg/util/xgo"
	"github.com/fsnotify/fsnotify"
)

// AppConf ...
type AppConf map[string]interface{}

// AppConfScanner ...
type AppConfScanner struct {
	dirPath     string
	stop        chan struct{}
	chanAppConf chan AppConf
}

// NewScanner new app config scanner
func NewScanner(dirPath string) *AppConfScanner {
	return &AppConfScanner{
		dirPath:     dirPath,
		stop:        make(chan struct{}),
		chanAppConf: make(chan AppConf, 100),
	}
}

// Start ...
func (sc *AppConfScanner) Start() error {
	xgo.DelayGo(time.Second*5, sc.watch)
	return nil
}

// Close ...
func (sc *AppConfScanner) Close() error {
	sc.stop <- struct{}{}
	return nil
}

// C ...
func (sc *AppConfScanner) C() <-chan AppConf {
	return sc.chanAppConf
}

// Scan ...
func (sc *AppConfScanner) Scan() {
}

// watch ...
func (sc *AppConfScanner) watch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()

	err = watcher.Add(sc.dirPath)
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			log.Println("event:", event)
			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Println("modified file:", event.Name)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		case <-sc.stop:
			return
		}
	}
}
