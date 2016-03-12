package main

// Given model name X
// 1) Look at its struct representation. This MUST exist
// 2) Look at its database counterpart under db/*. This MUST exist

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/lestrrat/go-pdebug"
)

var ErrAnnotatedStructNotFound = errors.New("annotated struct was not found")

func snakeCase(s string) string {
	ret := []rune{}
	wasLower := false
	for len(s) > 0 {
		r, n := utf8.DecodeRuneInString(s)
		if r == utf8.RuneError {
			panic("yikes")
		}

		s = s[n:]
		if unicode.IsUpper(r) {
			if wasLower {
				ret = append(ret, '_')
			}
			wasLower = false
		} else {
			wasLower = true
		}

		ret = append(ret, unicode.ToLower(r))
	}
	return string(ret)
}

type Processor struct {
	Types []string
	Dir   string
}

func skipGenerated(fi os.FileInfo) bool {
	switch {
	case strings.HasSuffix(fi.Name(), "_gen.go"):
		return false
	case strings.HasSuffix(fi.Name(), "_gen.go"):
		return false
	}
	return true
}

type genctx struct {
	Dir         string
	DBRows      map[string]DBRow
	Models      []Model
	PkgName     string
	TargetTypes []string
}

type Field struct {
	JSONName string
	L10N     bool
	Name     string
	Type     string
}

type Model struct {
	Fields  []Field
	HasEID  bool
	HasL10N bool
	Name    string
	PkgName string
}

type DBColumn struct {
	BaseType string
	IsNullType bool
	Name string
	Type string
}

type DBRow struct {
	Columns map[string]DBColumn
	Name    string
	PkgName string
}

func (p *Processor) Do() error {
	ctx := genctx{
		Dir: p.Dir,
		DBRows: make(map[string]DBRow),
		TargetTypes: p.Types,
	}
	if err := parseModelDir(&ctx, ctx.Dir); err != nil {
		return err
	}

	if err := parseDBDir(&ctx, filepath.Join(ctx.Dir, "db")); err != nil {
		return err
	}

	if err := generateFiles(&ctx); err != nil {
		return err
	}

	return nil
}

func parseDir(ctx *genctx, dir string, cb func(*genctx, *ast.Package) error) error {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, skipGenerated, parser.ParseComments)
	if err != nil {
		return err
	}

	if len(pkgs) == 0 {
		return errors.New("no packages to process...")
	}

	for _, pkg := range pkgs {
		if strings.HasSuffix(pkg.Name, "_test") {
			continue
		}

		ctx.PkgName = pkg.Name
		if err := cb(ctx, pkg); err != nil {
			return err
		}
		return nil
	}

	return errors.New("only found test package")
}

func parseModelDir(ctx *genctx, dir string) error {
	return parseDir(ctx, dir, processModelPkg)
}

func parseDBDir(ctx *genctx, dir string) error {
	return parseDir(ctx, dir, processDBPkg)
}

func processPkg(ctx *genctx, pkg *ast.Package, cb func(ast.Node) bool) error {
	if pdebug.Enabled {
		g := pdebug.Marker("processPkg %s", pkg.Name)
		defer g.End()
	}

	for fn, f := range pkg.Files {
		if pdebug.Enabled {
			pdebug.Printf("Checking file %s", fn)
		}

		ast.Inspect(f, cb)
	}
	return nil
}

func processModelPkg(ctx *genctx, pkg *ast.Package) error {
	if pdebug.Enabled {
		g := pdebug.Marker("processModelPkg %s", pkg.Name)
		defer g.End()
	}

	if err := processPkg(ctx, pkg, ctx.extractModelStructs); err != nil {
		return err
	}
	return nil
}

func processDBPkg(ctx *genctx, pkg *ast.Package) error {
	if pdebug.Enabled {
		g := pdebug.Marker("processDBPkg %s", pkg.Name)
		defer g.End()
	}

	if err := processPkg(ctx, pkg, ctx.extractDBStructs); err != nil {
		return err
	}
	return nil
}

func shouldProceed(ctx *genctx, name string) bool {
	if len(ctx.TargetTypes) == 0 {
		return true
	}

	for _, t := range ctx.TargetTypes {
		if t == name {
			return true
		}
	}
	return false
}

