package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"go/format"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/lestrrat/go-jshschema"
	"github.com/lestrrat/go-jsschema"
)

type Cmd struct {
	Name    string   `json:"name"`    // "conference"
	Actions []Action `json:"actions"` // []string{ "create", "lookup", "delete" }
	Subcmds []Cmd    `json:"subcmds"` // Cmd{ "dates", "admin" }
}

type Action struct {
	Name     string                       `json:"name"`
	ArgsHint map[string]map[string]string `json:"args_hint"`
}

func (a *Action) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch v.(type) {
	case string:
		*a = Action{Name: v.(string)}
	default:
		m := v.(map[string]interface{})
		ah := make(map[string]map[string]string)
		for k, v := range m["args_hint"].(map[string]interface{}) {
			m2 := make(map[string]string)
			ah[k] = m2
			for k2, v2 := range v.(map[string]interface{}) {
				m2[k2] = v2.(string)
			}
		}
		*a = Action{Name: m["name"].(string), ArgsHint: ah}
	}
	return nil
}

type genctx struct {
	cmdchain []string
	schema   *hschema.HyperSchema
	out      io.Writer
}

func main() {
	os.Exit(_main())
}

func _main() int {
	var specfile string
	var tmplfile string
	var outputfile string
	flag.StringVar(&specfile, "s", "", "JSON Hyper Schema spec file to use")
	flag.StringVar(&tmplfile, "t", "", "Template JSON file that contains the command structure")
	flag.StringVar(&outputfile, "o", "", "Output file")
	flag.Parse()

	if specfile == "" || tmplfile == "" || outputfile == "" {
		log.Printf("Usage: genoctavctl -s /path/to/hyperschema.json -t /path/to/octavctl.json -o /path/to/output.go")
		return 1
	}

	var cmd Cmd
	if err := readCmdSpec(tmplfile, &cmd); err != nil {
		log.Printf("%s", err)
		return 1
	}

	s, err := hschema.ReadFile(specfile)
	if err != nil {
		log.Printf("%s", err)
		return 1
	}

	buf := bytes.Buffer{}
	buf.WriteString(preamble)
	ctx := genctx{
		out:    &buf,
		schema: s,
	}
	if err := processCmd(&ctx, cmd); err != nil {
		log.Printf("%s", err)
		return 1
	}

	fsrc, err := format.Source(buf.Bytes())
	if err != nil {
		log.Printf("%s", buf.Bytes())
		log.Printf("%s", err)
		return 1
	}

	out, err := os.Create(outputfile)
	if err != nil {
		log.Printf("failed to create file %s: %s", outputfile, err)
		return 1
	}
	defer out.Close()

	out.Write(fsrc)

	return 0
}

const preamble = `
package main

// AUTO GENERATED FILE! DO NOT EDIT

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/builderscon/octav/octav/client"
	"github.com/builderscon/octav/octav/internal/homedir"
	"github.com/builderscon/octav/octav/model"
	"github.com/builderscon/octav/octav/validator"
)

// Global options
var endpoint string

func prepGlobalFlags(fs *flag.FlagSet) {
	fs.StringVar(&endpoint, "endpoint", endpoint, "Base URL of the octav server (required)")
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
		if i != len(*l) - 1 {
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

func readConfigFile() error {
	// This is a placeholder. We probably should implement a real
	// Config object later
	dir, err := homedir.Get()
	if err != nil {
		return err
	}

	f := filepath.Join(dir, ".octavrc")
	if _, err := os.Stat(f); err != nil {
		return nil // doesn't exist. that's fine
	}

	cf, err := os.Open(f)
	if err != nil {
		return err
	}
	defer cf.Close()

	config := map[string]string{}
	if err := json.NewDecoder(cf).Decode(&config); err != nil {
		return err
	}
	if v, ok := config["endpoint"]; ok {
		endpoint = v
	}
	return nil
}

func main() {
	if err := readConfigFile(); err != nil {
		os.Exit(errOut(err))
	}

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
`

func readCmdSpec(fname string, cmd *Cmd) error {
	fh, err := os.Open(fname)
	if err != nil {
		log.Printf("failed to read %s: %s", fname, err)
		return err
	}
	defer fh.Close()

	if err := json.NewDecoder(fh).Decode(&cmd); err != nil {
		return err
	}

	return nil
}

