package tools

import (
	"sync"
)

type LocalizedFields struct {
	lock sync.RWMutex
	// Language -> field/value
	fields map[string]map[string]string
}

