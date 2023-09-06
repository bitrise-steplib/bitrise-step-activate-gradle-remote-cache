package step

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-io/go-utils/v2/pathutil"
)

// Sync the major version of this step and the library.
// Use the latest 1.x version of our dependency, so we don't have to update this definition after every lib release.
// But don't forget to update this to `2.+` if the library reaches version 2.0!
const gradleDepVersion = "1.+"
const apiEndpoint = "grpcs://pluggable.services.bitrise.io"

// Sync the major version of this step and the plugin.
// Use the latest 1.x version of the plugin, so we don't have to update this definition after every plugin release.
// But don't forget to update this to `2.+` if the library reaches version 2.0!
const metricsDepVersion = "main-SNAPSHOT" // TODO: change to 1.+
const metricsEndpoint = "gradle-analytics.services.bitrise.io"
const metricsPort = 443

type Input struct {
	Push            bool   `env:"push,required"`
	Verbose         bool   `env:"verbose,required"`
	ValidationLevel string `env:"validation_level,opt[none,warning,error]"`
	CollectMetrics  bool   `env:"collect_metrics"`
}

type RemoteCacheStep struct {
	logger         log.Logger
	inputParser    stepconf.InputParser
	commandFactory command.Factory
	envRepo        env.Repository
	pathModifier   pathutil.PathModifier
}

func New(logger log.Logger, inputParser stepconf.InputParser, commandFactory command.Factory, envRepo env.Repository, pathModifier pathutil.PathModifier) RemoteCacheStep {
	return RemoteCacheStep{
		logger:         logger,
		inputParser:    inputParser,
		commandFactory: commandFactory,
		envRepo:        envRepo,
		pathModifier:   pathModifier,
	}
}

func (step RemoteCacheStep) Run() error {
	var input Input
	if err := step.inputParser.Parse(&input); err != nil {
		return fmt.Errorf("failed to parse inputs: %w", err)
	}
	stepconf.Print(input)
	step.logger.Println()

	step.logger.EnableDebugLog(input.Verbose)

	enabled, err := step.ensureFeatureEnabled()
	if err != nil {
		step.logger.Warnf("failed to check if feature is available: %s", err)
		return nil
	}
	if !enabled {
		return nil
	}

	token := step.envRepo.Get("BITRISEIO_BITRISE_SERVICES_ACCESS_TOKEN")
	if token == "" {
		return fmt.Errorf("$BITRISEIO_BITRISE_SERVICES_ACCESS_TOKEN is empty. This step is only supposed to run in Bitrise CI builds")
	}

	step.logger.Printf("Adding Gradle init script to ~/.gradle/init.gradle")
	if err := step.ensureGradleHome(); err != nil {
		return fmt.Errorf("failed to create .gradle directory in user home: %w", err)
	}
	if err := step.addInitScript(gradleDepVersion, apiEndpoint, token, input); err != nil {
		return fmt.Errorf("failed to set up remote caching: %w", err)
	}
	if err := step.addGlobalGradleProperties(); err != nil {
		return fmt.Errorf("failed to apply additional Gradle properties: %w", err)
	}
	step.logger.Donef("Init script added, remote cache enabled for subsequent builds")

	return nil
}

func (step RemoteCacheStep) addInitScript(version, endpoint, token string, input Input) error {
	inventory := templateInventory{
		CacheVersion:    version,
		CacheEndpoint:   endpoint,
		AuthToken:       token,
		PushEnabled:     input.Push,
		DebugEnabled:    input.Verbose,
		ValidationLevel: input.ValidationLevel,
		MetricsEnabled:  input.CollectMetrics,
		MetricsVersion:  metricsDepVersion,
		MetricsEndpoint: metricsEndpoint,
		MetricsPort:     metricsPort,
	}
	scriptContent, err := renderTemplate(inventory)
	if err != nil {
		return err
	}

	gradleHome, err := step.pathModifier.AbsPath(gradleHome)
	if err != nil {
		return err
	}
	err = os.WriteFile(path.Join(gradleHome, "init.gradle"), []byte(scriptContent), 0755)
	if err != nil {
		return fmt.Errorf("failed to add init script to %s, error: %w", gradleHome, err)
	}

	return nil
}

// addGlobalGradleProperties creates additional settings at ~/.gradle/gradle.properties, overriding
// any properties defined in the project root directory.
// https://docs.gradle.org/current/userguide/build_environment.html#sec:gradle_configuration_properties
func (step RemoteCacheStep) addGlobalGradleProperties() error {
	gradleHome, err := step.pathModifier.AbsPath(gradleHome)
	if err != nil {
		return err
	}

	// Enable build caching - some projects enable this, but it's disabled by default
	// It needs to be enabled for the remote cache config to take effect
	gradleProperties := "org.gradle.caching=true"
	return os.WriteFile(filepath.Join(gradleHome, "gradle.properties"), []byte(gradleProperties), 0755)
}
