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

func init() {
	SetLocale(Languages[0])
	SetDomain("messages")
}

func Default() *Gettext {
	initOnce.Do(func() {
		defaultGettext = New("locales")
		defaultGettext.AddDomain("messages")
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

func (v *Gettext) Get(locale, domain, s string, args ...interface{}) string {
	v.mu.RLock()
	defer v.mu.RUnlock()

	l, ok := v.Locales[locale]
	if !ok {
		return s
	}
	return l.GetD(domain, s, args...)
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
	l := Default()

	defaultDomainMutex.RLock()
	dd := defaultDomain
	defaultDomainMutex.RUnlock()

	defaultLocaleMutex.RLock()
	dl := defaultLocale
	defaultLocaleMutex.RUnlock()

	return l.Get(dl, dd, s, args...)
}
