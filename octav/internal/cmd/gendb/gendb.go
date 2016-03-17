// +build !appengine

package main

import (
	"log"
	"os"

	"github.com/jessevdk/go-flags"
)

type options struct {
	Types []string `short:"t" long:"type" description:"name of the target type"`
	Dir   string   `short:"d" long:"dir" required:"true" description:"directory to process"`
}

func main() {
	os.Exit(_main())
}

func _main() int {
	var opts options
	if _, err := flags.Parse(&opts); err != nil {
		log.Printf("%s", err)
		return 1
	}

	p := Processor{}
	p.Types = opts.Types
	p.Dir = opts.Dir
	if err := p.Do(); err != nil {
		log.Printf("%s", err)
		return 1
	}
	return 0
}