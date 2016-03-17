package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/builderscon/octav/octav"
	"github.com/builderscon/octav/octav/db"
	"github.com/jessevdk/go-flags"
)

type options struct {
	Config string `short:"c" long:"config" description:"path to a config file"`
	Listen string `short:"l" long:"listen" default:":8080" description:"Listen address"`
}

type config struct {
	Listen   string    `json:"listen"`
	Database db.Config `json:"database"`
}

func main() { os.Exit(_main()) }
func _main() int {
	var opts options
	if _, err := flags.Parse(&opts); err != nil {
		log.Printf("%s", err)
		return 1
	}

	fn := opts.Config
	if fn == "" {
		// No configuration? do the best we can
		if err := db.Init(""); err != nil {
			log.Printf("%s", err)
			return 1
		}
	} else {
		var c config
		fh, err := os.Open(fn)
		if err != nil {
			log.Printf("%s", err)
			return 1
		}
		if err := json.NewDecoder(fh).Decode(&c); err != nil {
			log.Printf("%s", err)
			return 1
		}

		if err := db.Init(c.Database.DSN); err != nil {
			log.Printf("%s", err)
			return 1
		}
		if l := c.Listen; l != "" {
			opts.Listen = l
		}
	}

	log.Printf("Server listening on %s", opts.Listen)
	if err := octav.Run(opts.Listen); err != nil {
		log.Printf("%s", err)
		return 1
	}
	return 0
}
