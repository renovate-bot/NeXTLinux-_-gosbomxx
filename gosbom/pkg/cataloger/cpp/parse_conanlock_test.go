package cpp

import (
	"testing"

	"github.com/nextlinux/gosbom/gosbom/artifact"
	"github.com/nextlinux/gosbom/gosbom/pkg"
	"github.com/nextlinux/gosbom/gosbom/pkg/cataloger/internal/pkgtest"
	"github.com/nextlinux/gosbom/gosbom/source"
)

func TestParseConanlock(t *testing.T) {
	fixture := "test-fixtures/conan.lock"
	expected := []pkg.Package{
		{
			Name:         "zlib",
			Version:      "1.2.12",
			PURL:         "pkg:conan/zlib@1.2.12",
			Locations:    source.NewLocationSet(source.NewLocation(fixture)),
			Language:     pkg.CPP,
			Type:         pkg.ConanPkg,
			MetadataType: pkg.ConanLockMetadataType,
			Metadata: pkg.ConanLockMetadata{
				Ref: "zlib/1.2.12",
				Options: map[string]string{
					"fPIC":   "True",
					"shared": "False",
				},
				Path:    "all/conanfile.py",
				Context: "host",
			},
		},
	}

	// TODO: relationships are not under test
	var expectedRelationships []artifact.Relationship

	pkgtest.TestFileParser(t, fixture, parseConanlock, expected, expectedRelationships)
}
