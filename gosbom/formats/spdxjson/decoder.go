package spdxjson

import (
	"fmt"
	"io"

	"github.com/spdx/tools-golang/json"

	"github.com/nextlinux/gosbom/gosbom/formats/common/spdxhelpers"
	"github.com/nextlinux/gosbom/gosbom/sbom"
)

func decoder(reader io.Reader) (s *sbom.SBOM, err error) {
	doc, err := json.Read(reader)
	if err != nil {
		return nil, fmt.Errorf("unable to decode spdx-json: %w", err)
	}

	return spdxhelpers.ToGosbomModel(doc)
}
