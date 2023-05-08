package elixir

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"

	"github.com/nextlinux/gosbom/internal/log"
	"github.com/nextlinux/gosbom/gosbom/artifact"
	"github.com/nextlinux/gosbom/gosbom/pkg"
	"github.com/nextlinux/gosbom/gosbom/pkg/cataloger/generic"
	"github.com/nextlinux/gosbom/gosbom/source"
)

// integrity check
var _ generic.Parser = parseMixLock

var mixLockDelimiter = regexp.MustCompile(`[%{}\n" ,:]+`)

// parseMixLock parses a mix.lock and returns the discovered Elixir packages.
func parseMixLock(_ source.FileResolver, _ *generic.Environment, reader source.LocationReadCloser) ([]pkg.Package, []artifact.Relationship, error) {
	r := bufio.NewReader(reader)

	var packages []pkg.Package
	for {
		line, err := r.ReadString('\n')
		switch {
		case errors.Is(io.EOF, err):
			return packages, nil, nil
		case err != nil:
			return nil, nil, fmt.Errorf("failed to parse mix.lock file: %w", err)
		}
		tokens := mixLockDelimiter.Split(line, -1)
		if len(tokens) < 6 {
			continue
		}
		name, version, hash, hashExt := tokens[1], tokens[4], tokens[5], tokens[len(tokens)-2]

		if name == "" {
			log.WithFields("path", reader.RealPath).Debug("skipping empty package name from mix.lock file")
			continue
		}

		packages = append(packages,
			newPackage(
				pkg.MixLockMetadata{
					Name:       name,
					Version:    version,
					PkgHash:    hash,
					PkgHashExt: hashExt,
				},
				reader.Location.WithAnnotation(pkg.EvidenceAnnotationKey, pkg.PrimaryEvidenceAnnotation),
			),
		)
	}
}
