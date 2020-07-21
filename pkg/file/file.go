package file

import (
	"fmt"
	"io/ioutil"
	"os"
)

func ReadFile(path string) (content string, err error) {
	stat, err := os.Stat(path)
	if err != nil {
		return
	}

	if stat.IsDir() {
		return "", fmt.Errorf(path + " is directory, not file")
	}

	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	content = string(fileBytes)

	return
}
