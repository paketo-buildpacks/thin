package thin_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	"github.com/paketo-buildpacks/thin"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testBuild(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		layersDir  string
		workingDir string
		cnbDir     string
		buffer     *bytes.Buffer

		build        packit.BuildFunc
		buildContext packit.BuildContext
	)

	it.Before(func() {
		var err error
		layersDir, err = os.MkdirTemp("", "layers")
		Expect(err).NotTo(HaveOccurred())

		cnbDir, err = os.MkdirTemp("", "cnb")
		Expect(err).NotTo(HaveOccurred())

		workingDir, err = os.MkdirTemp("", "working-dir")
		Expect(err).NotTo(HaveOccurred())

		buffer = bytes.NewBuffer(nil)
		logger := scribe.NewEmitter(buffer)

		build = thin.Build(logger)
		buildContext = packit.BuildContext{
			WorkingDir: workingDir,
			CNBPath:    cnbDir,
			Stack:      "some-stack",
			BuildpackInfo: packit.BuildpackInfo{
				Name:    "Some Buildpack",
				Version: "some-version",
			},
			Plan: packit.BuildpackPlan{
				Entries: []packit.BuildpackPlanEntry{},
			},
			Layers: packit.Layers{Path: layersDir},
		}
	})

	it.After(func() {
		Expect(os.RemoveAll(layersDir)).To(Succeed())
		Expect(os.RemoveAll(cnbDir)).To(Succeed())
		Expect(os.RemoveAll(workingDir)).To(Succeed())
	})

	it("returns a result that provides a thin start command", func() {
		result, err := build(buildContext)
		Expect(err).NotTo(HaveOccurred())

		Expect(result).To(Equal(packit.BuildResult{
			Plan: packit.BuildpackPlan{
				Entries: nil,
			},
			Layers: nil,
			Launch: packit.LaunchMetadata{
				Processes: []packit.Process{
					{
						Type:    "web",
						Command: "bash",
						Args:    []string{"-c", `bundle exec thin -p "${PORT:-3000}" start`},
						Default: true,
						Direct:  true,
					},
				},
			},
		}))

		Expect(buffer.String()).To(ContainSubstring("Some Buildpack some-version"))
		Expect(buffer.String()).To(ContainSubstring("Assigning launch processes:"))
	})

	context("when a thin.yml file exists in the working directory", func() {
		it.Before(func() {
			Expect(os.WriteFile(filepath.Join(workingDir, "thin.yml"), []byte{}, os.ModePerm)).To(Succeed())
		})

		it("returns a result with that file provided to the thin start command", func() {
			result, err := build(buildContext)
			Expect(err).NotTo(HaveOccurred())

			Expect(result).To(Equal(packit.BuildResult{
				Plan: packit.BuildpackPlan{
					Entries: nil,
				},
				Layers: nil,
				Launch: packit.LaunchMetadata{
					Processes: []packit.Process{
						{
							Type:    "web",
							Command: "bash",
							Args: []string{
								"-c",
								fmt.Sprintf(`bundle exec thin -C %s -p "${PORT:-3000}" start`, filepath.Join(workingDir, "thin.yml")),
							},
							Default: true,
							Direct:  true,
						},
					},
				},
			}))
		})
	})

	context("when the BP_THIN_CONFIG_LOCATION environment variable points to a valid file", func() {
		it.Before(func() {
			thinConfigFilepath := filepath.Join(workingDir, "some-thin-config.yml")
			Expect(os.WriteFile(thinConfigFilepath, []byte{}, os.ModePerm)).To(Succeed())
			Expect(os.Setenv("BP_THIN_CONFIG_LOCATION", thinConfigFilepath)).To(Succeed())
		})

		it.After(func() {
			Expect(os.Unsetenv("BP_THIN_CONFIG_LOCATION")).To(Succeed())
		})

		it("returns a result with that file provided to the thin start command", func() {
			result, err := build(buildContext)
			Expect(err).NotTo(HaveOccurred())

			Expect(result).To(Equal(packit.BuildResult{
				Plan: packit.BuildpackPlan{
					Entries: nil,
				},
				Layers: nil,
				Launch: packit.LaunchMetadata{
					Processes: []packit.Process{
						{
							Type:    "web",
							Command: "bash",
							Args: []string{
								"-c",
								fmt.Sprintf(`bundle exec thin -C %s -p "${PORT:-3000}" start`, filepath.Join(workingDir, "some-thin-config.yml")),
							},
							Default: true,
							Direct:  true,
						},
					},
				},
			}))
		})

		context("when the BP_THIN_CONFIG_LOCATION environment variable is a relative path", func() {
			it.Before(func() {
				relativeFilepath := filepath.Join("some-dir", "some-thin-config.yml")
				Expect(os.MkdirAll(filepath.Join(workingDir, "some-dir"), os.ModePerm)).To(Succeed())
				Expect(os.WriteFile(filepath.Join(workingDir, "some-dir", "some-thin-config.yml"), []byte{}, os.ModePerm)).To(Succeed())
				Expect(os.Setenv("BP_THIN_CONFIG_LOCATION", relativeFilepath)).To(Succeed())
			})

			it.After(func() {
				Expect(os.Unsetenv("BP_THIN_CONFIG_LOCATION")).To(Succeed())
			})

			it("returns a result with that file relative to the working directory", func() {
				result, err := build(buildContext)
				Expect(err).NotTo(HaveOccurred())

				Expect(result).To(Equal(packit.BuildResult{
					Plan: packit.BuildpackPlan{
						Entries: nil,
					},
					Layers: nil,
					Launch: packit.LaunchMetadata{
						Processes: []packit.Process{
							{
								Type:    "web",
								Command: "bash",
								Args: []string{
									"-c",
									fmt.Sprintf(`bundle exec thin -C %s -p "${PORT:-3000}" start`, filepath.Join(workingDir, "some-dir", "some-thin-config.yml")),
								},
								Default: true,
								Direct:  true,
							},
						},
					},
				}))
			})
		})
	})

	context("when both the BP_THIN_CONFIG_LOCATION env var is set and a thin.yml file is present", func() {
		it.Before(func() {
			Expect(os.WriteFile(filepath.Join(workingDir, "thin.yml"), []byte{}, os.ModePerm)).To(Succeed())

			envVarThinConfigFilepath := filepath.Join(workingDir, "some-thin-config.yml")
			Expect(os.WriteFile(envVarThinConfigFilepath, []byte{}, os.ModePerm)).To(Succeed())
			Expect(os.Setenv("BP_THIN_CONFIG_LOCATION", envVarThinConfigFilepath)).To(Succeed())
		})

		it.After(func() {
			Expect(os.Unsetenv("BP_THIN_CONFIG_LOCATION")).To(Succeed())
		})

		it("prioritizes the environment variable", func() {
			result, err := build(buildContext)
			Expect(err).NotTo(HaveOccurred())

			Expect(result).To(Equal(packit.BuildResult{
				Plan: packit.BuildpackPlan{
					Entries: nil,
				},
				Layers: nil,
				Launch: packit.LaunchMetadata{
					Processes: []packit.Process{
						{
							Type:    "web",
							Command: "bash",
							Args: []string{
								"-c",
								fmt.Sprintf(`bundle exec thin -C %s -p "${PORT:-3000}" start`, filepath.Join(workingDir, "some-thin-config.yml")),
							},
							Default: true,
							Direct:  true,
						},
					},
				},
			}))
		})
	})

	context("failure cases", func() {
		context("when there is an error determining if the default thin config file exists", func() {
			it.Before(func() {
				envVarThinConfigFilepath := filepath.Join(workingDir, "some-thin-config.yml")
				Expect(os.WriteFile(envVarThinConfigFilepath, []byte{}, os.ModePerm)).To(Succeed())
				Expect(os.Setenv("BP_THIN_CONFIG_LOCATION", envVarThinConfigFilepath)).To(Succeed())

				Expect(os.Chmod(workingDir, 0000)).To(Succeed())
			})

			it.After(func() {
				Expect(os.Chmod(workingDir, os.ModePerm)).To(Succeed())
			})

			it("returns an error", func() {
				_, err := build(buildContext)
				Expect(err).To(HaveOccurred())
			})
		})

		context("when there is an error determining if the env var thin config file exists", func() {
			it.Before(func() {
				envVarThinConfigFilepath := filepath.Join(workingDir, "some-thin-config.yml")
				Expect(os.WriteFile(envVarThinConfigFilepath, []byte{}, os.ModePerm)).To(Succeed())
				Expect(os.Setenv("BP_THIN_CONFIG_LOCATION", envVarThinConfigFilepath)).To(Succeed())

				Expect(os.Chmod(workingDir, 0000)).To(Succeed())
			})

			it.After(func() {
				Expect(os.Chmod(workingDir, os.ModePerm)).To(Succeed())
				Expect(os.Unsetenv("BP_THIN_CONFIG_LOCATION")).To(Succeed())
			})

			it("returns an error", func() {
				_, err := build(buildContext)
				Expect(err).To(HaveOccurred())
			})
		})

		context("when the BP_THIN_CONFIG_LOCATION environment variable points to a non-existent file", func() {
			it.Before(func() {
				Expect(os.Setenv("BP_THIN_CONFIG_LOCATION", filepath.Join(workingDir, "non-existent-file"))).To(Succeed())
			})

			it.After(func() {
				Expect(os.Unsetenv("BP_THIN_CONFIG_LOCATION")).To(Succeed())
			})

			it("returns an error", func() {
				_, err := build(buildContext)
				Expect(err).To(MatchError(packit.Fail.WithMessage("thin config file does not exist at: %s", filepath.Join(workingDir, "non-existent-file"))))
			})
		})
	})
}
