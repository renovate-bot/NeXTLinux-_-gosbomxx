package integration

import (
	"testing"

	"github.com/nextlinux/gosbom/gosbom/source"
)

func Test_AllLayersIncludesSquashed(t *testing.T) {
	// This is a verification test for issue #894 (https://github.com/nextlinux/gosbom/issues/894)
	allLayers, _ := catalogFixtureImage(t, "image-suse-all-layers", source.AllLayersScope, nil)
	squashed, _ := catalogFixtureImage(t, "image-suse-all-layers", source.SquashedScope, nil)

	lenAllLayers := len(allLayers.Artifacts.Packages.Sorted())
	lenSquashed := len(squashed.Artifacts.Packages.Sorted())

	if lenAllLayers < lenSquashed {
		t.Errorf("squashed has more packages than all-layers: %d > %d", lenSquashed, lenAllLayers)
	}
}
