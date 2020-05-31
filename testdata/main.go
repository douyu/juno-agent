package main

import (
	"fmt"
	"strings"
)

func main() {

	result := `root 2750175 0.6 0.1 1640248 77596 ? Sl 7月04 123:09 /home/www/server/wsd-penaten-admin-juno-go/bin/wsd-penaten-admin-juno-go --host=10.1.56.26 --config=config/config-prod.toml
www 2750176 1.5 0.0 617680 34476 ? Sl 7月04 276:47 /home/www/server/wsd-penaten-admin-generatorcode-go/bin/wsd-penaten-admin-generatorcode-go --host=10.1.56.26 --config=config/config-prod.toml`
	arr := strings.Split(result, "\n")
	for _, line := range arr {
		items := strings.SplitAfterN(line, " ", 11)
		fmt.Println(len(items))
		fmt.Println(items[10])
	}
}
