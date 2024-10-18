package envmanager_test

import (
	"github.com/hbttundar/diabuddy-api-config/config/envmanager"
	testmain "github.com/hbttundar/diabuddy-api-config/test"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewEnvManager(t *testing.T) {
	manager, err := envmanager.NewEnvManager()
	assert.NoError(t, err, "expected no error while creating EnvManager")
	assert.NotNil(t, manager, "expected a valid EnvManager instance")
}

func TestLoadEnvironmentVariables(t *testing.T) {
	t.Run("Load environment variables for test environment", func(t *testing.T) {
		testmain.EnvVars[envmanager.AppEnvKey] = "test"

		// Apply environment setup
		testmain.Setup()
		defer testmain.TearDown()

		envManager, err := envmanager.NewEnvManager(envmanager.WithEnvironment("test"))
		assert.NoError(t, err, "expected no error while creating env manager")

		err = envManager.LoadEnvironmentVariables()
		assert.NoError(t, err, "expected no error while loading environment variables for test environment")

		appName := envManager.Get(envmanager.AppNameKey)
		assert.NotEmpty(t, appName, "expected APP_NAME to be loaded from test environment")
	})

	t.Run("Not fail to load environment variables for non-existent environment as use its default .env file", func(t *testing.T) {
		testmain.EnvVars[envmanager.AppEnvKey] = "nonexistent"

		// Apply environment setup
		testmain.Setup()
		defer testmain.TearDown()

		envManager, err := envmanager.NewEnvManager(envmanager.WithEnvironment("nonexistent"))
		assert.NoError(t, err, "expected no error while creating env manager")

		err = envManager.LoadEnvironmentVariables()
		assert.NoError(t, err, "expected no error while loading environment variables for a non-existent environment as it's by default consider .env file")
	})
}

func TestEnvManager_ReadEnvironmentVariables(t *testing.T) {
	envManager, err := envmanager.NewEnvManager(envmanager.WithEnvironment("test"))
	assert.NoError(t, err, "expected no error while creating environment manager")

	envVars, err := envManager.ReadEnvironmentVariables()
	assert.NoError(t, err, "expected no error while reading environment variables")
	assert.NotNil(t, envVars, "expected non-nil map of environment variables")

	expectedKeys := []string{
		envmanager.DbHostKey,
		envmanager.DbPortKey,
		envmanager.DbDatabaseKey,
		envmanager.DbUsernameKey,
		envmanager.DbPasswordKey,
	}

	for _, key := range expectedKeys {
		value, ok := envVars[key]
		assert.True(t, ok, "expected key %s to be present in the environment variables", key)
		assert.NotEmpty(t, value, "expected key %s to have a non-empty value", key)
	}

	for key, expectedValue := range envVars {
		err := os.Setenv(key, expectedValue)
		assert.NoError(t, err, "expected no error while setting environment variable %s", key)

		actualValue := os.Getenv(key)
		assert.Equal(t, expectedValue, actualValue, "expected environment variable %s to match the value set by ReadEnvironmentVariables", key)
	}

	for key := range envVars {
		_ = os.Unsetenv(key)
	}
}

func TestEnvManager_GetWithDefaults(t *testing.T) {
	tests := []struct {
		name          string
		setupEnv      map[string]string
		useDefault    bool
		expectedValue string
	}{
		{
			name:          "test using default value when env variable is not provided or empty",
			setupEnv:      map[string]string{"APP_NAME": ""},
			useDefault:    true,
			expectedValue: "default_app",
		},
		{
			name:          "test using default value when use default is false",
			setupEnv:      map[string]string{"APP_NAME": ""},
			useDefault:    false,
			expectedValue: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTestEnvironment(tt.setupEnv)

			manager, _ := envmanager.NewEnvManager(envmanager.WithUseDefault(tt.useDefault))
			value := manager.Get(envmanager.AppNameKey)
			assert.Equal(t, tt.expectedValue, value)
		})
	}
}

func TestEnvManager_GetWithEnvironmentVariable(t *testing.T) {
	testmain.EnvVars["APP_ENV"] = "test-env"

	testmain.Setup()
	defer testmain.TearDown()

	manager, _ := envmanager.NewEnvManager()

	value := manager.Get("APP_ENV")
	expectedValue := "test-env"

	assert.Equal(t, expectedValue, value, "expected environment value to match")
}

func TestEnvManager_WithUseCache(t *testing.T) {
	t.Run("Verify that cached value remains consistent", func(t *testing.T) {
		testmain.EnvVars["CACHE_TEST_KEY"] = "initial-value"
		testmain.Setup()
		defer testmain.TearDown()

		manager, _ := envmanager.NewEnvManager(envmanager.WithUseCache(true))

		initialValue := manager.Get("CACHE_TEST_KEY")
		assert.Equal(t, "initial-value", initialValue, "expected initial cached value to be 'initial-value'")

		testmain.EnvVars["CACHE_TEST_KEY"] = "updated-value"
		testmain.Setup()

		cachedValue := manager.Get("CACHE_TEST_KEY")
		assert.Equal(t, "initial-value", cachedValue, "expected cached value to remain 'initial-value' even after environment variable update")
	})

	t.Run("Verify that clearing cache works correctly", func(t *testing.T) {
		testmain.EnvVars["CACHE_TEST_KEY"] = "initial-value"
		testmain.Setup()
		defer testmain.TearDown()

		manager, _ := envmanager.NewEnvManager(envmanager.WithUseCache(true))

		initialValue := manager.Get("CACHE_TEST_KEY")
		assert.Equal(t, "initial-value", initialValue, "expected initial cached value to be 'initial-value'")

		manager.ClearCache()

		testmain.EnvVars["CACHE_TEST_KEY"] = "new-value"
		testmain.Setup()

		newValue := manager.Get("CACHE_TEST_KEY")
		assert.Equal(t, "new-value", newValue, "expected new value after clearing cache")
	})
}

func TestEnvManager_WithExtendedDefaults(t *testing.T) {
	extendedDefaults := map[string]string{
		"EXTENDED_KEY_1": "extended_value_1",
		"EXTENDED_KEY_2": "extended_value_2",
	}

	// Define a DefaultExtender function that modifies the default map.
	extender := func(defaults map[string]string) {
		for key, value := range extendedDefaults {
			defaults[key] = value
		}
	}

	// Create an EnvManager using the DefaultExtender.
	envManager, err := envmanager.NewEnvManager(envmanager.WithExtendedDefaults(extender))
	assert.NoError(t, err, "expected no error while creating environment manager with extended defaults")

	t.Run("Verify that extended defaults are available", func(t *testing.T) {
		for key, expectedValue := range extendedDefaults {
			value := envManager.Get(key)
			assert.Equal(t, expectedValue, value, "expected value for key %s to match the extended default", key)
		}
	})

	t.Run("Verify that extended defaults are overridden by actual environment variables", func(t *testing.T) {
		testmain.EnvVars["EXTENDED_KEY_1"] = "overridden_value_1"

		testmain.Setup()
		defer testmain.TearDown()

		envManager, err = envmanager.NewEnvManager(envmanager.WithExtendedDefaults(extender))
		assert.NoError(t, err, "expected no error while creating environment manager with extended defaults and overridden variables")

		value := envManager.Get("EXTENDED_KEY_1")
		assert.Equal(t, "overridden_value_1", value, "expected overridden value for EXTENDED_KEY_1")
	})
}

func setupTestEnvironment(setupEnv map[string]string) {
	// Clear and set up the environment based on the provided map
	testmain.TearDown()
	for key, value := range setupEnv {
		testmain.EnvVars[key] = value
	}
	testmain.Setup()
}
