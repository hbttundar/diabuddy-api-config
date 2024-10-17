package appconfig_test

import (
	appconfig "github.com/hbttundar/diabuddy-api-config/config/app"
	"github.com/hbttundar/diabuddy-api-config/config/envmanager"
	testmain "github.com/hbttundar/diabuddy-api-config/test"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestAppConfig(t *testing.T) {
	tests := []struct {
		name         string
		envVariables map[string]string
		expected     string
		testFunction func(appConfig *appconfig.AppConfig, t *testing.T, expected string)
	}{
		{
			name:         "Get() method with defined key from env resolver",
			envVariables: map[string]string{"APP_NAME": "diabuddy-user-api-test"},
			expected:     "diabuddy-user-api-test",
			testFunction: func(appConfig *appconfig.AppConfig, t *testing.T, expected string) {
				value := appConfig.Get("APP_NAME")
				assert.Equal(t, expected, value, "Expected APP_NAME to be 'diabuddy-user-api-test'")
			},
		},
		{
			name:         "Get() method with defined key from os environment",
			envVariables: map[string]string{"APP_NAME": "TestApp"},
			expected:     "TestApp",
			testFunction: func(appConfig *appconfig.AppConfig, t *testing.T, expected string) {
				value := appConfig.Get("APP_NAME")
				assert.Equal(t, expected, value, "Expected APP_NAME to be 'TestApp' from environment")
			},
		},
		{
			name:         "Get() method when both os env and default value exist but env is empty",
			envVariables: map[string]string{"APP_NAME": ""},
			expected:     "Diabuddy",
			testFunction: func(appConfig *appconfig.AppConfig, t *testing.T, expected string) {
				value := appConfig.Get("APP_NAME")
				assert.Equal(t, expected, value, "Expected APP_NAME to fall back to default from Default() map")
			},
		},
		{
			name:         "Get() method with undefined key",
			envVariables: nil,
			expected:     "",
			testFunction: func(appConfig *appconfig.AppConfig, t *testing.T, expected string) {
				value := appConfig.Get("NON_EXISTENT_KEY")
				assert.Equal(t, expected, value, "Expected default for an undefined key to be an empty string")
			},
		},
		{
			name:         "Get() method with passing default value",
			envVariables: map[string]string{"APP_NAME": ""},
			expected:     "MyDefaultApp",
			testFunction: func(appConfig *appconfig.AppConfig, t *testing.T, expected string) {
				value := appConfig.Get("APP_NAME", "MyDefaultApp")
				assert.Equal(t, expected, value, "Expected APP_NAME to return default 'MyDefaultApp' value which was passed")
			},
		},
		{
			name:         "BasePath() method",
			envVariables: nil,
			expected:     "",
			testFunction: func(appConfig *appconfig.AppConfig, t *testing.T, expected string) {
				basePath, err := appConfig.BasePath()
				assert.NoError(t, err, "Expected no error while fetching base path")
				assert.NotEmpty(t, basePath, "Expected base path to not be empty")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.envVariables != nil {
				for key, value := range tt.envVariables {
					testmain.EnvVars[key] = value
				}
			}
			testmain.Setup()
			defer testmain.TearDown()
			envManager, err := envmanager.NewEnvManager()
			// Create a new AppConfig and execute the test function
			appConfig, err := appconfig.NewAppConfig(envManager)
			assert.NoError(t, err, "Unexpected error while creating AppConfig")
			tt.testFunction(appConfig, t, tt.expected)
		})
	}
}

func TestAppConfig_Validate(t *testing.T) {
	t.Run("Valid configuration with all required fields", func(t *testing.T) {
		// Setup environment variables for a valid configuration
		envVariables := map[string]string{
			"APP_NAME":  "Diabuddy",
			"APP_ENV":   "production",
			"APP_URL":   "http://localhost",
			"APP_DEBUG": "true",
		}

		// Set the environment variables for the test
		for key, value := range envVariables {
			testmain.EnvVars[key] = value
		}
		testmain.Setup()
		defer testmain.TearDown()
		// Create a new EnvManager with default values enabled
		envManager, err := envmanager.NewEnvManager(envmanager.WithUseDefault(true))
		assert.NoError(t, err, "Expected no error during env manager initialization")

		// Create a new AppConfig
		appConfig, err := appconfig.NewAppConfig(envManager)
		assert.NoError(t, err, "Expected no error during app c initialization")

		// Validate the configuration
		err = appConfig.Validate()
		assert.NoError(t, err, "Expected no validation error for valid c")
	})

	t.Run("Missing required environment variable (APP_NAME)", func(t *testing.T) {
		// Setup environment variables without "APP_NAME" to trigger failure
		envVariables := map[string]string{
			"APP_ENV":   "production",
			"APP_URL":   "http://localhost",
			"APP_DEBUG": "true",
		}

		// Set the environment variables for the test
		for key, value := range envVariables {
			testmain.EnvVars[key] = value
		}
		testmain.Setup()
		defer testmain.TearDown()

		// Create a new EnvManager without default values (useDefaults = false)
		envManager, err := envmanager.NewEnvManager(envmanager.WithUseDefault(false))
		assert.NoError(t, err, "Expected no error during env manager initialization")

		// Create a new AppConfig
		appConfig, err := appconfig.NewAppConfig(envManager)
		assert.NoError(t, err, "Expected no error during app c initialization")
		// force missing APP_NAME
		os.Unsetenv("APP_NAME")
		// Validate the configuration, expecting an error because APP_NAME is missing
		err = appConfig.Validate()
		assert.Error(t, err, "Expected validation error due to missing APP_NAME")
		if err != nil {
			assert.Contains(t, err.Error(), "APP_NAME is required", "Expected error message for missing APP_NAME")
		}
	})
}
