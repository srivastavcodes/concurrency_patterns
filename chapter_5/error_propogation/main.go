package main

import (
	"errors"
	"fmt"
	"log"
	"os"
)

func handleError(key int, err error, message string) {
	log.SetPrefix(fmt.Sprintf("[logID: %v]: ", key))
	log.Printf("%#v", err)
	log.Printf("[%v] %v", key, message)
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	err := runJob("1")
	if err != nil {
		msg := "there was an unexpected issue; please report this as bug"
		if errors.As(err, &IntermediateLevelErr{}) {
			msg = err.Error()
		}
		handleError(1, err, msg)
	}
}
