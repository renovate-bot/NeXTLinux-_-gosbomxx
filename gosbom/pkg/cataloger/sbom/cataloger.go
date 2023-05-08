package sbom

import (
	"github.com/nextlinux/gosbom/internal/log"
	"github.com/nextlinux/gosbom/gosbom/artifact"
	"github.com/nextlinux/gosbom/gosbom/formats"
	"github.com/nextlinux/gosbom/gosbom/pkg"
	"github.com/nextlinux/gosbom/gosbom/pkg/cataloger/generic"
	"github.com/nextlinux/gosbom/gosbom/source"
)

const catalogerName = "sbom-cataloger"

// NewSBOMCataloger returns a new SBOM cataloger object loaded from saved SBOM JSON.
func NewSBOMCataloger() *generic.Cataloger {
	return generic.NewCataloger(catalogerName).
		WithParserByGlobs(parseSBOM,
			"**/*.gosbom.json",
			"**/*.bom.*",
			"**/*.bom",
			"**/bom",
			"**/*.sbom.*",
			"**/*.sbom",
			"**/sbom",
			"**/*.cdx.*",
			"**/*.cdx",
			"**/*.spdx.*",
			"**/*.spdx",
		)
}

func parseSBOM(_ source.FileResolver, _ *generic.Environment, reader source.LocationReadCloser) ([]pkg.Package, []artifact.Relationship, error) {
	s, _, err := formats.Decode(reader)
	if err != nil {
		return nil, nil, err
	}

	if s == nil {
		log.WithFields("path", reader.Location.RealPath).Trace("file is not an SBOM")
		return nil, nil, nil
	}

	var pkgs []pkg.Package
	var relationships []artifact.Relationship
	for _, p := range s.Artifacts.Packages.Sorted() {
		// replace all locations on the package with the location of the SBOM file.
		// Why not keep the original list of locations? Since the "locations" field is meant to capture
		// where there is evidence of this file, and the catalogers have not run against any file other than,
		// the SBOM, this is the only location that is relevant for this cataloger.
		p.Locations = source.NewLocationSet(
			reader.Location.WithAnnotation(pkg.EvidenceAnnotationKey, pkg.PrimaryEvidenceAnnotation),
		)
		p.FoundBy = catalogerName

		pkgs = append(pkgs, p)
		relationships = append(relationships, artifact.Relationship{
			From: p,
			To:   reader.Location.Coordinates,
			Type: artifact.DescribedByRelationship,
		})
	}

	return pkgs, relationships, nil
}
