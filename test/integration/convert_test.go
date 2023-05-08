package integration

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nextlinux/gosbom/cmd/gosbom/cli/convert"
	"github.com/nextlinux/gosbom/internal/config"
	"github.com/nextlinux/gosbom/gosbom/formats"
	"github.com/nextlinux/gosbom/gosbom/formats/cyclonedxjson"
	"github.com/nextlinux/gosbom/gosbom/formats/cyclonedxxml"
	"github.com/nextlinux/gosbom/gosbom/formats/spdxjson"
	"github.com/nextlinux/gosbom/gosbom/formats/spdxtagvalue"
	"github.com/nextlinux/gosbom/gosbom/formats/gosbomjson"
	"github.com/nextlinux/gosbom/gosbom/formats/table"
	"github.com/nextlinux/gosbom/gosbom/sbom"
	"github.com/nextlinux/gosbom/gosbom/source"
)

// TestConvertCmd tests if the converted SBOM is a valid document according
// to spec.
// TODO: This test can, but currently does not, check the converted SBOM content. It
// might be useful to do that in the future, once we gather a better understanding of
// what users expect from the convert command.
func TestConvertCmd(t *testing.T) {
	tests := []struct {
		name   string
		format sbom.Format
	}{
		{
			name:   "gosbom-json",
			format: gosbomjson.Format(),
		},
		{
			name:   "spdx-json",
			format: spdxjson.Format(),
		},
		{
			name:   "spdx-tag-value",
			format: spdxtagvalue.Format(),
		},
		{
			name:   "cyclonedx-json",
			format: cyclonedxjson.Format(),
		},
		{
			name:   "cyclonedx-xml",
			format: cyclonedxxml.Format(),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gosbomSbom, _ := catalogFixtureImage(t, "image-pkg-coverage", source.SquashedScope, nil)
			gosbomFormat := gosbomjson.Format()

			gosbomFile, err := os.CreateTemp("", "test-convert-sbom-")
			require.NoError(t, err)
			defer func() {
				_ = os.Remove(gosbomFile.Name())
			}()

			err = gosbomFormat.Encode(gosbomFile, gosbomSbom)
			require.NoError(t, err)

			formatFile, err := os.CreateTemp("", "test-convert-sbom-")
			require.NoError(t, err)
			defer func() {
				_ = os.Remove(gosbomFile.Name())
			}()

			ctx := context.Background()
			app := &config.Application{
				Outputs: []string{fmt.Sprintf("%s=%s", test.format.ID().String(), formatFile.Name())},
			}

			// stdout reduction of test noise
			rescue := os.Stdout // keep backup of the real stdout
			os.Stdout, _ = os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
			defer func() {
				os.Stdout = rescue
			}()

			err = convert.Run(ctx, app, []string{gosbomFile.Name()})
			require.NoError(t, err)
			contents, err := os.ReadFile(formatFile.Name())
			require.NoError(t, err)

			formatFound := formats.Identify(contents)
			if test.format.ID() == table.ID {
				require.Nil(t, formatFound)
				return
			}
			require.Equal(t, test.format.ID(), formatFound.ID())
		})
	}
}
