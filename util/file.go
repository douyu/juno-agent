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

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path"
)

func ReadDirFiles(dir string, appName string) (result map[string]string, err error) {
	result = make(map[string]string)
	if files, readerr := ioutil.ReadDir(dir); readerr != nil {
		err = readerr
		return
	} else {
		if appName != "" {
			newFiles := make([]os.FileInfo, 0)
			for _, info := range files {
				if info.Name() == appName {
					newFiles = append(newFiles, info)
					break
				}
			}
			files = newFiles
		}
		for _, fileInfo := range files {
			if !fileInfo.IsDir() {
				fileContent, readerr := ioutil.ReadFile(dir + "/" + fileInfo.Name())
				if readerr != nil {
					err = readerr
					return
				}
				result[fileInfo.Name()] = string(fileContent)
			}
		}
	}
	return result, nil
}

func MkdirAll(filePath string) error {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		dir := path.Dir(filePath)
		if mkdErr := os.MkdirAll(dir, os.ModePerm); mkdErr != nil {
			return mkdErr
		}
	} else {
		return err
	}
	return nil
}

func WriteFile(filePath string, content string) error {
	if err := MkdirAll(filePath); err != nil {
		return err
	}
	if err := ioutil.WriteFile(filePath, []byte(content), os.ModePerm); err != nil {
		return err
	}
	return nil
}

// 生成32位MD5
func MD5(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}
