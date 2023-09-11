module github.com/paketo-buildpacks/thin

go 1.16

require (
	github.com/BurntSushi/toml v1.3.2
	github.com/onsi/gomega v1.27.10
	github.com/paketo-buildpacks/occam v0.17.0
	github.com/paketo-buildpacks/packit/v2 v2.11.0
	github.com/sclevine/spec v1.4.0
	golang.org/x/net v0.15.0 // indirect
	gotest.tools/v3 v3.4.0 // indirect
)

replace github.com/CycloneDX/cyclonedx-go => github.com/CycloneDX/cyclonedx-go v0.6.0
