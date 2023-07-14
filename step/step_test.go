package step

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
)

func TestRemoteCacheStep_Run(t *testing.T) {
	tests := []struct {
		name    string
		envRepo env.Repository
		wantErr bool
	}{
		{
			name: "happy path",
			envRepo: fakeEnvRepo{envVars: map[string]string{
				"BITRISEIO_BITRISE_SERVICES_ACCESS_TOKEN": "fake access token",
				"BITRISEIO_BUILD_CACHE_ENABLED": "true",
				"verbose":          "false",
				"push":             "true",
				"validation_level": "warning",
			}},
		},
		{
			name: "missing auth token",
			envRepo: fakeEnvRepo{envVars: map[string]string{
				"BITRISEIO_BUILD_CACHE_ENABLED": "true",
				"verbose":          "false",
				"push":             "true",
				"validation_level": "warning",
			}},
			wantErr: true,
		},
		{
			name: "feature not enabled",
			envRepo: fakeEnvRepo{envVars: map[string]string{
				"verbose":          "false",
				"push":             "true",
				"validation_level": "warning",
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakePath := t.TempDir()
			step := RemoteCacheStep{
				logger:         log.NewLogger(),
				inputParser:    stepconf.NewInputParser(tt.envRepo),
				commandFactory: command.NewFactory(tt.envRepo),
				envRepo:        tt.envRepo,
				pathModifier:   fakePathModifier{path: fakePath},
			}
			err := step.Run()
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			} else if err == nil {
				_, err := os.ReadFile(path.Join(fakePath, "init.gradle"))
				if err != nil {
					t.Errorf("failed to open generated file: %s", err)
				}
			}
		})
	}
}

type fakeEnvRepo struct {
	envVars map[string]string
}

func (repo fakeEnvRepo) Get(key string) string {
	value, ok := repo.envVars[key]
	if ok {
		return value
	} else {
		return ""
	}
}

func (repo fakeEnvRepo) Set(key, value string) error {
	repo.envVars[key] = value
	return nil
}

func (repo fakeEnvRepo) Unset(key string) error {
	repo.envVars[key] = ""
	return nil
}

func (repo fakeEnvRepo) List() []string {
	envs := []string{}
	for k, v := range repo.envVars {
		envs = append(envs, fmt.Sprintf("%s=%s", k, v))
	}
	return envs
}

type fakePathModifier struct {
	path string
}

func (pm fakePathModifier) AbsPath(pth string) (string, error) {
	return pm.path, nil
}
