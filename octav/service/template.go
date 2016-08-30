package service

import (
	"io"
	"sync"
	"text/template"

	"github.com/builderscon/octav/octav/assets"
	"github.com/builderscon/octav/octav/gettext"
	pdebug "github.com/lestrrat/go-pdebug"
)

var templateSvc TemplateSvc
var templateOnce sync.Once

func Template() *TemplateSvc {
	templateOnce.Do(templateSvc.Init)
	return &templateSvc
}

func (v *TemplateSvc) Init() {
	if pdebug.Enabled {
		g := pdebug.Marker("service.Template.Init")
		defer g.End()
	}

	var t *template.Template

	var parsed int
	for _, n := range assets.AssetNames() {
		b, err := assets.Asset(n)
		if err != nil {
			panic(err.Error())
		}

		if t == nil {
			t = template.New(n).Funcs(map[string]interface{}{
				"gettext": gettext.Get,
			})
		}

		var tmpl *template.Template
		if n == t.Name() {
			tmpl = t
		} else {
			tmpl = t.New(n)
		}

		if pdebug.Enabled {
			pdebug.Printf("Parsing template %s", n)
		}
		if _, err := tmpl.Parse(string(b)); err != nil {
			panic(err.Error())
		}
		parsed++
	}

	if pdebug.Enabled {
		pdebug.Printf("Parsed %d templates", parsed)
	}

	v.template = t
}

func (v *TemplateSvc) Execute(dst io.Writer, name string, vars interface{}) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.Template.Execute").BindError(&err)
		defer g.End()
	}

	return v.template.ExecuteTemplate(dst, name, vars)
}
