package apiconfig_test

import (
	apiconfig "github.com/hbttundar/diabuddy-api-config/config/api"
	"github.com/hbttundar/diabuddy-api-config/config/envmanager"
	testmain "github.com/hbttundar/diabuddy-api-config/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestNewApiConfig refactored to use global Setup and TearDown
func TestNewApiConfig(t *testing.T) {
	tests := []struct {
		name              string
		setupEnv          map[string]string
		useDefaultOptions bool
		expectedAppEnv    string
		expectedDBHost    string
		expectedError     bool
	}{
		{
			name: "Default configuration without environment variables",
			setupEnv: map[string]string{
				"APP_ENV": "",
				"DB_HOST": "",
			},
			useDefaultOptions: true,
			expectedAppEnv:    "local",
			expectedDBHost:    "127.0.0.1",
			expectedError:     false,
		},
		{
			name: "Custom configuration with environment variables",
			setupEnv: map[string]string{
				"APP_NAME":    "diabuddy-user-api",
				"APP_ENV":     "production",
				"APP_URL":     "localhost",
				"APP_DEBUG":   "false",
				"DB_HOST":     "192.168.10.10",
				"DB_PORT":     "5432",
				"DB_DATABASE": "diabuddy",
			},
			useDefaultOptions: false,
			expectedAppEnv:    "production",
			expectedDBHost:    "192.168.10.10",
			expectedError:     false,
		},
		{
			name:              "Missing configuration causing validation failure",
			setupEnv:          map[string]string{"APP_NAME": "", "APP_URL": "", "APP_DEBUG": "", "APP_ENV": ""},
			useDefaultOptions: false,
			expectedAppEnv:    "",
			expectedDBHost:    "",
			expectedError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Modify global EnvVars based on specific test setup
			for k, v := range tt.setupEnv {
				testmain.EnvVars[k] = v
			}

			// Call Setup() to apply the environment variables
			testmain.Setup()
			defer testmain.TearDown() // Clean up after the test

			envManager, err := envmanager.NewEnvManager(envmanager.WithUseDefault(tt.useDefaultOptions))
			assert.NoError(t, err, "Expect no error during env manager initialization")

			apiConfig, err := apiconfig.NewApiConfig(envManager)
			if tt.expectedError {
				assert.Error(t, err, "Expected an error due to missing configuration.")
			} else {
				assert.NoError(t, err, "Did not expect an error with valid configuration.")
				assert.Equal(t, tt.expectedAppEnv, apiConfig.App.Get("APP_ENV"), "Expected APP_ENV to match.")
				assert.Equal(t, tt.expectedDBHost, apiConfig.DB.Get("DB_HOST", "127.0.0.1"), "Expected DB_HOST to match.")
			}
		})
	}
}

// TestApiConfig_Validate refactored to use global Setup and TearDown
func TestApiConfig_Validate(t *testing.T) {
	tests := []struct {
		name              string
		setupEnv          map[string]string
		useDefaultOptions bool
		expectedError     bool
		expectedErrMsg    string
	}{
		{
			name: "Valid configuration with required fields",
			setupEnv: map[string]string{
				"APP_NAME":  "Diabuddy",
				"APP_ENV":   "production",
				"APP_URL":   "http://localhost",
				"APP_DEBUG": "true",
			},
			useDefaultOptions: true,
			expectedError:     false,
		},
		{
			name: "Missing required environment variable (APP_NAME)",
			setupEnv: map[string]string{
				"APP_NAME": "",
				"APP_ENV":  "production",
			},
			expectedError:     true,
			useDefaultOptions: false,
			expectedErrMsg:    "Error 400: APP_NAME is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Modify global EnvVars based on specific test setup
			for k, v := range tt.setupEnv {
				testmain.EnvVars[k] = v
			}

			// Call Setup() to apply the environment variables
			testmain.Setup()
			defer testmain.TearDown() // Clean up after the test

			envManager, err := envmanager.NewEnvManager(envmanager.WithUseDefault(tt.useDefaultOptions))
			assert.NoError(t, err, "Expect no error during env manager initialization")
			_, err = apiconfig.NewApiConfig(envManager)
			if tt.expectedError {
				assert.Error(t, err, tt.name)
				assert.Equal(t, tt.expectedErrMsg, err.Error(), tt.name)
			} else {
				assert.NoError(t, err, tt.name)
			}
		})
	}
}

// TestDBConfig_Validate refactored to use global Setup and TearDown
func TestDBConfig_Validate(t *testing.T) {
	tests := []struct {
		name              string
		setupEnv          map[string]string
		useDefaultOptions bool
		expectedError     bool
		expectedErrMsg    string
	}{
		{
			name: "Valid DB configuration with required fields",
			setupEnv: map[string]string{
				"APP_NAME":    "Diabuddy",
				"APP_ENV":     "production",
				"APP_URL":     "http://localhost",
				"APP_DEBUG":   "true",
				"DB_HOST":     "localhost",
				"DB_PORT":     "5432",
				"DB_DATABASE": "diabuddy",
				"DB_USERNAME": "root",
				"DB_PASSWORD": "password",
			},
			useDefaultOptions: false,
			expectedError:     false,
		},
		{
			name: "Missing required DB environment variable DB_HOST",
			setupEnv: map[string]string{
				"APP_NAME":     "Diabuddy",
				"APP_ENV":      "production",
				"APP_URL":      "http://localhost",
				"APP_DEBUG":    "true",
				"DATABASE_URL": "",
				"DB_PORT":      "5432",
				"DB_HOST":      "",
				"DB_DATABASE":  "diabuddy",
				"DB_USERNAME":  "root",
				"DB_PASSWORD":  "password",
			},
			useDefaultOptions: false,
			expectedError:     true,
			expectedErrMsg:    "Error 500: db configuration is invalid required key(s):DB_HOST is/are missed or empty.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Modify global EnvVars based on specific test setup
			for k, v := range tt.setupEnv {
				testmain.EnvVars[k] = v
			}

			// Call Setup() to apply the environment variables
			testmain.Setup()
			defer testmain.TearDown() // Clean up after the test

			envManager, err := envmanager.NewEnvManager(envmanager.WithUseDefault(tt.useDefaultOptions))
			assert.NoError(t, err, "Expect no error during env manager initialization")
			_, err = apiconfig.NewApiConfig(envManager)
			if tt.expectedError {
				assert.Error(t, err, tt.name)
			} else {
				assert.NoError(t, err, tt.name)
			}
		})
	}
}