func (ctx *genctx) extractModelStructs(n ast.Node) bool {
	decl, ok := n.(*ast.GenDecl)
	if !ok {
		return true
	}

	if decl.Tok != token.TYPE {
		return true
	}

	for _, spec := range decl.Specs {
		var t *ast.TypeSpec
		var s *ast.StructType
		var ok bool

		if t, ok = spec.(*ast.TypeSpec); !ok {
			continue
		}

		if !shouldProceed(ctx, t.Name.Name) {
			continue
		}

		if s, ok = t.Type.(*ast.StructType); !ok {
			continue
		}

		st := Model{
			Fields:  make([]Field, 0, len(s.Fields.List)),
			Name:    t.Name.Name,
			PkgName: ctx.PkgName,
		}

	LoopFields:
		for _, f := range s.Fields.List {
			if len(f.Names) == 0 {
				continue
			}

			if unicode.IsLower(rune(f.Names[0].Name[0])) {
				continue
			}

			var jsname string
			var l10n bool
			if f.Tag != nil {
				v := f.Tag.Value
				if len(v) >= 2 {
					if v[0] == '`' {
						v = v[1:]
					}
					if v[len(v)-1] == '`' {
						v = v[:len(v)-1]
					}
				}

				st := reflect.StructTag(v)
				tag := st.Get("json")
				if tag == "-" {
					continue LoopFields
				}
				if tag == "" || tag[0] == ',' {
					jsname = f.Names[0].Name
				} else {
					tl := strings.SplitN(tag, ",", 2)
					jsname = tl[0]
				}

				tag = st.Get("l10n")
				if b, err := strconv.ParseBool(tag); err == nil && b {
					l10n = true
				}
			}

			typ, err := getTypeName(f.Type)
			if err != nil {
				return true
			}

			field := Field{
				L10N:     l10n,
				Name:     f.Names[0].Name,
				JSONName: jsname,
				Type:     typ,
			}

			st.Fields = append(st.Fields, field)
		}
		ctx.Models = append(ctx.Models, st)
	}

	return true
}

func getTypeName(ref ast.Expr) (string, error) {
	var typ string
	switch ref.(type) {
	case *ast.Ident:
		typ = ref.(*ast.Ident).Name
	case *ast.SelectorExpr:
		typ = ref.(*ast.SelectorExpr).Sel.Name
	case *ast.StarExpr:
		return getTypeName(ref.(*ast.StarExpr).X)
	case *ast.ArrayType:
		typ = "[]" + ref.(*ast.ArrayType).Elt.(*ast.Ident).Name
	case *ast.MapType:
		mt := ref.(*ast.MapType)
		typ = "map[" + mt.Key.(*ast.Ident).Name + "]" + mt.Value.(*ast.Ident).Name
	default:
		fmt.Printf("%#v\n", ref)
		return "", errors.New("field type not supported")
	}
	return typ, nil
}

func (ctx *genctx) extractDBStructs(n ast.Node) bool {
	decl, ok := n.(*ast.GenDecl)
	if !ok {
		return true
	}

	if decl.Tok != token.TYPE {
		return true
	}

	for _, spec := range decl.Specs {
		var t *ast.TypeSpec
		var s *ast.StructType
		var ok bool

		if t, ok = spec.(*ast.TypeSpec); !ok {
			continue
		}

		if !shouldProceed(ctx, t.Name.Name) {
			continue
		}

		if s, ok = t.Type.(*ast.StructType); !ok {
			continue
		}

		st := DBRow{
			Columns:  make(map[string]DBColumn),
			Name:    t.Name.Name,
			PkgName: ctx.PkgName,
		}

		for _, f := range s.Fields.List {
			if len(f.Names) == 0 {
				continue
			}

			if unicode.IsLower(rune(f.Names[0].Name[0])) {
				continue
			}

			typ, err := getTypeName(f.Type)
			if err != nil {
				return true
			}

			// If this is a Null* field, record that
			var nulltyp bool
			var basetyp string

			{
				// extract the package portion
				var prefix string
				if dotpos := strings.IndexRune(typ, '.'); dotpos > -1 {
					prefix = typ[:dotpos]
				}

				if i := strings.Index(typ, prefix + "Null"); i > -1 {
					nulltyp = true
					basetyp = typ[len(prefix)+i+4:]
				}
			}

			column := DBColumn{
				BaseType: basetyp,
				IsNullType: nulltyp,
				Name:     f.Names[0].Name,
				Type:     typ,
			}

			st.Columns[column.Name] = column
		}
		ctx.DBRows[st.Name] = st
	}

	return true
}

func generateFiles(ctx *genctx) error {
	for _, m := range ctx.Models {
		if pdebug.Enabled {
			pdebug.Printf("Checking model %s", m.Name)
		}
		if !shouldProceed(ctx, m.Name) {
			if pdebug.Enabled {
				pdebug.Printf("Skipping model %s", m.Name)
			}
			continue
		}

		if err := generateModelFile(ctx, m); err != nil {
			return err
		}
	}
	return nil
}

