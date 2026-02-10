package appstore

import _ "embed"

//go:embed VERSION
var version string

// Version returns the version of the library.
func Version() string { return version }
