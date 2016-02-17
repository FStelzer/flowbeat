package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/fstelzer/flowbeat/beater"
)

var Name = "flowbeat"

func main() {
	if err := beat.Run(Name, "", beater.New()); err != nil {
		os.Exit(1)
	}
}
