package thin

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

type GemfileParser struct{}

func NewGemfileParser() GemfileParser {
	return GemfileParser{}
}

func (p GemfileParser) Parse(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, fmt.Errorf("failed to parse Gemfile: %w", err)
	}
	defer file.Close()

	quotes := `["']`
	thinRe := regexp.MustCompile(fmt.Sprintf(`^gem %sthin%s`, quotes, quotes))

	hasThin := false
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := []byte(scanner.Text())

		if !hasThin {
			hasThin = thinRe.Match(line)
		}

		if hasThin {
			return true, nil
		}
	}

	return false, nil
}
