package service

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/builderscon/octav/octav/gettext"
)

func Template() *TemplateSvc {
	var t *template.Template
	filepath.Walk("template", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if info != nil && info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		if t == nil {
			t = template.New(path).Funcs(map[string]interface{}{
				"gettext": gettext.Get,
			})
		}

		if _, err := t.Parse(string(b)); err != nil {
			return err
		}
		return nil
	})

	return &TemplateSvc{
		template: t,
	}
}

func (v *TemplateSvc) Execute(dst io.Writer, name string, vars interface{}) error {
	return v.template.ExecuteTemplate(dst, name, vars)
}
