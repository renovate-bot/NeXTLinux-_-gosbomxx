package integration

import (
	"testing"

	"github.com/nextlinux/gosbom/gosbom/pkg"
	"github.com/nextlinux/gosbom/gosbom/source"
)

func TestRustAudit(t *testing.T) {
	sbom, _ := catalogFixtureImage(t, "image-rust-auditable", source.SquashedScope, []string{"all"})

	expectedPkgs := 2
	actualPkgs := 0
	for range sbom.Artifacts.Packages.Enumerate(pkg.RustPkg) {
		actualPkgs += 1
	}

	if actualPkgs != expectedPkgs {
		t.Errorf("unexpected number of Rust packages: %d != %d", expectedPkgs, actualPkgs)
	}
}
