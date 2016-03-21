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

func newClient() (*client.Client, error) {
	if endpoint == "" {
		return nil, errors.New("-endpoint is required")
	}
	return client.New(endpoint), nil
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

type stringList []string

func (l *stringList) Set(s string) error {
	*l = append(*l, s)
	return nil
}

func (l *stringList) String() string {
	buf := bytes.Buffer{}
	for i, v := range *l {
		buf.WriteString(v)
		if i != len(*l)-1 {
			buf.WriteByte(' ')
		}
	}
	return buf.String()
}

func (l stringList) Valid() bool {
	return len(l) > 0
}

func (l stringList) Get() interface{} {
	return []string(l)
}

func main() {
	args := cmdargs(os.Args).WithFrontPopped()
	os.Exit(doSubcmd(args))
}

func errOut(err error) int {
	log.Printf("%s", err)
	return 1
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

func doConferenceCreate(args cmdargs) int {
	fs := flag.NewFlagSet("octavctl conference create", flag.ContinueOnError)
	var description string
	fs.StringVar(&description, "description", "", "")
	var slug string
	fs.StringVar(&slug, "slug", "", "")
	var sub_title string
	fs.StringVar(&sub_title, "sub_title", "", "")
	var title string
	fs.StringVar(&title, "title", "", "")
	var user_id string
	fs.StringVar(&user_id, "user_id", "", "")
	prepGlobalFlags(fs)
	if err := fs.Parse([]string(args)); err != nil {
		return errOut(err)
	}

	m := make(map[string]interface{})
	if description != "" {
		m["description"] = description
	}
	if slug != "" {
		m["slug"] = slug
	}
	if sub_title != "" {
		m["sub_title"] = sub_title
	}
	if title != "" {
		m["title"] = title
	}
	if user_id != "" {
		m["user_id"] = user_id
	}
	r := model.CreateConferenceRequest{}
	if err := r.Populate(m); err != nil {
		return errOut(err)
	}

	if err := validator.HTTPCreateConferenceRequest.Validate(&r); err != nil {
		return errOut(err)
	}

	cl, err := newClient()
	if err != nil {
		return errOut(err)
	}
	res, err := cl.CreateConference(&r)
	if err != nil {
		return errOut(err)
	}
	if err := printJSON(res); err != nil {
		return errOut(err)
	}

	return 0
}

func doConferenceLookup(args cmdargs) int {
	fs := flag.NewFlagSet("octavctl conference lookup", flag.ContinueOnError)
	var id string
	fs.StringVar(&id, "id", "", "")
	prepGlobalFlags(fs)
	if err := fs.Parse([]string(args)); err != nil {
		return errOut(err)
	}

	m := make(map[string]interface{})
	if id != "" {
		m["id"] = id
	}
	r := model.LookupConferenceRequest{}
	if err := r.Populate(m); err != nil {
		return errOut(err)
	}

	if err := validator.HTTPLookupConferenceRequest.Validate(&r); err != nil {
		return errOut(err)
	}

	cl, err := newClient()
	if err != nil {
		return errOut(err)
	}
	res, err := cl.LookupConference(&r)
	if err != nil {
		return errOut(err)
	}
	if err := printJSON(res); err != nil {
		return errOut(err)
	}

	return 0
}

func doConferenceDelete(args cmdargs) int {
	fs := flag.NewFlagSet("octavctl conference delete", flag.ContinueOnError)
	var id string
	fs.StringVar(&id, "id", "", "")
	prepGlobalFlags(fs)
	if err := fs.Parse([]string(args)); err != nil {
		return errOut(err)
	}

	m := make(map[string]interface{})
	if id != "" {
		m["id"] = id
	}
	r := model.DeleteConferenceRequest{}
	if err := r.Populate(m); err != nil {
		return errOut(err)
	}

	if err := validator.HTTPDeleteConferenceRequest.Validate(&r); err != nil {
		return errOut(err)
	}

	cl, err := newClient()
	if err != nil {
		return errOut(err)
	}
	if err := cl.DeleteConference(&r); err != nil {
		return errOut(err)
	}

	return 0
}

func doConferenceList(args cmdargs) int {
	fs := flag.NewFlagSet("octavctl conference list", flag.ContinueOnError)
	var lang string
	fs.StringVar(&lang, "lang", "", "")
	var limit int64
	fs.Int64Var(&limit, "limit", 0, "")
	var range_end string
	fs.StringVar(&range_end, "range_end", "", "")
	var range_start string
	fs.StringVar(&range_start, "range_start", "", "")
	var since string
	fs.StringVar(&since, "since", "", "")
	prepGlobalFlags(fs)
	if err := fs.Parse([]string(args)); err != nil {
		return errOut(err)
	}

	m := make(map[string]interface{})
	if lang != "" {
		m["lang"] = lang
	}
	if limit != 0 {
		m["limit"] = limit
	}
	if range_end != "" {
		m["range_end"] = range_end
	}
	if range_start != "" {
		m["range_start"] = range_start
	}
	if since != "" {
		m["since"] = since
	}
	r := model.ListConferenceRequest{}
	if err := r.Populate(m); err != nil {
		return errOut(err)
	}

	if err := validator.HTTPListConferenceRequest.Validate(&r); err != nil {
		return errOut(err)
	}

	cl, err := newClient()
	if err != nil {
		return errOut(err)
	}
	res, err := cl.ListConference(&r)
	if err != nil {
		return errOut(err)
	}
	if err := printJSON(res); err != nil {
		return errOut(err)
	}

	return 0
}

func doConferenceDatesAdd(args cmdargs) int {
	fs := flag.NewFlagSet("octavctl conference dates add", flag.ContinueOnError)
	var id string
	fs.StringVar(&id, "id", "", "")
	var dates stringList
	fs.Var(&dates, "dates", "")
	prepGlobalFlags(fs)
	if err := fs.Parse([]string(args)); err != nil {
		return errOut(err)
	}

	m := make(map[string]interface{})
	if id != "" {
		m["conference_id"] = id
	}
	if dates.Valid() {
		m["dates"] = dates.Get()
	}
	r := model.AddConferenceDatesRequest{}
	if err := r.Populate(m); err != nil {
		return errOut(err)
	}

	if err := validator.HTTPAddConferenceDatesRequest.Validate(&r); err != nil {
		return errOut(err)
	}

	cl, err := newClient()
	if err != nil {
		return errOut(err)
	}
	if err := cl.AddConferenceDates(&r); err != nil {
		return errOut(err)
	}

	return 0
}

func doConferenceDatesDelete(args cmdargs) int {
	fs := flag.NewFlagSet("octavctl conference dates delete", flag.ContinueOnError)
	var id string
	fs.StringVar(&id, "id", "", "")
	var dates stringList
	fs.Var(&dates, "dates", "")
	prepGlobalFlags(fs)
	if err := fs.Parse([]string(args)); err != nil {
		return errOut(err)
	}

	m := make(map[string]interface{})
	if id != "" {
		m["conference_id"] = id
	}
	if dates.Valid() {
		m["dates"] = dates.Get()
	}
	r := model.DeleteConferenceDatesRequest{}
	if err := r.Populate(m); err != nil {
		return errOut(err)
	}

	if err := validator.HTTPDeleteConferenceDatesRequest.Validate(&r); err != nil {
		return errOut(err)
	}

	cl, err := newClient()
	if err != nil {
		return errOut(err)
	}
	if err := cl.DeleteConferenceDates(&r); err != nil {
		return errOut(err)
	}

	return 0
}

func doConferenceDatesSubcmd(args cmdargs) int {
	switch v := args.Get(0); v {
	case "add":
		return doConferenceDatesAdd(args.WithFrontPopped())
	case "delete":
		return doConferenceDatesDelete(args.WithFrontPopped())
	default:
		log.Printf("unimplemented (conference): %s", v)
		return 1
	}
	return 0
}

func doConferenceAdminAdd(args cmdargs) int {
	fs := flag.NewFlagSet("octavctl conference admin add", flag.ContinueOnError)
	var id string
	fs.StringVar(&id, "id", "", "")
	var user_id string
	fs.StringVar(&user_id, "user_id", "", "")
	prepGlobalFlags(fs)
	if err := fs.Parse([]string(args)); err != nil {
		return errOut(err)
	}

	m := make(map[string]interface{})
	if id != "" {
		m["conference_id"] = id
	}
	if user_id != "" {
		m["user_id"] = user_id
	}
	r := model.AddConferenceAdminRequest{}
	if err := r.Populate(m); err != nil {
		return errOut(err)
	}

	if err := validator.HTTPAddConferenceAdminRequest.Validate(&r); err != nil {
		return errOut(err)
	}

	cl, err := newClient()
	if err != nil {
		return errOut(err)
	}
	if err := cl.AddConferenceAdmin(&r); err != nil {
		return errOut(err)
	}

	return 0
}

func doConferenceAdminDelete(args cmdargs) int {
	fs := flag.NewFlagSet("octavctl conference admin delete", flag.ContinueOnError)
	var id string
	fs.StringVar(&id, "id", "", "")
	var user_id string
	fs.StringVar(&user_id, "user_id", "", "")
	prepGlobalFlags(fs)
	if err := fs.Parse([]string(args)); err != nil {
		return errOut(err)
	}

	m := make(map[string]interface{})
	if id != "" {
		m["conference_id"] = id
	}
	if user_id != "" {
		m["user_id"] = user_id
	}
	r := model.DeleteConferenceAdminRequest{}
	if err := r.Populate(m); err != nil {
		return errOut(err)
	}

	if err := validator.HTTPDeleteConferenceAdminRequest.Validate(&r); err != nil {
		return errOut(err)
	}

	cl, err := newClient()
	if err != nil {
		return errOut(err)
	}
	if err := cl.DeleteConferenceAdmin(&r); err != nil {
		return errOut(err)
	}

	return 0
}

func doConferenceAdminSubcmd(args cmdargs) int {
	switch v := args.Get(0); v {
	case "add":
		return doConferenceAdminAdd(args.WithFrontPopped())
	case "delete":
		return doConferenceAdminDelete(args.WithFrontPopped())
	default:
		log.Printf("unimplemented (conference): %s", v)
		return 1
	}
	return 0
}

func doConferenceSubcmd(args cmdargs) int {
	switch v := args.Get(0); v {
	case "create":
		return doConferenceCreate(args.WithFrontPopped())
	case "lookup":
		return doConferenceLookup(args.WithFrontPopped())
	case "delete":
		return doConferenceDelete(args.WithFrontPopped())
	case "list":
		return doConferenceList(args.WithFrontPopped())
	case "dates":
		return doConferenceDatesSubcmd(args.WithFrontPopped())
	case "admin":
		return doConferenceAdminSubcmd(args.WithFrontPopped())
	default:
		log.Printf("unimplemented (conference): %s", v)
		return 1
	}
	return 0
}

func doVenueCreate(args cmdargs) int {
	fs := flag.NewFlagSet("octavctl venue create", flag.ContinueOnError)
	var address string
	fs.StringVar(&address, "address", "", "")
	var latitude float64
	fs.Float64Var(&latitude, "latitude", 0, "")
	var longitude float64
	fs.Float64Var(&longitude, "longitude", 0, "")
	var name string
	fs.StringVar(&name, "name", "", "")
	prepGlobalFlags(fs)
	if err := fs.Parse([]string(args)); err != nil {
		return errOut(err)
	}

	m := make(map[string]interface{})
	if address != "" {
		m["address"] = address
	}
	if latitude != 0 {
		m["latitude"] = latitude
	}
	if longitude != 0 {
		m["longitude"] = longitude
	}
	if name != "" {
		m["name"] = name
	}
	r := model.CreateVenueRequest{}
	if err := r.Populate(m); err != nil {
		return errOut(err)
	}

	if err := validator.HTTPCreateVenueRequest.Validate(&r); err != nil {
		return errOut(err)
	}

	cl, err := newClient()
	if err != nil {
		return errOut(err)
	}
	res, err := cl.CreateVenue(&r)
	if err != nil {
		return errOut(err)
	}
	if err := printJSON(res); err != nil {
		return errOut(err)
	}

	return 0
}

func doVenueList(args cmdargs) int {
	fs := flag.NewFlagSet("octavctl venue list", flag.ContinueOnError)
	var lang string
	fs.StringVar(&lang, "lang", "", "")
	var since string
	fs.StringVar(&since, "since", "", "")
	prepGlobalFlags(fs)
	if err := fs.Parse([]string(args)); err != nil {
		return errOut(err)
	}

	m := make(map[string]interface{})
	if lang != "" {
		m["lang"] = lang
	}
	if since != "" {
		m["since"] = since
	}
	r := model.ListVenueRequest{}
	if err := r.Populate(m); err != nil {
		return errOut(err)
	}

	if err := validator.HTTPListVenueRequest.Validate(&r); err != nil {
		return errOut(err)
	}

	cl, err := newClient()
	if err != nil {
		return errOut(err)
	}
	res, err := cl.ListVenue(&r)
	if err != nil {
		return errOut(err)
	}
	if err := printJSON(res); err != nil {
		return errOut(err)
	}

	return 0
}

func doVenueLookup(args cmdargs) int {
	fs := flag.NewFlagSet("octavctl venue lookup", flag.ContinueOnError)
	var id string
	fs.StringVar(&id, "id", "", "")
	prepGlobalFlags(fs)
	if err := fs.Parse([]string(args)); err != nil {
		return errOut(err)
	}

	m := make(map[string]interface{})
	if id != "" {
		m["id"] = id
	}
	r := model.LookupVenueRequest{}
	if err := r.Populate(m); err != nil {
		return errOut(err)
	}

	if err := validator.HTTPLookupVenueRequest.Validate(&r); err != nil {
		return errOut(err)
	}

	cl, err := newClient()
	if err != nil {
		return errOut(err)
	}
	res, err := cl.LookupVenue(&r)
	if err != nil {
		return errOut(err)
	}
	if err := printJSON(res); err != nil {
		return errOut(err)
	}

	return 0
}

func doVenueDelete(args cmdargs) int {
	fs := flag.NewFlagSet("octavctl venue delete", flag.ContinueOnError)
	var id string
	fs.StringVar(&id, "id", "", "")
	prepGlobalFlags(fs)
	if err := fs.Parse([]string(args)); err != nil {
		return errOut(err)
	}

	m := make(map[string]interface{})
	if id != "" {
		m["id"] = id
	}
	r := model.DeleteVenueRequest{}
	if err := r.Populate(m); err != nil {
		return errOut(err)
	}

	if err := validator.HTTPDeleteVenueRequest.Validate(&r); err != nil {
		return errOut(err)
	}

	cl, err := newClient()
	if err != nil {
		return errOut(err)
	}
	if err := cl.DeleteVenue(&r); err != nil {
		return errOut(err)
	}

	return 0
}

func doVenueSubcmd(args cmdargs) int {
	switch v := args.Get(0); v {
	case "create":
		return doVenueCreate(args.WithFrontPopped())
	case "list":
		return doVenueList(args.WithFrontPopped())
	case "lookup":
		return doVenueLookup(args.WithFrontPopped())
	case "delete":
		return doVenueDelete(args.WithFrontPopped())
	default:
		log.Printf("unimplemented (conference): %s", v)
		return 1
	}
	return 0
}

func doRoomCreate(args cmdargs) int {
	fs := flag.NewFlagSet("octavctl room create", flag.ContinueOnError)
	var capacity int64
	fs.Int64Var(&capacity, "capacity", 0, "")
	var name string
	fs.StringVar(&name, "name", "", "")
	var venue_id string
	fs.StringVar(&venue_id, "venue_id", "", "")
	prepGlobalFlags(fs)
	if err := fs.Parse([]string(args)); err != nil {
		return errOut(err)
	}

	m := make(map[string]interface{})
	if capacity != 0 {
		m["capacity"] = capacity
	}
	if name != "" {
		m["name"] = name
	}
	if venue_id != "" {
		m["venue_id"] = venue_id
	}
	r := model.CreateRoomRequest{}
	if err := r.Populate(m); err != nil {
		return errOut(err)
	}

	if err := validator.HTTPCreateRoomRequest.Validate(&r); err != nil {
		return errOut(err)
	}

	cl, err := newClient()
	if err != nil {
		return errOut(err)
	}
	res, err := cl.CreateRoom(&r)
	if err != nil {
		return errOut(err)
	}
	if err := printJSON(res); err != nil {
		return errOut(err)
	}

	return 0
}

func doRoomList(args cmdargs) int {
	fs := flag.NewFlagSet("octavctl room list", flag.ContinueOnError)
	var lang string
	fs.StringVar(&lang, "lang", "", "")
	var limit int64
	fs.Int64Var(&limit, "limit", 0, "")
	var venue_id string
	fs.StringVar(&venue_id, "venue_id", "", "")
	prepGlobalFlags(fs)
	if err := fs.Parse([]string(args)); err != nil {
		return errOut(err)
	}

	m := make(map[string]interface{})
	if lang != "" {
		m["lang"] = lang
	}
	if limit != 0 {
		m["limit"] = limit
	}
	if venue_id != "" {
		m["venue_id"] = venue_id
	}
	r := model.ListRoomRequest{}
	if err := r.Populate(m); err != nil {
		return errOut(err)
	}

	if err := validator.HTTPListRoomRequest.Validate(&r); err != nil {
		return errOut(err)
	}

	cl, err := newClient()
	if err != nil {
		return errOut(err)
	}
	res, err := cl.ListRoom(&r)
	if err != nil {
		return errOut(err)
	}
	if err := printJSON(res); err != nil {
		return errOut(err)
	}

	return 0
}

func doRoomLookup(args cmdargs) int {
	fs := flag.NewFlagSet("octavctl room lookup", flag.ContinueOnError)
	var id string
	fs.StringVar(&id, "id", "", "")
	prepGlobalFlags(fs)
	if err := fs.Parse([]string(args)); err != nil {
		return errOut(err)
	}

	m := make(map[string]interface{})
	if id != "" {
		m["id"] = id
	}
	r := model.LookupRoomRequest{}
	if err := r.Populate(m); err != nil {
		return errOut(err)
	}

	if err := validator.HTTPLookupRoomRequest.Validate(&r); err != nil {
		return errOut(err)
	}

	cl, err := newClient()
	if err != nil {
		return errOut(err)
	}
	res, err := cl.LookupRoom(&r)
	if err != nil {
		return errOut(err)
	}
	if err := printJSON(res); err != nil {
		return errOut(err)
	}

	return 0
}

func doRoomDelete(args cmdargs) int {
	fs := flag.NewFlagSet("octavctl room delete", flag.ContinueOnError)
	var id string
	fs.StringVar(&id, "id", "", "")
	prepGlobalFlags(fs)
	if err := fs.Parse([]string(args)); err != nil {
		return errOut(err)
	}

	m := make(map[string]interface{})
	if id != "" {
		m["id"] = id
	}
	r := model.DeleteRoomRequest{}
	if err := r.Populate(m); err != nil {
		return errOut(err)
	}

	if err := validator.HTTPDeleteRoomRequest.Validate(&r); err != nil {
		return errOut(err)
	}

	cl, err := newClient()
	if err != nil {
		return errOut(err)
	}
	if err := cl.DeleteRoom(&r); err != nil {
		return errOut(err)
	}

	return 0
}

func doRoomSubcmd(args cmdargs) int {
	switch v := args.Get(0); v {
	case "create":
		return doRoomCreate(args.WithFrontPopped())
	case "list":
		return doRoomList(args.WithFrontPopped())
	case "lookup":
		return doRoomLookup(args.WithFrontPopped())
	case "delete":
		return doRoomDelete(args.WithFrontPopped())
	default:
		log.Printf("unimplemented (conference): %s", v)
		return 1
	}
	return 0
}

func doSessionCreate(args cmdargs) int {
	fs := flag.NewFlagSet("octavctl session create", flag.ContinueOnError)
	var abstract string
	fs.StringVar(&abstract, "abstract", "", "")
	var category string
	fs.StringVar(&category, "category", "", "")
	var conference_id string
	fs.StringVar(&conference_id, "conference_id", "", "")
	var duration int64
	fs.Int64Var(&duration, "duration", 0, "")
	var material_level string
	fs.StringVar(&material_level, "material_level", "", "")
	var memo string
	fs.StringVar(&memo, "memo", "", "")
	var photo_permission string
	fs.StringVar(&photo_permission, "photo_permission", "", "")
	var slide_language string
	fs.StringVar(&slide_language, "slide_language", "", "")
	var slide_subtitles string
	fs.StringVar(&slide_subtitles, "slide_subtitles", "", "")
	var slide_url string
	fs.StringVar(&slide_url, "slide_url", "", "")
	var speaker_id string
	fs.StringVar(&speaker_id, "speaker_id", "", "")
	var spoken_language string
	fs.StringVar(&spoken_language, "spoken_language", "", "")
	var tags string
	fs.StringVar(&tags, "tags", "", "")
	var title string
	fs.StringVar(&title, "title", "", "")
	var video_permission string
	fs.StringVar(&video_permission, "video_permission", "", "")
	var video_url string
	fs.StringVar(&video_url, "video_url", "", "")
	prepGlobalFlags(fs)
	if err := fs.Parse([]string(args)); err != nil {
		return errOut(err)
	}

	m := make(map[string]interface{})
	if abstract != "" {
		m["abstract"] = abstract
	}
	if category != "" {
		m["category"] = category
	}
	if conference_id != "" {
		m["conference_id"] = conference_id
	}
	if duration != 0 {
		m["duration"] = duration
	}
	if material_level != "" {
		m["material_level"] = material_level
	}
	if memo != "" {
		m["memo"] = memo
	}
	if photo_permission != "" {
		m["photo_permission"] = photo_permission
	}
	if slide_language != "" {
		m["slide_language"] = slide_language
	}
	if slide_subtitles != "" {
		m["slide_subtitles"] = slide_subtitles
	}
	if slide_url != "" {
		m["slide_url"] = slide_url
	}
	if speaker_id != "" {
		m["speaker_id"] = speaker_id
	}
	if spoken_language != "" {
		m["spoken_language"] = spoken_language
	}
	if tags != "" {
		m["tags"] = tags
	}
	if title != "" {
		m["title"] = title
	}
	if video_permission != "" {
		m["video_permission"] = video_permission
	}
	if video_url != "" {
		m["video_url"] = video_url
	}
	r := model.CreateSessionRequest{}
	if err := r.Populate(m); err != nil {
		return errOut(err)
	}

	if err := validator.HTTPCreateSessionRequest.Validate(&r); err != nil {
		return errOut(err)
	}

	cl, err := newClient()
	if err != nil {
		return errOut(err)
	}
	res, err := cl.CreateSession(&r)
	if err != nil {
		return errOut(err)
	}
	if err := printJSON(res); err != nil {
		return errOut(err)
	}

	return 0
}

func doSessionLookup(args cmdargs) int {
	fs := flag.NewFlagSet("octavctl session lookup", flag.ContinueOnError)
	var id string
	fs.StringVar(&id, "id", "", "")
	prepGlobalFlags(fs)
	if err := fs.Parse([]string(args)); err != nil {
		return errOut(err)
	}

	m := make(map[string]interface{})
	if id != "" {
		m["id"] = id
	}
	r := model.LookupSessionRequest{}
	if err := r.Populate(m); err != nil {
		return errOut(err)
	}

	if err := validator.HTTPLookupSessionRequest.Validate(&r); err != nil {
		return errOut(err)
	}

	cl, err := newClient()
	if err != nil {
		return errOut(err)
	}
	res, err := cl.LookupSession(&r)
	if err != nil {
		return errOut(err)
	}
	if err := printJSON(res); err != nil {
		return errOut(err)
	}

	return 0
}

func doSessionDelete(args cmdargs) int {
	fs := flag.NewFlagSet("octavctl session delete", flag.ContinueOnError)
	var id string
	fs.StringVar(&id, "id", "", "")
	prepGlobalFlags(fs)
	if err := fs.Parse([]string(args)); err != nil {
		return errOut(err)
	}

	m := make(map[string]interface{})
	if id != "" {
		m["id"] = id
	}
	r := model.DeleteSessionRequest{}
	if err := r.Populate(m); err != nil {
		return errOut(err)
	}

	if err := validator.HTTPDeleteSessionRequest.Validate(&r); err != nil {
		return errOut(err)
	}

	cl, err := newClient()
	if err != nil {
		return errOut(err)
	}
	if err := cl.DeleteSession(&r); err != nil {
		return errOut(err)
	}

	return 0
}

func doSessionSubcmd(args cmdargs) int {
	switch v := args.Get(0); v {
	case "create":
		return doSessionCreate(args.WithFrontPopped())
	case "lookup":
		return doSessionLookup(args.WithFrontPopped())
	case "delete":
		return doSessionDelete(args.WithFrontPopped())
	default:
		log.Printf("unimplemented (conference): %s", v)
		return 1
	}
	return 0
}

func doSubcmd(args cmdargs) int {
	switch v := args.Get(0); v {
	case "conference":
		return doConferenceSubcmd(args.WithFrontPopped())
	case "venue":
		return doVenueSubcmd(args.WithFrontPopped())
	case "room":
		return doRoomSubcmd(args.WithFrontPopped())
	case "session":
		return doSessionSubcmd(args.WithFrontPopped())
	default:
		log.Printf("unimplemented (conference): %s", v)
		return 1
	}
	return 0
}
