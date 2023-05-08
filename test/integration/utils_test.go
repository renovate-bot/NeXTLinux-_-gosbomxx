package integration

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nextlinux/stereoscope/pkg/imagetest"
	"github.com/nextlinux/gosbom/gosbom"
	"github.com/nextlinux/gosbom/gosbom/pkg/cataloger"
	"github.com/nextlinux/gosbom/gosbom/sbom"
	"github.com/nextlinux/gosbom/gosbom/source"
)

func catalogFixtureImage(t *testing.T, fixtureImageName string, scope source.Scope, catalogerCfg []string) (sbom.SBOM, *source.Source) {
	imagetest.GetFixtureImage(t, "docker-archive", fixtureImageName)
	tarPath := imagetest.GetFixtureImageTarPath(t, fixtureImageName)
	userInput := "docker-archive:" + tarPath
	sourceInput, err := source.ParseInput(userInput, "")
	require.NoError(t, err)
	theSource, cleanupSource, err := source.New(*sourceInput, nil, nil)
	t.Cleanup(cleanupSource)
	require.NoError(t, err)

	c := cataloger.DefaultConfig()
	c.Catalogers = catalogerCfg

	c.Search.Scope = scope
	pkgCatalog, relationships, actualDistro, err := gosbom.CatalogPackages(theSource, c)
	if err != nil {
		t.Fatalf("failed to catalog image: %+v", err)
	}

	return sbom.SBOM{
		Artifacts: sbom.Artifacts{
			Packages:          pkgCatalog,
			LinuxDistribution: actualDistro,
		},
		Relationships: relationships,
		Source:        theSource.Metadata,
		Descriptor: sbom.Descriptor{
			Name:    "gosbom",
			Version: "v0.42.0-bogus",
			// the application configuration should be persisted here, however, we do not want to import
			// the application configuration in this package (it's reserved only for ingestion by the cmd package)
			Configuration: map[string]string{
				"config-key": "config-value",
			},
		},
	}, theSource
}

func catalogDirectory(t *testing.T, dir string) (sbom.SBOM, *source.Source) {
	userInput := "dir:" + dir
	sourceInput, err := source.ParseInput(userInput, "")
	require.NoError(t, err)
	theSource, cleanupSource, err := source.New(*sourceInput, nil, nil)
	t.Cleanup(cleanupSource)
	require.NoError(t, err)

	// TODO: this would be better with functional options (after/during API refactor)
	c := cataloger.DefaultConfig()
	c.Search.Scope = source.AllLayersScope
	pkgCatalog, relationships, actualDistro, err := gosbom.CatalogPackages(theSource, c)
	if err != nil {
		t.Fatalf("failed to catalog image: %+v", err)
	}

	return sbom.SBOM{
		Artifacts: sbom.Artifacts{
			Packages:          pkgCatalog,
			LinuxDistribution: actualDistro,
		},
		Relationships: relationships,
		Source:        theSource.Metadata,
	}, theSource
}
