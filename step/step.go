package step

import (
	"fmt"
	"os"
	"path"

	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-io/go-utils/v2/pathutil"
)

const gradleDepVersion = "main-SNAPSHOT"                   // TODO: we should change this to "1.+" after publishing 1.0
const apiEndpoint = "grpcs://cache-v3.bitrise.flare.build" // TODO: set up Secret Manager based on env

type Input struct {
	Verbose bool `env:"verbose,required"`
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

	token := step.envRepo.Get("BITRISEIO_BITRISE_SERVICES_ACCESS_TOKEN")
	if token == "" {
		return fmt.Errorf("$BITRISEIO_BITRISE_SERVICES_ACCESS_TOKEN is empty. This step is only supposed to run in Bitrise CI builds")
	}

	step.logger.Printf("Adding Gradle init script to ~/.gradle/init.gradle")
	if err := step.addInitScript(gradleDepVersion, apiEndpoint, token); err != nil {
		return fmt.Errorf("failed to set up remote caching: %w", err)
	}
	step.logger.Donef("Init script added, remote cache enabled for subsequent builds")

	return nil
}

func (step RemoteCacheStep) addInitScript(version, endpoint, token string) error {
	inventory := templateInventory{
		Version:   version,
		Endpoint:  endpoint,
		AuthToken: token,
	}
	scriptContent, err := renderTemplate(inventory)
	if err != nil {
		return err
	}

	gradleHome, err := step.pathModifier.AbsPath("~/.gradle")
	if err != nil {
		return err
	}
	err = os.MkdirAll(gradleHome, 0755)
	if err != nil {
		return fmt.Errorf("failed to create .gradle directory in user home: %w", err)
	}
	err = os.WriteFile(path.Join(gradleHome, "init.gradle"), []byte(scriptContent), 0755)
	if err != nil {
		return fmt.Errorf("failed to add init script to %s, error: %w", gradleHome, err)
	}

	return nil
}