func generateModelFile(ctx *genctx, m Model) error {
	if pdebug.Enabled {
		g := pdebug.Marker("generateModelFile %s", m.Name)
		defer g.End()
	}

	// Find the matching DBRow
	row, ok  := ctx.DBRows[m.Name]
	if !ok {
		return errors.New("could not find matching row for " + m.Name)
	}

	buf := bytes.Buffer{}

	varname := 'v'
	hasID := false

	buf.WriteString("// Automatically generated by genmodel utility. DO NOT EDIT!\n")
	buf.WriteString("package ")
	buf.WriteString(m.PkgName)
	buf.WriteString("\n\n")
	buf.WriteString("\nimport (\n")
	buf.WriteString("\n" + strconv.Quote("encoding/json"))
	buf.WriteString("\n" + strconv.Quote("github.com/builderscon/octav/octav/db"))
	buf.WriteString("\n" + strconv.Quote("github.com/lestrrat/go-pdebug"))
	buf.WriteString("\n)")

	fmt.Fprintf(&buf, "\n\nfunc (%c %s) GetPropNames() ([]string, error) {", varname, m.Name)
	fmt.Fprintf(&buf, "\nl, _ := %c.L10N.GetPropNames()", varname)
	buf.WriteString("\nreturn append(l, ")

	l10nfields := bytes.Buffer{}
	for _, f := range m.Fields {
		buf.WriteString(strconv.Quote(f.JSONName))
		buf.WriteString(",")
		if f.Name == "ID" {
			hasID = true
		}
		if f.L10N {
			l10nfields.WriteString(strconv.Quote(f.JSONName))
			l10nfields.WriteString(",")
		}
	}
	buf.WriteString("), nil")
	buf.WriteString("\n}")

	fmt.Fprintf(&buf, "\n\nfunc (%c %s) GetPropValue(s string) (interface{}, error) {", varname, m.Name)
	buf.WriteString("\nswitch s {")
	for _, f := range m.Fields {
		fmt.Fprintf(&buf, "\ncase %s:", strconv.Quote(f.JSONName))
		fmt.Fprintf(&buf, "\nreturn %c.%s, nil", varname, f.Name)
	}
	buf.WriteString("\ndefault:")
	fmt.Fprintf(&buf, "\nreturn %c.L10N.GetPropValue(s)", varname)
	buf.WriteString("\n}\n}")

	fmt.Fprintf(&buf, "\n\nfunc (%c %s) MarshalJSON() ([]byte, error) {", varname, m.Name)
	buf.WriteString("\nm := make(map[string]interface{})")
	for _, f := range m.Fields {
		fmt.Fprintf(&buf, "\nm[%s] = %c.%s", strconv.Quote(f.JSONName), varname, f.Name)
	}
	buf.WriteString("\nbuf, err := json.Marshal(m)")
	buf.WriteString("\nif err != nil {")
	buf.WriteString("\nreturn nil, err")
	buf.WriteString("\n}")
	fmt.Fprintf(&buf, "\nreturn marshalJSONWithL10N(buf, %c.L10N)", varname)
	buf.WriteString("\n}")

	fmt.Fprintf(&buf, "\n\nfunc (%c *%s) UnmarshalJSON(data []byte) error {", varname, m.Name)
	buf.WriteString("\nm := make(map[string]interface{})")
	buf.WriteString("\nif err := json.Unmarshal(data, &m); err != nil {")
	buf.WriteString("\nreturn err")
	buf.WriteString("\n}")

	for _, f := range m.Fields {
		fmt.Fprintf(&buf, "\n\nif jv, ok := m[%s]; ok {", strconv.Quote(f.JSONName))
		buf.WriteString("\nswitch jv.(type) {")
		if strings.Contains(f.Type, "int") {
			buf.WriteString("\ncase float64:")
			fmt.Fprintf(&buf, "\n%c.%s = %s(jv.(float64))", varname, f.Name, f.Type)
		} else {
			fmt.Fprintf(&buf, "\ncase %s:", f.Type)
			fmt.Fprintf(&buf, "\n%c.%s = jv.(%s)", varname, f.Name, f.Type)
		}
		fmt.Fprintf(&buf, "\ndelete(m, %s)", strconv.Quote(f.JSONName))
		buf.WriteString("\ndefault:")
		fmt.Fprintf(&buf, "\nreturn ErrInvalidFieldType{Field: %s}", strconv.Quote(f.JSONName))
		buf.WriteString("\n}")
		buf.WriteString("\n}")
	}

	if l10nfields.Len() > 0 {
		fmt.Fprintf(&buf, "\n\nif err := ExtractL10NFields(m, &v.L10N, []string{%s}); err != nil {", l10nfields.String())
		buf.WriteString("\nreturn err")
		buf.WriteString("\n}")
	}

	buf.WriteString("\nreturn nil")
	buf.WriteString("\n}")

	if hasID {
		fmt.Fprintf(&buf, "\n\nfunc (v *%s) Load(tx *db.Tx, id string) error {", m.Name)
		fmt.Fprintf(&buf, "\nvdb := db.%s{}", m.Name)
		buf.WriteString("\nif err := vdb.LoadByEID(tx, id); err != nil {")
		buf.WriteString("\nreturn err")
		buf.WriteString("\n}\n")
		buf.WriteString("\nif err := v.FromRow(vdb); err != nil {")
		buf.WriteString("\nreturn err")
		buf.WriteString("\n}")
		buf.WriteString("\nif err := v.LoadLocalizedFields(tx); err != nil {")
		buf.WriteString("\nreturn err")
		buf.WriteString("\n}")
		buf.WriteString("\nreturn nil")
		buf.WriteString("\n}")

		fmt.Fprintf(&buf, "\n\nfunc (v *%s) LoadLocalizedFields(tx *db.Tx) error {", m.Name)
		fmt.Fprintf(&buf, "\nls, err := db.LoadLocalizedStringsForParent(tx, v.ID, %s)", strconv.Quote(m.Name))
		buf.WriteString("\nif err != nil {")
		buf.WriteString("\nreturn err")
		buf.WriteString("\n}")
		buf.WriteString("\n\nif len(ls) > 0 {")
		buf.WriteString("\nv.L10N = LocalizedFields{}")
		buf.WriteString("\nfor _, l := range ls {")
		buf.WriteString("\nv.L10N.Set(l.Language, l.Name, l.Localized)")
		buf.WriteString("\n}")
		buf.WriteString("\n}")
		buf.WriteString("\nreturn nil")
		buf.WriteString("\n}")

		fmt.Fprintf(&buf, "\n\nfunc (v *%s) FromRow(vdb db.%s) error {", m.Name, m.Name)
		buf.WriteString("\nv.ID = vdb.EID")
		for _, f := range m.Fields {
			if f.Name == "ID" {
				continue
			}

			c, ok := row.Columns[f.Name]
			if !ok {
				continue
			}

			if c.IsNullType {
				fmt.Fprintf(&buf, "\nif vdb.%s.Valid {", f.Name)
				fmt.Fprintf(&buf, "\nv.%s = vdb.%s.%s", f.Name, f.Name, c.BaseType)
				buf.WriteString("\n}")
			} else {
				fmt.Fprintf(&buf, "\nv.%s = vdb.%s", f.Name, f.Name)
			}
		}
		buf.WriteString("\nreturn nil")
		buf.WriteString("\n}")
	}

	fmt.Fprintf(&buf, "\n\nfunc (v *%s) ToRow(vdb *db.%s) error {", m.Name, m.Name)
	for _, f := range m.Fields {
		if f.Name == "ID" {
			buf.WriteString("\nvdb.EID = v.ID")
		} else {
			c, ok := row.Columns[f.Name]
			if !ok {
				continue
			}

			if c.IsNullType {
				fmt.Fprintf(&buf, "\nvdb.%s.Valid = true", f.Name)
				fmt.Fprintf(&buf, "\nvdb.%s.%s = v.%s", f.Name, c.BaseType, f.Name)
			} else {
				fmt.Fprintf(&buf, "\nvdb.%s = v.%s", f.Name, f.Name)
			}
		}
	}
	buf.WriteString("\nreturn nil")
	buf.WriteString("\n}")

	fmt.Fprintf(&buf, "\n\nfunc (v *%s) Create(tx *db.Tx) error {", m.Name)
	if hasID {
		buf.WriteString("\n" + `if v.ID == "" {`)
		buf.WriteString("\nv.ID = UUID()")
		buf.WriteString("\n}")
	}
	fmt.Fprintf(&buf, "\n\nvdb := db.%s{}", m.Name)
	buf.WriteString("\nv.ToRow(&vdb)")
	buf.WriteString("\nif err := vdb.Create(tx); err != nil {")
	buf.WriteString("\nreturn err")
	buf.WriteString("\n}\n")
	if hasID {
		fmt.Fprintf(&buf, "\nif err := v.L10N.CreateLocalizedStrings(tx, %s, v.ID); err != nil {", strconv.Quote(m.Name))
		buf.WriteString("\nreturn err")
		buf.WriteString("\n}")
	}
	buf.WriteString("\nreturn nil")
	buf.WriteString("\n}")

	fmt.Fprintf(&buf, "\n\nfunc (v *%s) Update(tx *db.Tx) (err error) {", m.Name)
	buf.WriteString("\nif pdebug.Enabled {")
	fmt.Fprintf(&buf, "\n" + `g := pdebug.Marker("%s.Update (%%s)", v.ID).BindError(&err)`, m.Name)
	buf.WriteString("\ndefer g.End()")
	buf.WriteString("\n}")
	fmt.Fprintf(&buf, "\n\nvdb := db.%s{}", m.Name)
	buf.WriteString("\nv.ToRow(&vdb)")
	buf.WriteString("\nif err := vdb.Update(tx); err != nil {")
	buf.WriteString("\nreturn err")
	buf.WriteString("\n}")
	buf.WriteString("\n\nreturn v.L10N.Foreach(func(l, k, x string) error {")
	buf.WriteString("\nls := db.LocalizedString{")
	fmt.Fprintf(&buf, "\nParentType: %s,", strconv.Quote(m.Name))
	buf.WriteString("\nParentID: v.ID,")
	buf.WriteString("\nLanguage: l,")
	buf.WriteString("\nName: k,")
	buf.WriteString("\nLocalized: x,")
	buf.WriteString("\n}")
	buf.WriteString("\nreturn ls.Upsert(tx)")
	buf.WriteString("\n})")
	buf.WriteString("\n}")

	fmt.Fprintf(&buf, "\n\nfunc (v *%s) Delete(tx *db.Tx) error {", m.Name)
	buf.WriteString("\nif pdebug.Enabled {")
	fmt.Fprintf(&buf, "\n"+`g := pdebug.Marker("%s.Delete (%%s)", v.ID)`, m.Name)
	buf.WriteString("\ndefer g.End()")
	buf.WriteString("\n}")
	fmt.Fprintf(&buf, "\n\nvdb := db.%s{EID: v.ID}", m.Name)
	buf.WriteString("\nif err := vdb.Delete(tx); err != nil {")
	buf.WriteString("\nreturn err")
	buf.WriteString("\n}")
	fmt.Fprintf(&buf, "\nif err := db.DeleteLocalizedStringsForParent(tx, v.ID, %s); err != nil {", strconv.Quote(m.Name))
	buf.WriteString("\nreturn err")
	buf.WriteString("\n}")
	buf.WriteString("\nreturn nil")
	buf.WriteString("\n}")

	if hasID {
		fmt.Fprintf(&buf, "\n\nfunc (v *%sList) Load(tx *db.Tx, since string, limit int) error {", m.Name)
		fmt.Fprintf(&buf, "\nvdbl := db.%sList{}", m.Name)
		buf.WriteString("\nif err := vdbl.LoadSinceEID(tx, since, limit); err != nil {")
		buf.WriteString("\nreturn err")
		buf.WriteString("\n}")
		fmt.Fprintf(&buf, "\nres := make([]%s, len(vdbl))", m.Name)
		buf.WriteString("\nfor i, vdb := range vdbl {")
		fmt.Fprintf(&buf, "\nv := %s{}", m.Name)
		buf.WriteString("\nif err := v.FromRow(vdb); err != nil {")
		buf.WriteString("\nreturn err")
		buf.WriteString("\n}")
		buf.WriteString("\nif err := v.LoadLocalizedFields(tx); err != nil {")
		buf.WriteString("\nreturn err")
		buf.WriteString("\n}")
		buf.WriteString("\nres[i] = v")
		buf.WriteString("\n}")
		buf.WriteString("\n*v = res")
		buf.WriteString("\nreturn nil")
		buf.WriteString("\n}")
	}

	fsrc, err := format.Source(buf.Bytes())
	if err != nil {
		log.Printf("%s", buf.Bytes())
		return err
	}

	fn := filepath.Join(ctx.Dir, snakeCase(m.Name)+"_gen.go")
	if pdebug.Enabled {
		pdebug.Printf("Generating file %s", fn)
	}
	fi, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer fi.Close()

	if _, err := fi.Write(fsrc); err != nil {
		return err
	}

	return nil
}
