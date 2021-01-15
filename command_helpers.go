package helpers

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

const DEBUG = false

func Execute(cmd string) ([]byte, error) {
	if DEBUG {
		log.Println(cmd)
	}
	command := exec.Command("bash", "-c", cmd)
	command.Stderr = nil
	return command.Output()
}

func DirectExecution(cmd string, args string) ([]byte, error) {
	if DEBUG {
		log.Println(cmd, args)
	}
	command := exec.Command(cmd, args)
	command.Stderr = nil
	return command.Output()
}

func ExecuteOrPanic(cmd string) string {
	out, err := Execute(cmd)
	if err != nil {
		log.Println(fmt.Sprintf("Error while running cmd: bash -c '%s', output: %s", cmd, out))
		log.Panic(err)
	}
	return strings.TrimRight(string(out), " \r\n")
}

func ExecuteTillSuccessful(cmd string, max int, panicIfUnsuccessful bool) string {
	var results []string
	for i := 0; i < max; i++ {
		result, err := Execute(cmd)
		if err == nil {
			return string(result)
		}
		time.Sleep(time.Second * 1)
		results = append(results, string(result))
	}
	if panicIfUnsuccessful {
		errs := strings.Join(results, "; ")
		log.Panic(fmt.Errorf("failed to execute cmd: %s after %d attempts; Errors: %s", cmd, max, errs))
	}
	return ""
}
