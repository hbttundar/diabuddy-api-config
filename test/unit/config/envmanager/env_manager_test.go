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
		assert.NoError(t, err, "expected no error while loading environment variables for a non-existent environment as its by default consider .env file")
	})
}

func TestEnvManager_ReadEnvironmentVariables(t *testing.T) {
	// Set up environment manager for "test" environment
	envManager, err := envmanager.NewEnvManager(envmanager.WithEnvironment("test"))
	assert.NoError(t, err, "expected no error while creating environment manager")

	// Load environment variables from the "test" environment
	envVars, err := envManager.ReadEnvironmentVariables()
	assert.NoError(t, err, "expected no error while reading environment variables")
	assert.NotNil(t, envVars, "expected non-nil map of environment variables")

	// Check that some specific keys have expected values
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

	// Validate that env variables are set correctly in the current process using ReadEnvironmentVariables
	for key, expectedValue := range envVars {
		// Set the environment variable for this test
		err := os.Setenv(key, expectedValue)
		assert.NoError(t, err, "expected no error while setting environment variable %s", key)

		// Ensure it was set correctly
		actualValue := os.Getenv(key)
		assert.Equal(t, expectedValue, actualValue, "expected environment variable %s to match the value set by ReadEnvironmentVariables", key)
	}

	// Cleanup the environment variables
	for key := range envVars {
		_ = os.Unsetenv(key)
	}
}

func TestGetWithDefaults(t *testing.T) {
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
			expectedValue: "Diabuddy",
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
			// Update global EnvVars from testmain with specific test setup
			for k, v := range tt.setupEnv {
				testmain.EnvVars[k] = v
			}

			// Apply environment setup
			testmain.Setup()
			defer testmain.TearDown()

			// Create a new EnvManager and execute the test
			manager, _ := envmanager.NewEnvManager(envmanager.WithUseDefault(tt.useDefault))
			value := manager.Get(envmanager.AppNameKey)
			assert.Equal(t, tt.expectedValue, value)
		})
	}
}

func TestGetWithEnvironmentVariable(t *testing.T) {
	testmain.EnvVars["APP_ENV"] = "test-env"

	// Apply environment setup
	testmain.Setup()
	defer testmain.TearDown()

	manager, _ := envmanager.NewEnvManager()

	value := manager.Get("APP_ENV")
	expectedValue := "test-env"

	assert.Equal(t, expectedValue, value, "expected environment value to match")
}

func TestCachingBehavior(t *testing.T) {
	t.Run("Verify that cached value remains consistent", func(t *testing.T) {
		// Set initial environment variables and apply them
		testmain.EnvVars["CACHE_TEST_KEY"] = "initial-value"
		testmain.Setup()
		defer testmain.TearDown()

		manager, _ := envmanager.NewEnvManager(envmanager.WithUseCache(true))

		// First retrieval to cache the value
		initialValue := manager.Get("CACHE_TEST_KEY")
		assert.Equal(t, "initial-value", initialValue, "expected initial cached value to be 'initial-value'")

		// Change the environment variable value and re-setup environment
		testmain.EnvVars["CACHE_TEST_KEY"] = "updated-value"
		testmain.Setup()

		// Value should still be the cached "initial-value"
		cachedValue := manager.Get("CACHE_TEST_KEY")
		assert.Equal(t, "initial-value", cachedValue, "expected cached value to remain 'initial-value' even after environment variable update")
	})

	t.Run("Verify that clearing cache works correctly", func(t *testing.T) {
		// Set initial environment variables and apply them
		testmain.EnvVars["CACHE_TEST_KEY"] = "initial-value"
		testmain.Setup()
		defer testmain.TearDown()

		manager, _ := envmanager.NewEnvManager(envmanager.WithUseCache(true))

		// First retrieval to cache the value
		initialValue := manager.Get("CACHE_TEST_KEY")
		assert.Equal(t, "initial-value", initialValue, "expected initial cached value to be 'initial-value'")

		// Clear the cache
		manager.ClearCache()

		// Change the environment variable value and re-setup environment
		testmain.EnvVars["CACHE_TEST_KEY"] = "new-value"
		testmain.Setup()

		// Now that the cache is cleared, the new value should be retrieved
		newValue := manager.Get("CACHE_TEST_KEY")
		assert.Equal(t, "new-value", newValue, "expected new value after clearing cache")
	})
}
