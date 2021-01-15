package helpers

import (
	"io"
	"log"
	"os"
)

func CloseIt(body io.ReadCloser) {
	if body != nil {
		err := body.Close()
		CheckError(err)
	}
}

func CheckError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func CheckErrorWithResult(something interface{}, err error) interface{} {
	CheckError(err)
	return something
}

func FileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		return false
	}
}
