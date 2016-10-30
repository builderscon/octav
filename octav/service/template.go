package service

import (
	"sync"

	"github.com/builderscon/octav/octav/assets"
	"github.com/builderscon/octav/octav/gettext"
	pdebug "github.com/lestrrat/go-pdebug"
	tmplbox "github.com/lestrrat/go-tmplbox"
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

	box := tmplbox.New(tmplbox.AssetSourceFunc(assets.Asset))
	box.Funcs(map[string]interface{}{
		"gettext": gettext.Get,
	})
	v.box = box
}

func (v *TemplateSvc) Get(name string, deps ...string) (tmplbox.Template, error) {
	return v.box.GetOrCompose(name, deps...)
}
