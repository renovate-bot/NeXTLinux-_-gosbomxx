package cpp

import (
	"encoding/json"
	"strings"

	"github.com/nextlinux/gosbom/gosbom/artifact"
	"github.com/nextlinux/gosbom/gosbom/pkg"
	"github.com/nextlinux/gosbom/gosbom/pkg/cataloger/generic"
	"github.com/nextlinux/gosbom/gosbom/source"
)

var _ generic.Parser = parseConanlock

type conanLock struct {
	GraphLock struct {
		Nodes map[string]struct {
			Ref            string   `json:"ref"`
			PackageID      string   `json:"package_id"`
			Context        string   `json:"context"`
			Prev           string   `json:"prev"`
			Requires       []string `json:"requires"`
			PythonRequires string   `json:"py_requires"`
			Options        string   `json:"options"`
			Path           string   `json:"path"`
		} `json:"nodes"`
	} `json:"graph_lock"`
	Version     string `json:"version"`
	ProfileHost string `json:"profile_host"`
}

// parseConanlock is a parser function for conan.lock contents, returning all packages discovered.
func parseConanlock(_ source.FileResolver, _ *generic.Environment, reader source.LocationReadCloser) ([]pkg.Package, []artifact.Relationship, error) {
	var pkgs []pkg.Package
	var cl conanLock
	if err := json.NewDecoder(reader).Decode(&cl); err != nil {
		return nil, nil, err
	}
	for _, node := range cl.GraphLock.Nodes {
		metadata := pkg.ConanLockMetadata{
			Ref:     node.Ref,
			Options: parseOptions(node.Options),
			Path:    node.Path,
			Context: node.Context,
		}

		p := newConanlockPackage(
			metadata,
			reader.Location.WithAnnotation(pkg.EvidenceAnnotationKey, pkg.PrimaryEvidenceAnnotation),
		)

		if p != nil {
			pkgs = append(pkgs, *p)
		}
	}

	return pkgs, nil, nil
}

func parseOptions(options string) map[string]string {
	o := make(map[string]string)
	if len(options) == 0 {
		return nil
	}

	kvps := strings.Split(options, "\n")
	for _, kvp := range kvps {
		kv := strings.Split(kvp, "=")
		if len(kv) == 2 {
			o[kv[0]] = kv[1]
		}
	}

	return o
}
