package utils

import (
	"log"
	"runtime"
)

func Check(err error) {
	if err != nil {
		log.Printf("\nError: %+v", err)
	}
}

func CheckFatal(err error) {
	if err != nil {
		log.Printf("\nError: %+v", err)
		buf := make([]byte, 1<<20)
		stacklen := runtime.Stack(buf, true)
		log.Fatalf("=== received Fatal ERROR ===\n*** goroutine dump...\n%s\n*** end\n", buf[:stacklen])
	}
}
