package step

import (
	"fmt"
	"os"
)

func (step RemoteCacheStep) ensureFeatureEnabled() error {
	isEnabled := step.envRepo.Get("BITRISEIO_BUILD_CACHE_ENABLED") == "true"

	if !isEnabled {
		cmd := step.commandFactory.Create("bitrise", []string{
			":annotations",
			"annotate",
			"Test message",
			"--style", "error",
		}, nil)
		cmdOut, err := cmd.RunAndReturnTrimmedCombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to create annotation: %s", cmdOut)
		}
	}

	return nil
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
