package main

type Scanner interface {
	Scan(...interface{}) error
}

type StructField struct {
	JSONName string
	L10N     bool
	Name     string
	Type     string
}

type Struct struct {
	AutoIncrementField *StructField
	CacheEnabled       bool
	CacheExpires       string
	Fields             []StructField
	Name               string
	NoScanner          bool
	PackageName        string
	PreCreate          string
	PreDelete          string
	PrimaryKey         *StructField
	PostCreate         string
	Tablename          string
}

type InspectionCtx struct {
	Marker  string
	Package string
	Structs []Struct
}