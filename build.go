package thin

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/fs"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

func Build(logger scribe.Emitter) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		logger.Title("%s %s", context.BuildpackInfo.Name, context.BuildpackInfo.Version)

		rackConfigFilepath := filepath.Join(context.WorkingDir, "config.ru")
		thinConfigFilepath := os.Getenv("BP_THIN_CONFIG_LOCATION")
		if thinConfigFilepath != "" {
			if !filepath.IsAbs(thinConfigFilepath) {
				thinConfigFilepath = filepath.Join(context.WorkingDir, thinConfigFilepath)
			}

			thinConfigFilepathExists, err := fs.Exists(thinConfigFilepath)
			if err != nil {
				return packit.BuildResult{}, err
			}

			if !thinConfigFilepathExists {
				return packit.BuildResult{}, packit.Fail.WithMessage("thin config file does not exist at: %s", thinConfigFilepath)
			}
		} else {
			thinConfigFilepath = filepath.Join(context.WorkingDir, "thin.yml")
		}

		args := "bundle exec thin"

		exists, err := fs.Exists(thinConfigFilepath)
		if err != nil {
			return packit.BuildResult{}, err
		}

		if exists {
			args = args + fmt.Sprintf(" -C %s", thinConfigFilepath)
		}

		exists, err = fs.Exists(rackConfigFilepath)
		if err != nil {
			return packit.BuildResult{}, err
		}

		if exists {
			args = args + fmt.Sprintf(" -R %s", rackConfigFilepath)
		}

		// 3000 is the default thin port
		args = args + ` -p "${PORT:-3000}" start`
		processes := []packit.Process{
			{
				Type:    "web",
				Command: "bash",
				Args:    []string{"-c", args},
				Default: true,
				Direct:  true,
			},
		}
		logger.LaunchProcesses(processes)

		return packit.BuildResult{
			Launch: packit.LaunchMetadata{
				Processes: processes,
			},
		}, nil
	}
}
