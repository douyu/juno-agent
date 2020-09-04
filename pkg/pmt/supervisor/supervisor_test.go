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

package supervisor

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var content = `
[program:testing]
directory=/home/www/server/testing
environment=LD_LIBRARY_PATH="/home/www/server/testing/lib",JUNO_TIME="{{JunoTime}}"
command=/home/www/server/%(program_name)s/bin/%(program_name)s --host={{ServerName}} --config=/home/www/.config/piquet/wsd-live-app-home-go/config/config-trunk.toml
autostart=true
autorestart=true
startsecs=10
stdout_logfile=/home/www/server/testing/stdout.log
stdout_logfile_maxbytes=100MB
stdout_logfile_backups=10
stdout_capture_maxbytes=100MB
stderr_logfile=/home/www/server/testing/stderr.log
stderr_logfile_maxbytes=100MB
stderr_logfile_backups=10
stopsignal=QUIT
`

func Test_supervisorScanner(t *testing.T) {
	config := &Config{
		Dir: "../../../../testdata/testdata/supervisor.d",
	}
	scanner := config.Build()
	t.Run("parse conf file", func(t *testing.T) {
		result, err := scanner.parse([]byte(content))
		assert.Nil(t, err)
		t.Log("result", result)
	})

	t.Run("parse conf dir", func(t *testing.T) {
		targets, err := scanner.ListPrograms()
		assert.Nil(t, err)
		t.Log("targets", targets)
		assert.NotNil(t, targets["test1"])
		assert.NotNil(t, targets["test2"])
		for key, val := range targets {
			t.Logf("==> %+v = %+v\n", key, val)
		}
		t.Logf("===> len(targets) = %d\n", len(targets))
	})

	t.Run("parse ini file", func(t *testing.T) {
		task, _, err := scanner.parseFile("../../../testdata/testdata/supervisor.d/test1.conf")
		assert.Nil(t, err)
		fmt.Printf("task = %+v\n", task)
	})
}
