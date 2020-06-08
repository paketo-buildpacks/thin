package thin_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnitThin(t *testing.T) {
	suite := spec.New("thin", spec.Report(report.Terminal{}), spec.Parallel())
	suite("Build", testBuild)
	suite("Detect", testDetect)
	suite("GemfileParser", testGemfileParser)
	suite.Run(t)
}
