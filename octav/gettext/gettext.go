package gettext

import (
	"sync"

	gotext "gopkg.in/leonelquinteros/gotext.v1"
)

var Languages = []string{"en", "ja"}
var defaultGettext *Gettext
var defaultLocale string
var defaultLocaleMutex sync.RWMutex
var defaultDomain string
var defaultDomainMutex sync.RWMutex
var initOnce sync.Once

func Default() *Gettext {
	initOnce.Do(func() {
		defaultGettext = New("locales")
		defaultGettext.AddDomain("messages")

		SetLocale(Languages[0])
		SetDomain("messages")
	})
	return defaultGettext
}

func New(path string) *Gettext {
	var v Gettext

	for _, l := range Languages {
		v.AddLocale(l, gotext.NewLocale(path, l))
	}

	return &v
}

func (v *Gettext) AddDomain(domain string) {
	v.mu.Lock()
	defer v.mu.Unlock()

	for _, l := range v.Locales {
		l.AddDomain(domain)
	}
}

func (v *Gettext) AddLocale(name string, l *gotext.Locale) {
	v.mu.Lock()
	defer v.mu.Unlock()

	if v.Locales == nil {
		v.Locales = map[string]*gotext.Locale{}
	}
	v.Locales[name] = l
}

func (v *Gettext) Get(locale, domain, name string, args ...interface{}) string {
	v.mu.RLock()
	defer v.mu.RUnlock()

	return v.Locales[locale].GetD(domain, name, args...)
}

func SetDomain(name string) {
	defaultDomainMutex.Lock()
	defer defaultDomainMutex.Unlock()

	defaultDomain = name
}

func SetLocale(name string) {
	defaultLocaleMutex.Lock()
	defer defaultLocaleMutex.Unlock()

	defaultLocale = name
}

func Get(s string, args ...interface{}) string {
	defaultDomainMutex.RLock()
	defaultLocaleMutex.RLock()
	defer defaultDomainMutex.RUnlock()
	defer defaultLocaleMutex.RUnlock()

	return defaultGettext.Get(defaultLocale, defaultDomain, s, args...)
}
