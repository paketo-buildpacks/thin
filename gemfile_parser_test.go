package thin_test

import (
	"os"
	"testing"

	"github.com/paketo-buildpacks/thin"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testGemfileParser(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		path   string
		parser thin.GemfileParser
	)

	it.Before(func() {
		file, err := os.CreateTemp("", "Gemfile")
		Expect(err).NotTo(HaveOccurred())
		defer file.Close()

		path = file.Name()
		parser = thin.NewGemfileParser()
	})

	it.After(func() {
		Expect(os.RemoveAll(path)).To(Succeed())
	})

	context("Parse", func() {
		context("when using thin", func() {
			it("parses correctly", func() {
				Expect(os.WriteFile(path, []byte(`source 'https://rubygems.org'

gem 'thin'`), 0600)).To(Succeed())

				hasThin, err := parser.Parse(path)
				Expect(err).NotTo(HaveOccurred())
				Expect(hasThin).To(BeTrue())
			})
		})

		context("when not using thin", func() {
			it("parses correctly", func() {
				Expect(os.WriteFile(path, []byte(`source 'https://rubygems.org'
ruby '~> 2.0'`), 0600)).To(Succeed())

				hasThin, err := parser.Parse(path)
				Expect(err).NotTo(HaveOccurred())
				Expect(hasThin).To(BeFalse())
			})
		})

		context("when the Gemfile file does not exist", func() {
			it.Before(func() {
				Expect(os.Remove(path)).To(Succeed())
			})

			it("returns all false", func() {
				hasThin, err := parser.Parse(path)
				Expect(err).NotTo(HaveOccurred())
				Expect(hasThin).To(BeFalse())
			})
		})

		context("failure cases", func() {
			context("when the Gemfile cannot be opened", func() {
				it.Before(func() {
					Expect(os.Chmod(path, 0000)).To(Succeed())
				})

				it("returns an error", func() {
					_, err := parser.Parse(path)
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(ContainSubstring("failed to parse Gemfile:")))
					Expect(err).To(MatchError(ContainSubstring("permission denied")))
				})
			})
		})
	})
}
