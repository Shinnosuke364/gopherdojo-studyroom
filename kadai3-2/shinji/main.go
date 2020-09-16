package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/shinnosuke/gopherdojo-studyroom/kadai3-2/shinji/spdl"
)

var (
	url      *string = flag.String("u", "https://golangcode.com/logo.svg", "download source URL")
	filepath *string = flag.String("f", "./logo.svg", "filepath of the destination")
	procs    *int    = flag.Int("p", 5, "number of processes")
)

func main() {
	flag.Parse()

	start := time.Now()

	err := spdl.Download(*filepath, *url, *procs)
	if err != nil {
		log.Fatal(err)
	}

	end := time.Now()

	fmt.Printf("File %s downlaod in current working directory\n", *filepath)
	fmt.Printf("%vprocks: %f秒\n", *procs, (end.Sub(start)).Seconds())
}
