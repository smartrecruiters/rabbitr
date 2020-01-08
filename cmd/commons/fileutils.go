package commons

import (
	"os"
)

func MakeDir(path string) (err error) {
	err = nil
	if _, err = os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			Debugf("Make dir error occurred, err: %s", err)
		}
	}
	return
}
