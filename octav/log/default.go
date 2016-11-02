// +build !gcp
// +build !debug

package log

func init() {
	DefaultLogger = nullLog{}
}