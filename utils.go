package main

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"
)

func startprofile(filename string) {
	fmt.Println("Profiling")
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
}
func stopprofile() {
	pprof.StopCPUProfile()
}
