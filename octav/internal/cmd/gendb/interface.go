package main

type Scanner interface {
  Scan(...interface{}) error
}

type StructField struct {
  AutoIncrement bool
  Converter     string
  ColumnName    string
  Name          string
  Type          string
  Unique        bool
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