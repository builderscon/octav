package cache

import "time"

type Option interface {
	Name() string
	Value() interface{}
}

type optionWithValue struct {
	name  string
	value interface{}
}

func (o *optionWithValue) Name() string {
	return o.name
}

func (o *optionWithValue) Value() interface{} {
	return o.value
}

func WithExpires(t time.Duration) Option {
	return &optionWithValue{name:"expires", value:t}
}
