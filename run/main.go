package main

import (
	"os"

	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/scribe"
	"github.com/paketo-community/thin"
)

func main() {
	parser := thin.NewGemfileParser()
	logger := scribe.NewLogger(os.Stdout)

	packit.Run(
		thin.Detect(parser),
		thin.Build(logger),
	)
}
