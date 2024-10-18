package appconfig

import (
	"github.com/hbttundar/diabuddy-api-config/config/envmanager"
	"github.com/hbttundar/diabuddy-api-config/util/resolver/rootpath"
	diabuddyErrors "github.com/hbttundar/diabuddy-errors"
	"path/filepath"
)

type AppConfig struct {
	envManager   *envmanager.EnvManager
	pathResolver *rootpath.RootPathResolver
}

func NewAppConfig(envManager *envmanager.EnvManager) (*AppConfig, diabuddyErrors.ApiErrors) {
	ac := &AppConfig{
		envManager:   envManager,
		pathResolver: rootpath.NewRootPathResolver(),
	}
	return ac, nil
}

func (ac *AppConfig) Get(key string, defaultValue ...string) string {
	return ac.envManager.Get(key, defaultValue...)
}

func (ac *AppConfig) BasePath() (string, diabuddyErrors.ApiErrors) {
	basePath, err := filepath.Abs("../")
	if err != nil {
		return "", diabuddyErrors.NewApiError(diabuddyErrors.InternalServerErrorType, "could not find appconfig root directory", diabuddyErrors.WithInternalError(err))
	}
	return ac.pathResolver.Resolve(basePath)
}

func (ac *AppConfig) Validate() diabuddyErrors.ApiErrors {
	requiredKeys := []string{envmanager.AppNameKey, envmanager.AppEnvKey, envmanager.AppUrlKey, envmanager.AppDebugKey}
	for _, key := range requiredKeys {
		if ac.Get(key) == "" {
			return diabuddyErrors.NewApiError(diabuddyErrors.BadRequestErrorType, key+" is required")
		}
	}
	return nil
}
