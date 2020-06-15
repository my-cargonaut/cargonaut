// Package version provides variables which, when used correctly, provide
// build information about the application referencing this package. At compile
// time the appropriate linker flags must be passed:
//
// 	go build -ldflags "-X axicode.axiom.co/watchmakers/axiomdb/internal/common/version.release=1.0.0"
//
// Adapt the flags for all other exported variables.
package version
