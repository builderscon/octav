package main

import (
	"log"
	"os"

	"github.com/builderscon/octav/octav"
	"github.com/jessevdk/go-flags"
)

type options struct {
	Listen string `short:"l" long:"listen" default:":8080" description:"Listen address"`
}

func main() { os.Exit(_main()) }
func _main() int {
	var opts options
	if _, err := flags.Parse(&opts); err != nil {
		log.Printf("%s", err)
		return 1
	}
	log.Printf("Server listening on %s", opts.Listen)
	if err := octav.Run(opts.Listen); err != nil {
		log.Printf("%s", err)
		return 1
	}
	return 0
}
