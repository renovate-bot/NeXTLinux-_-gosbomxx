package erlang

import (
	"github.com/nextlinux/gosbom/gosbom/pkg"
	"github.com/nextlinux/gosbom/gosbom/source"
	"github.com/nextlinux/packageurl-go"
)

func newPackage(d pkg.RebarLockMetadata, locations ...source.Location) pkg.Package {
	p := pkg.Package{
		Name:         d.Name,
		Version:      d.Version,
		Language:     pkg.Erlang,
		Locations:    source.NewLocationSet(locations...),
		PURL:         packageURL(d),
		Type:         pkg.HexPkg,
		MetadataType: pkg.RebarLockMetadataType,
		Metadata:     d,
	}

	p.SetID()

	return p
}

func packageURL(m pkg.RebarLockMetadata) string {
	var qualifiers packageurl.Qualifiers

	return packageurl.NewPackageURL(
		packageurl.TypeHex,
		"",
		m.Name,
		m.Version,
		qualifiers,
		"",
	).ToString()
}
