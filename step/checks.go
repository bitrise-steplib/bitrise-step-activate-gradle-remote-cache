package step

import (
	"fmt"
	"os"
)

const unavailableAnnotation = `You have added the **Activate Build Cache** add-on step to your workflow.

However, it has not been activated for this workspace yet. Please contact [support@bitrise.io](mailto:support@bitrise.io) to activate it.

Build cache is not going to be activated in this build.`

const gradleHome = "~/.gradle"

func (step RemoteCacheStep) ensureFeatureEnabled() (bool, error) {
	isEnabled := step.envRepo.Get("BITRISEIO_BUILD_CACHE_ENABLED") == "true"

	if !isEnabled {
		step.logger.Warnf(unavailableAnnotation)

		cmd := step.commandFactory.Create("bitrise", []string{
			":annotations",
			"annotate",
			unavailableAnnotation,
			"--style", "warning",
		}, nil)
		cmdOut, err := cmd.RunAndReturnTrimmedCombinedOutput()
		if err != nil {
			return false, fmt.Errorf("failed to create annotation: %s", cmdOut)
		}
	}

	return true, nil
}

func (step RemoteCacheStep) ensureGradleHome() error {
	gradleHome, err := step.pathModifier.AbsPath(gradleHome)
	if err != nil {
		return err
	}
	err = os.MkdirAll(gradleHome, 0755)
	if err != nil {
		return err
	}
	return nil
}
