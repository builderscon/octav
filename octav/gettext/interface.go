package gettext

import (
	"sync"

	gotext "gopkg.in/leonelquinteros/gotext.v1"
)

type Gettext struct {
	mu      sync.RWMutex
	Locales map[string]*gotext.Locale
}
