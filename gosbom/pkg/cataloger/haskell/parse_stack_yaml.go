package haskell

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v3"

	"github.com/nextlinux/gosbom/gosbom/artifact"
	"github.com/nextlinux/gosbom/gosbom/pkg"
	"github.com/nextlinux/gosbom/gosbom/pkg/cataloger/generic"
	"github.com/nextlinux/gosbom/gosbom/source"
)

var _ generic.Parser = parseStackYaml

type stackYaml struct {
	ExtraDeps []string `yaml:"extra-deps"`
}

// parseStackYaml is a parser function for stack.yaml contents, returning all packages discovered.
func parseStackYaml(_ source.FileResolver, _ *generic.Environment, reader source.LocationReadCloser) ([]pkg.Package, []artifact.Relationship, error) {
	bytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load stack.yaml file: %w", err)
	}

	var stackFile stackYaml

	if err := yaml.Unmarshal(bytes, &stackFile); err != nil {
		return nil, nil, fmt.Errorf("failed to parse stack.yaml file: %w", err)
	}

	var pkgs []pkg.Package
	for _, dep := range stackFile.ExtraDeps {
		pkgName, pkgVersion, pkgHash := parseStackPackageEncoding(dep)
		pkgs = append(
			pkgs,
			newPackage(
				pkgName,
				pkgVersion,
				&pkg.HackageMetadata{
					PkgHash: pkgHash,
				},
				reader.Location.WithAnnotation(pkg.EvidenceAnnotationKey, pkg.PrimaryEvidenceAnnotation),
			),
		)
	}

	return pkgs, nil, nil
}
