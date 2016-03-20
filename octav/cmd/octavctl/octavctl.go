package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/builderscon/octav/octav/client"
	"github.com/builderscon/octav/octav/model"
	"github.com/builderscon/octav/octav/validator"
)

// TODO: We should consider auto-generating this whole file,
// but I'm not sure if we can do it from the JSON Hyper Schema file :/

// Global options
var endpoint string

func prepGlobalFlags(fs *flag.FlagSet) {
	fs.StringVar(&endpoint, "endpoint", "", "Base URL of the octav server (required)")
}

// Special case for l10n strings, where users need say:
//   -l10n "title#ja=ハロー、ワールド" -l10n "sub_title#ja=サブタイトル"
type l10nvars map[string]string

func (l l10nvars) String() string {
	ks := make([]string, len(l))
	for k := range l {
		ks = append(ks, k)
	}
	sort.Strings(ks)

	buf := bytes.Buffer{}
	for _, k := range ks {
		v := l[k]
		buf.WriteString(k)
		buf.WriteByte('=')
		buf.WriteString(v)
	}
	return buf.String()
}

var errl10nvarfmt = errors.New("value must be in key#lang=value format")

func (l *l10nvars) Set(v string) error {
	eqloc := strings.IndexByte(v, '=')
	if eqloc == -1 || eqloc == len(v)-1 {
		return errl10nvarfmt
	}

	key := v[:eqloc]
	value := v[eqloc+1:]

	lbloc := strings.IndexByte(key, '#')
	if lbloc == -1 {
		return errl10nvarfmt
	}

	(*l)[key] = value
	return nil
}

func main() {
	os.Exit(_main())
}

func newClient() *client.Client {
	return client.New(endpoint)
}

type cmdargs []string

func (a cmdargs) WithFrontPopped() cmdargs {
	if len(a) > 0 {
		return cmdargs(a[1:])
	}
	return nil
}

func (a cmdargs) Get(i int) string {
	if i > -1 && len(a) > i {
		return a[i]
	}
	return ""
}

func (a cmdargs) Len() int {
	return len(a)
}

func _main() int {
	args := cmdargs(os.Args).WithFrontPopped()
	switch args.Get(0) {
	case "conference":
		return doConferenceSubcmd(args.WithFrontPopped())
	default:
		log.Printf("unimplemented (main)")
		return 1
	}
}

func printJSON(v interface{}) error {
	buf, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	os.Stdout.Write(buf)
	os.Stdout.Write([]byte{'\n'})
	return nil
}

func doConferenceSubcmd(args cmdargs) int {
	switch v := args.Get(0); v {
	case "create":
		return doConferenceCreate(args.WithFrontPopped())
	case "lookup":
		return doConferenceLookup(args.WithFrontPopped())
	case "dates":
		return doConferenceDatesSubcmd(args.WithFrontPopped())
	default:
		log.Printf("unimplemented (conference): %s", v)
		return 1
	}
}

func doConferenceCreate(args cmdargs) int {
	l10n := l10nvars{}
	fs := flag.NewFlagSet("octavctl conference create", flag.ContinueOnError)
	fs.Var(&l10n, "l10n", "localized strings")
	prepGlobalFlags(fs)
	if err := fs.Parse([]string(args)); err != nil {
		log.Printf("%s", err)
		return 1
	}

	log.Printf("%#v", l10n)
	/*
		cl := newClient()
		conf, err := cl.LookupConference(&model.LookupConferenceRequest{
			ID: id,
		})
		if err != nil {
			log.Printf("%s", err)
			return 1
		}

		if err := printJSON(conf); err != nil {
			log.Printf("%s", err)
			return 1
		}
	*/
	return 0
}

func doConferenceLookup(args cmdargs) int {
	var id string
	fs := flag.NewFlagSet("octavctl conference lookup", flag.ContinueOnError)
	fs.StringVar(&id, "id", "", "ID of the conference to lookup")
	prepGlobalFlags(fs)
	if err := fs.Parse([]string(args)); err != nil {
		log.Printf("%s", err)
		return 1
	}

	m := make(map[string]interface{})
	m["id"] = id
	r := model.LookupConferenceRequest{}
	if err := r.Populate(m); err != nil {
		return errOut(err)
	}
	if err := validator.HTTPLookupConferenceRequest.Validate(r); err != nil {
		return errOut(err)
	}

	cl := newClient()
	conf, err := cl.LookupConference(&r)
	if err != nil {
		return errOut(err)
	}

	if err := printJSON(conf); err != nil {
		return errOut(err)
	}

	return 0
}

func doConferenceDatesSubcmd(args cmdargs) int {
	switch args.Get(0) {
	case "add":
		return doConferenceDatesAdd(args.WithFrontPopped())
	case "delete":
		return doConferenceDatesDelete(args.WithFrontPopped())
	default:
		log.Printf("unimplemented (conference dates)")
	}

	return 1
}

type stringList []string
func (l *stringList) Set(s string) error {
	*l = append(*l, s)
	return nil
}

func (l *stringList) String() string {
	buf := bytes.Buffer{}
	for i, v := range *l {
		buf.WriteString(v)
		if i != len(*l) - 1 {
			buf.WriteByte(' ')
		}
	}
	return buf.String()
}

func (l stringList) List() []string {
	return []string(l)
}

func doConferenceDatesAdd(args cmdargs) int {
	var id string
	fs := flag.NewFlagSet("octavctl conference add dates", flag.ContinueOnError)
	fs.StringVar(&id, "id", "", "ID of the target conference")
	var dtlist stringList
	fs.Var(&dtlist, "date", "Date(s) to add (may be repeated)")
	prepGlobalFlags(fs)
	if err := fs.Parse([]string(args)); err != nil {
		log.Printf("%s", err)
		return 1
	}

	m := make(map[string]interface{})
	m["conference_id"] = id
	m["dates"] = dtlist.List()
	r := model.AddConferenceDatesRequest{}
	if err := r.Populate(m); err != nil {
		return errOut(err)
	}

	if err := validator.HTTPAddConferenceDatesRequest.Validate(r); err != nil {
		return errOut(err)
	}
	cl := newClient()
	if err := cl.AddConferenceDates(&r); err != nil {
		return errOut(err)
	}

	return 0
}

func errOut(err error) int {
	log.Printf("%s", err)
	return 1
}

func doConferenceDatesDelete(args cmdargs) int {
	var id string
	fs := flag.NewFlagSet("octavctl conference delete dates", flag.ContinueOnError)
	fs.StringVar(&id, "id", "", "ID of the target conference")
	var dtlist stringList
	fs.Var(&dtlist, "date", "Date(s) to delete (may be repeated)")
	prepGlobalFlags(fs)
	if err := fs.Parse([]string(args)); err != nil {
		return errOut(err)
	}

	m := make(map[string]interface{})
	m["conference_id"] = id
	m["dates"] = dtlist.List()

	r := model.DeleteConferenceDatesRequest{}
	if err := r.Populate(m); err != nil {
		return errOut(err)
	}

log.Printf("%#v", r)

	if err := validator.HTTPDeleteConferenceDatesRequest.Validate(r); err != nil {
		return errOut(err)
	}

	cl := newClient()
	if err := cl.DeleteConferenceDates(&r); err != nil {
		return errOut(err)
	}

	return 0
}
