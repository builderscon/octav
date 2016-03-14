package main

type StructField struct {
	JSONName string
	URLName  string
	L10N     bool
	Maybe    bool
	Name     string
	Type     string
}

type Struct struct {
	Fields      []StructField
	Name        string
	PackageName string
	HasL10N     bool
}

type InspectionCtx struct {
	Marker  string
	Package string
	Structs []Struct
}