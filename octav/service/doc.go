// Package service contain objects that know how to interconnect
// the data from outside world (model) to internal data (db, et al).
//
// For example, this is the only component that understands the
// mapping between optional incoming parameters from forms
// and the database row.
package service