func processCmd(ctx *genctx, cmd Cmd) error {
	ctx.cmdchain = append(ctx.cmdchain, cmd.Name)
	defer func() {
		ctx.cmdchain = ctx.cmdchain[:len(ctx.cmdchain)-1]
	}()

	methodName := toMethodName(ctx.cmdchain[1:])
	buf := bytes.Buffer{}
	fmt.Fprintf(&buf, "\n\nfunc do%sSubcmd(args cmdargs) int{", methodName)
	buf.WriteString("\nswitch v := args.Get(0); v {")
	for _, action := range cmd.Actions {
		nextcmdchain := append(ctx.cmdchain[1:], action.Name)
		fmt.Fprintf(&buf, "\ncase %s:", strconv.Quote(action.Name))
		fmt.Fprintf(&buf, "\nreturn do%s(args.WithFrontPopped())", toMethodName(nextcmdchain))
		if err := processAction(ctx, action); err != nil {
			return err
		}
	}

	for _, subcmd := range cmd.Subcmds {
		nextcmdchain := append(ctx.cmdchain[1:], subcmd.Name)
		fmt.Fprintf(&buf, "\ncase %s:", strconv.Quote(subcmd.Name))
		fmt.Fprintf(&buf, "\nreturn do%sSubcmd(args.WithFrontPopped())", toMethodName(nextcmdchain))
		if err := processCmd(ctx, subcmd); err != nil {
			return err
		}
	}
	buf.WriteString("\ndefault:")
	buf.WriteString("\n" + `log.Printf("unimplemented (conference): %s", v)`)
	buf.WriteString("\nreturn 1")
	buf.WriteString("\n}")
	buf.WriteString("\nreturn 0")
	buf.WriteString("\n}")

	buf.WriteTo(ctx.out)
	return nil
}

func toMethodName(s []string) string {
	if len(s) == 0 {
		return ""
	}

	buf := bytes.Buffer{}
	for _, v := range s {
		buf.WriteRune(unicode.ToUpper(rune(v[0])))
		buf.WriteString(v[1:])
	}
	return buf.String()
}

var wsrx = regexp.MustCompile(`\s+`)

func titleToName(s string) string {
	buf := bytes.Buffer{}
	for _, p := range wsrx.Split(s, -1) {
		if len(p) > 0 {
			buf.WriteString(strings.ToUpper(p[:1]))
			buf.WriteString(p[1:])
		}
	}
	return buf.String()
}

func methodNameToRequestTransportName(s *hschema.HyperSchema, name string) string {
	transportNs, ok := s.Extras["hsup.transport_ns"]
	if !ok {
		transportNs = "model"
	}

	return fmt.Sprintf("%s.%sRequest", transportNs, name)
}

