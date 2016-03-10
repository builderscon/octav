package main

type Scanner interface {
	Scan(...interface{}) error
}

type StructField struct {
	JSONName string
	URLName string
	L10N     bool
	Maybe    bool
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
	HasL10N bool
}

type InspectionCtx struct {
	Marker  string
	Package string
	Structs []Struct
}