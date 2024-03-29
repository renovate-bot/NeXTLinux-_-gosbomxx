package elixir

import (
	"github.com/nextlinux/gosbom/gosbom/pkg"
	"github.com/nextlinux/gosbom/gosbom/source"
	"github.com/nextlinux/packageurl-go"
)

func newPackage(d pkg.MixLockMetadata, locations ...source.Location) pkg.Package {
	p := pkg.Package{
		Name:         d.Name,
		Version:      d.Version,
		Language:     pkg.Elixir,
		Locations:    source.NewLocationSet(locations...),
		PURL:         packageURL(d),
		Type:         pkg.HexPkg,
		MetadataType: pkg.MixLockMetadataType,
		Metadata:     d,
	}

	p.SetID()

	return p
}

func packageURL(m pkg.MixLockMetadata) string {
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
