package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/shinnosuke/gopherdojo-studyroom/kadai3-2/shinji/spdl"
)

var (
	url      *string = flag.String("u", "https://download.docker.com/win/stable/Docker%20Desktop%20Installer.exe", "download source URL")
	filepath *string = flag.String("f", "./docker.exe", "filepath of the destination")
	procs    *int    = flag.Int("p", 16, "number of processes")
)

func main() {
	flag.Parse()

	start := time.Now()

	err := spdl.Download(*filepath, *url, *procs)
	if err != nil {
		log.Fatal(err)
	}

	end := time.Now()

	fmt.Printf("\nFile %s downlaod in current working directory\n", *filepath)
	fmt.Printf("%vprocks: %f秒\n\n", *procs, (end.Sub(start)).Seconds())
}
