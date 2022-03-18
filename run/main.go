package main

import (
	"os"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	"github.com/paketo-buildpacks/thin"
)

func main() {
	parser := thin.NewGemfileParser()
	logger := scribe.NewEmitter(os.Stdout)

	packit.Run(
		thin.Detect(parser),
		thin.Build(logger),
	)
}
