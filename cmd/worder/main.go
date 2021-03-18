package main

import (
	"flag"
	"io/ioutil"
	"log"
	"time"

	"github.com/i7tsov/worder/pkg/worder"
)

func main() {
	file := flag.String("file", "", "file name to parse")
	workers := flag.Int("workers", 1, "number of parallel worker routines")
	flag.Parse()

	if *file == "" {
		log.Fatalf("Please provide file name")
	}

	pl, err := ioutil.ReadFile(*file)
	if err != nil {
		log.Fatal(err)
	}

	w := worder.Worder{
		Text:    string(pl),
		Workers: *workers,
		Path:    "./results",
	}

	start := time.Now()
	w.Run()
	elapsed := time.Now().Sub(start)

	log.Printf("Done. Elapsed time: %v", elapsed)
}
