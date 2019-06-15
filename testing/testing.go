// Package testing ...
// Pretty neat little trick from:
// https://brandur.org/fragments/testing-go-project-root
package testing

import (
	"os"
	"path"
	"runtime"
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	_ = os.Chdir(dir)
}