func processAction(ctx *genctx, action Action) error {
	ctx.cmdchain = append(ctx.cmdchain, action.Name)
	defer func() {
		ctx.cmdchain = ctx.cmdchain[:len(ctx.cmdchain)-1]
	}()

	cmdall := strings.Join(ctx.cmdchain, " ")
	maincmd := ctx.cmdchain[1]
	endpoint := "/v1/" + strings.Join(ctx.cmdchain[1:], "/")
	link, err := lookupLink(ctx, endpoint)
	if err != nil {
		return err
	}

	exclprefix := maincmd + "_"
	propnames := make([]string, 0, len(link.Schema.Properties))
	for key := range link.Schema.Properties {
		propnames = append(propnames, key)
	}
	sort.Strings(propnames)

	methodName := toMethodName(ctx.cmdchain[1:])

	buf := bytes.Buffer{}
	fmt.Fprintf(&buf, "\n\nfunc do%s(args cmdargs) int {", methodName)
	fmt.Fprintf(&buf, "\nfs := flag.NewFlagSet(%s, flag.ContinueOnError)", strconv.Quote(cmdall))

	setbuf := bytes.Buffer{}

	_, hasID := link.Schema.Properties["id"]

	for _, pname := range propnames {
		pdef, ok := link.Schema.Properties[pname]
		if !ok {
			panic("Could not find property definition for '" + pname + "'")
		}
		if !pdef.IsResolved() {
			rs, err := pdef.Resolve(ctx.schema)
			if err != nil {
				return err
			}
			pdef = rs
		}

		var sansprefix string
		if hasID {
			sansprefix = pname
		} else {
			sansprefix = strings.TrimPrefix(pname, exclprefix)
		}

		argt := "string"
		argm := "StringVar"
		argv := `""`
		argz := `""`
		if hint, ok := action.ArgsHint[sansprefix]; ok {
			if t, ok := hint["type"]; ok {
				argt = t // Must implement reflect.Setter and reflect.Getter
				argm = "Var"
			}
		} else if len(pdef.Type) == 1 { // Can't handle multiple types, too darn hard
			switch pdef.Type[0] {
			case schema.IntegerType:
				argt = "int64"
				argm = "Int64Var"
				argv = "0"
				argz = "0"
			case schema.NumberType:
				argt = "float64"
				argm = "Float64Var"
				argv = "0"
				argz = "0"
			}
		}
		fmt.Fprintf(&buf, "\nvar %s %s", sansprefix, argt)
		if argm == "Var" {
			fmt.Fprintf(&buf, "\n"+`fs.Var(&%s, %s, "")`, sansprefix, strconv.Quote(sansprefix))
		} else {
			fmt.Fprintf(&buf, "\n"+`fs.%s(&%s, %s, %s, "")`, argm, sansprefix, strconv.Quote(sansprefix), argv)
		}

		_, hasHint := action.ArgsHint[pname]
		if hasHint {
			fmt.Fprintf(&setbuf, "\nif %s.Valid() {", sansprefix)
		} else {
			fmt.Fprintf(&setbuf, "\nif %s != %s {", sansprefix, argz)
		}
		fmt.Fprintf(&setbuf, "\nm[%s] = %s", strconv.Quote(pname), sansprefix)
		if hasHint {
			// Must implement reflect.Setter and reflect.Getter
			setbuf.WriteString(".Get()")
		}
		setbuf.WriteString("\n}")
	}
	buf.WriteString("\nprepGlobalFlags(fs)")
	buf.WriteString("\nif err := fs.Parse([]string(args)); err != nil {")
	buf.WriteString("\nreturn errOut(err)")
	buf.WriteString("\n}")
	buf.WriteString("\n\nm := make(map[string]interface{})")
	setbuf.WriteTo(&buf)

	clMethodName := titleToName(link.Title)
	transport := methodNameToRequestTransportName(ctx.schema, clMethodName)

	fmt.Fprintf(&buf, "\nr := %s{}", transport)
	buf.WriteString("\nif err := r.Populate(m); err != nil {")
	buf.WriteString("\nreturn errOut(err)")
	buf.WriteString("\n}")

	validatorName := guessValidatorName(transport)
	fmt.Fprintf(&buf, "\n\nif err := %s.Validate(&r); err != nil {", validatorName)
	buf.WriteString("\nreturn errOut(err)")
	buf.WriteString("\n}")

	buf.WriteString("\n\ncl, err := newClient()")
	buf.WriteString("\nif err != nil {")
	buf.WriteString("\nreturn errOut(err)")
	buf.WriteString("\n}")
	targetSchema := link.TargetSchema

	var args bytes.Buffer
	args.WriteString("&r")
	if link.EncType == "multipart/form-data" {
		args.WriteString(", nil")
	}

	if targetSchema == nil {
		fmt.Fprintf(&buf, "\nif err := cl.%s(%s); err != nil {", clMethodName, args.String())
		buf.WriteString("\nreturn errOut(err)")
		buf.WriteString("\n}")
	} else {
		fmt.Fprintf(&buf, "\nres, err := cl.%s(%s)", clMethodName, args.String())
		buf.WriteString("\nif err != nil {")
		buf.WriteString("\nreturn errOut(err)")
		buf.WriteString("\n}")

		buf.WriteString("\nif err := printJSON(res); err != nil {")
		buf.WriteString("\nreturn errOut(err)")
		buf.WriteString("\n}")
	}
	buf.WriteString("\n\nreturn 0")

	buf.WriteString("\n}")
	buf.WriteTo(ctx.out)

	return nil
}

func guessValidatorName(s string) string {
	s = strings.TrimPrefix(s, "model.")
	return "validator.HTTP" + s
}

func lookupLink(ctx *genctx, endpoint string) (*hschema.Link, error) {
	for _, link := range ctx.schema.Links {
		href := link.Path()
		if href == endpoint {
			return link, nil
		}
	}
	return nil, errors.New("link '" + endpoint + "' not found")
}