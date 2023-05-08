package java

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nextlinux/gosbom/gosbom/pkg"
)

func Test_packageURL(t *testing.T) {
	tests := []struct {
		pkg    pkg.Package
		expect string
	}{
		{
			pkg: pkg.Package{
				Name:         "example-java-app-maven",
				Version:      "0.1.0",
				Language:     pkg.Java,
				Type:         pkg.JavaPkg,
				MetadataType: pkg.JavaMetadataType,
				Metadata: pkg.JavaMetadata{
					VirtualPath: "test-fixtures/java-builds/packages/example-java-app-maven-0.1.0.jar",
					Manifest: &pkg.JavaManifest{
						Main: map[string]string{
							"Manifest-Version": "1.0",
						},
					},
					PomProperties: &pkg.PomProperties{
						Path:       "META-INF/maven/org.nextlinux/example-java-app-maven/pom.properties",
						GroupID:    "org.nextlinux",
						ArtifactID: "example-java-app-maven",
						Version:    "0.1.0",
						Extra:      make(map[string]string),
					},
				},
			},
			expect: "pkg:maven/org.nextlinux/example-java-app-maven@0.1.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.expect, func(t *testing.T) {
			assert.Equal(t, tt.expect, packageURL(tt.pkg.Name, tt.pkg.Version, tt.pkg.Metadata.(pkg.JavaMetadata)))
		})
	}
}
