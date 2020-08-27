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

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/douyu/juno-agent/pkg/core"
	"github.com/douyu/juno-agent/util"
	"github.com/douyu/jupiter/pkg/flag"
)

func init() {

	flag.Register(
		&flag.BoolFlag{
			Name:    "confProxy",
			Usage:   "use --confProxy=true, show confProxy true or false",
			Default: false,
		},
		&flag.BoolFlag{
			Name:    "regProxy",
			Usage:   "use --regProxy=true, show regProxy plugin true or false",
			Default: false,
		},
		&flag.BoolFlag{
			Name:    "agentReport",
			Usage:   "use --agentReport=true, show agentReport plugin true or false",
			Default: false,
		},
		&flag.BoolFlag{
			Name:    "process",
			Usage:   "use --process=true, show process plugin true or false",
			Default: false,
		},
		&flag.BoolFlag{
			Name:    "supervisor",
			Usage:   "use --supervisor=true, show supervisor plugin true or false",
			Default: false,
		},
		&flag.BoolFlag{
			Name:    "systemd",
			Usage:   "use --systemd=true, show systemd plugin true or false",
			Default: false,
		},
		&flag.BoolFlag{
			Name:    "nginx",
			Usage:   "use --nginx=true, show nginx plugin true or false",
			Default: false,
		},
		&flag.BoolFlag{
			Name:    "worker",
			Usage:   "use --worker=true, show worker plugin true or false",
			Default: false,
		},
	)

}

func main() {
	args := os.Args
	if len(args) != 0 && len(args) >= 2 {
		if args[1] == "config" {
			fmt.Println(util.DefaultConfig)
			return
		}
	}
	eng := core.NewEngine()
	//eng.SetGovernor("127.0.0.1:9099")

	if err := eng.Run(); err != nil {
		log.Fatal(err)
	}
}
