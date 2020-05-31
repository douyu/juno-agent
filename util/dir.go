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


package util

//
// func WatchDir() error {
// 	w, err := fsnotify.NewWatcher()
// 	if err != nil {
// 		return  err
// 	}
//
// 	go func() {
// 	Loop:
// 		for {
// 			select {
// 			case ev, ok := <-w.Events:
// 				if !ok {
// 					return
// 				}
// 				log.Debug("supervisorConfigChange", "events", ev.String())
//
// 				if ev.Op == fsnotify.Create && strings.HasSuffix(ev.Name, ".conf") { // 新增
// 					filePath := ev.Name
// 					log.Debug("supervisorCreate", "filePath", filePath)
//
// 					buf, err := ioutil.ReadFile(filePath)
// 					if err != nil {
// 						log.Error("read file err", "msg", err.Error())
// 						return
// 					}
// 					fileName := filepath.Base(filePath)
// 					if err := c.setStatus(fileName, string(buf)); err != nil {
// 						log.Error("set status err", "msg", err.Error())
// 					}
// 					log.Debug("supervisorSetStatus", "fileName", fileName, "err", err)
//
// 				} else if ev.Op == fsnotify.Write && strings.HasSuffix(ev.Name, ".conf") { // 修改
// 					filePath := ev.Name
//
// 					log.Debug("supervisorWrite", "filePath", filePath)
//
// 					buf, err := ioutil.ReadFile(filePath)
// 					if err != nil {
// 						log.Error("read file err", "msg", err.Error())
// 						return
// 					}
// 					fileName := filepath.Base(filePath)
// 					if err := c.setStatus(fileName, string(buf)); err != nil {
// 						log.Error("set status err", "msg", err.Error())
// 					}
// 					log.Debug("supervisorSetStatus", "fileName", fileName, "err", err)
//
// 				} else if ev.Op == fsnotify.Remove && strings.HasSuffix(ev.Name, ".conf") { // 删除
//
// 					fileName := filepath.Base(ev.Name)
// 					appName, err := getAppName(fileName)
//
// 					log.Debug("supervisorRemove", "fileName", fileName, "appName", appName, "err", err)
//
// 					if err == nil {
// 						c.StatusMap.Delete(appName)
// 					}
// 				}
// 			case err, ok := <-w.Errors:
// 				log.Error("watch err", "msg", err.Error())
// 				if !ok {
// 					return
// 				}
// 			case <-c.stop:
// 				break Loop
// 			}
// 		}
// 	}()
// 	return nil
// }
