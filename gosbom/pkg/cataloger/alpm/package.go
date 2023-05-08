package alpm

import (
	"github.com/nextlinux/gosbom/gosbom/linux"
	"github.com/nextlinux/gosbom/gosbom/pkg"
	"github.com/nextlinux/gosbom/gosbom/source"
	"github.com/nextlinux/gosbom/internal"
	"github.com/nextlinux/packageurl-go"
)

func newPackage(m pkg.AlpmMetadata, release *linux.Release, locations ...source.Location) pkg.Package {
	p := pkg.Package{
		Name:         m.Package,
		Version:      m.Version,
		Locations:    source.NewLocationSet(locations...),
		Type:         pkg.AlpmPkg,
		Licenses:     internal.SplitAny(m.License, " \n"),
		PURL:         packageURL(m, release),
		MetadataType: pkg.AlpmMetadataType,
		Metadata:     m,
	}
	p.SetID()
	return p
}

func packageURL(m pkg.AlpmMetadata, distro *linux.Release) string {
	if distro == nil || distro.ID != "arch" {
		// note: there is no namespace variation (like with debian ID_LIKE for ubuntu ID, for example)
		return ""
	}

	qualifiers := map[string]string{
		pkg.PURLQualifierArch: m.Architecture,
	}

	if m.BasePackage != "" {
		qualifiers[pkg.PURLQualifierUpstream] = m.BasePackage
	}

	return packageurl.NewPackageURL(
		"alpm", // `alpm` for Arch Linux and other users of the libalpm/pacman package manager. (see https://github.com/package-url/purl-spec/pull/164)
		distro.ID,
		m.Package,
		m.Version,
		pkg.PURLQualifiers(
			qualifiers,
			distro,
		),
		"",
	).ToString()
}
