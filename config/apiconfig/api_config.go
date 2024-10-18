package apiconfig

import (
	"github.com/hbttundar/diabuddy-api-config/config"
	appconfig "github.com/hbttundar/diabuddy-api-config/config/appconfig"
	dbconfig "github.com/hbttundar/diabuddy-api-config/config/dbconfig"
	"github.com/hbttundar/diabuddy-api-config/config/envmanager"
	diabuddyErrors "github.com/hbttundar/diabuddy-errors"
)

type ApiConfig struct {
	DB  config.Config
	App config.Config
}

func NewApiConfig(envManager *envmanager.EnvManager) (*ApiConfig, diabuddyErrors.ApiErrors) {
	appConfig, err := appconfig.NewAppConfig(envManager)
	if err != nil {
		return nil, err
	}
	dbConfig, err := dbconfig.NewDBConfig(envManager)
	if err != nil {
		return nil, err
	}

	apiConfig := &ApiConfig{
		App: appConfig,
		DB:  dbConfig,
	}
	err = apiConfig.Validate()
	if err != nil {
		return nil, err
	}
	return apiConfig, nil
}

func (ac *ApiConfig) Validate() diabuddyErrors.ApiErrors {
	if err := ac.App.Validate(); err != nil {
		return err
	}
	if err := ac.DB.Validate(); err != nil {
		return err
	}
	return nil
}
