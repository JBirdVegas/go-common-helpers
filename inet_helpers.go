package helpers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func isConnected() (ok bool) {
	_, err := http.Head("https://clients3.google.com/generate_204")
	if err != nil {
		return false
	}
	return true
}

func WaitTillConnected(maxSeconds int) {
	maxDuration := time.Second * time.Duration(maxSeconds)
	start := time.Now()
	index := 1

	for ok := true; ok; ok = !isConnected() {
		index += 1
		time.Sleep(time.Millisecond * 25)
		now := time.Now()
		sub := start.Sub(now)
		if sub >= maxDuration {
			log.Println(fmt.Sprintf("Waited for %d before giving up", sub))
			os.Exit(4)
		}
	}
	result := time.Since(start) / time.Second
	if result > 1 {
		log.Println(fmt.Sprintf("Waited for %d seconds before internet connection was detectected", result))
	}
}
