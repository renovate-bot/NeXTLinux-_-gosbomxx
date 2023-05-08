package dotnet

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/nextlinux/gosbom/gosbom/artifact"
	"github.com/nextlinux/gosbom/gosbom/pkg"
	"github.com/nextlinux/gosbom/gosbom/pkg/cataloger/generic"
	"github.com/nextlinux/gosbom/gosbom/source"
)

var _ generic.Parser = parseDotnetDeps

type dotnetDeps struct {
	Libraries map[string]dotnetDepsLibrary `json:"libraries"`
}

type dotnetDepsLibrary struct {
	Type     string `json:"type"`
	Path     string `json:"path"`
	Sha512   string `json:"sha512"`
	HashPath string `json:"hashPath"`
}

func parseDotnetDeps(_ source.FileResolver, _ *generic.Environment, reader source.LocationReadCloser) ([]pkg.Package, []artifact.Relationship, error) {
	var pkgs []pkg.Package

	dec := json.NewDecoder(reader)

	var p dotnetDeps
	if err := dec.Decode(&p); err != nil {
		return nil, nil, fmt.Errorf("failed to parse deps.json file: %w", err)
	}

	var names []string

	for nameVersion := range p.Libraries {
		names = append(names, nameVersion)
	}

	// sort the names so that the order of the packages is deterministic
	sort.Strings(names)

	for _, nameVersion := range names {
		lib := p.Libraries[nameVersion]
		dotnetPkg := newDotnetDepsPackage(
			nameVersion,
			lib,
			reader.Location.WithAnnotation(pkg.EvidenceAnnotationKey, pkg.PrimaryEvidenceAnnotation),
		)

		if dotnetPkg != nil {
			pkgs = append(pkgs, *dotnetPkg)
		}
	}

	return pkgs, nil, nil
}
