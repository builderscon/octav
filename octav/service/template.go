package service

import (
	"io"
	"text/template"

	"github.com/builderscon/octav/octav/assets"
	"github.com/builderscon/octav/octav/gettext"
	pdebug "github.com/lestrrat/go-pdebug"
)

func Template() *TemplateSvc {
	var t *template.Template

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
	}

	return &TemplateSvc{
		template: t,
	}
}

func (v *TemplateSvc) Execute(dst io.Writer, name string, vars interface{}) error {
	return v.template.ExecuteTemplate(dst, name, vars)
}
