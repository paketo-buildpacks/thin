module github.com/paketo-buildpacks/thin

go 1.16

require (
	github.com/BurntSushi/toml v1.2.1
	github.com/onsi/gomega v1.27.4
	github.com/opencontainers/runc v1.1.4 // indirect
	github.com/paketo-buildpacks/occam v0.16.0
	github.com/paketo-buildpacks/packit/v2 v2.9.0
	github.com/sclevine/spec v1.4.0
	gotest.tools/v3 v3.4.0 // indirect
)

replace github.com/CycloneDX/cyclonedx-go => github.com/CycloneDX/cyclonedx-go v0.6.0
