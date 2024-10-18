package envmanager

import (
	"fmt"
	"github.com/hbttundar/diabuddy-api-config/util/resolver/rootpath"
	diabuddyErrors "github.com/hbttundar/diabuddy-errors"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"sync"
)

const (
	AppNameKey           = "APP_NAME"
	AppEnvKey            = "APP_ENV"
	AppEncryptionKey     = "APP_KEY"
	AppDebugKey          = "APP_DEBUG"
	AppUrlKey            = "APP_URL"
	AppTimezoneKey       = "APP_TIMEZONE"
	AppLocaleKey         = "APP_LOCALE"
	AppFallbackLocaleKey = "APP_FALLBACK_LOCALE"
	AppCipherKey         = "APP_CIPHER"
	AuthSecretKey        = "AUTH_SECRET"
	DbUrlKey             = "DATABASE_URL"
	DbHostKey            = "DB_HOST"
	DbPortKey            = "DB_PORT"
	DbDatabaseKey        = "DB_DATABASE"
	DbUsernameKey        = "DB_USERNAME"
	DbPasswordKey        = "DB_PASSWORD"
	DbSslModeKey         = "SSL_MODE"
)

type EnvManager struct {
	useDefaults  bool
	useCache     bool
	environment  string
	cache        sync.Map
	defaults     map[string]string
	pathResolver *rootpath.RootPathResolver
}

type EnvOption func(*EnvManager) diabuddyErrors.ApiErrors

type DefaultExtender func(map[string]string)

// WithEnvironment sets the environment for loading the corresponding .env resolver
func WithEnvironment(environment string) EnvOption {
	return func(em *EnvManager) diabuddyErrors.ApiErrors {
		em.environment = environment
		return nil
	}
}

// WithUseDefault sets the useDefaults flag for the EnvManager
func WithUseDefault(useDefaults bool) EnvOption {
	return func(em *EnvManager) diabuddyErrors.ApiErrors {
		em.useDefaults = useDefaults
		return nil
	}
}

// WithUseCache sets the useCache flag for the EnvManager
func WithUseCache(useCache bool) EnvOption {
	return func(em *EnvManager) diabuddyErrors.ApiErrors {
		em.useCache = useCache
		return nil
	}
}

// WithExtendedDefaults allows extending the default values during initialization
func WithExtendedDefaults(extender DefaultExtender) EnvOption {
	return func(em *EnvManager) diabuddyErrors.ApiErrors {
		extender(em.defaults)
		return nil
	}
}

// NewEnvManager creates an EnvManager with the specified options
func NewEnvManager(options ...EnvOption) (*EnvManager, diabuddyErrors.ApiErrors) {
	em := &EnvManager{
		useDefaults:  true,
		useCache:     false,
		environment:  "production",
		defaults:     defaultValues(),
		pathResolver: rootpath.NewRootPathResolver(),
	}

	// Apply the provided options
	for _, option := range options {
		if err := option(em); err != nil {
			return nil, err
		}
	}

	// Load environment variables from file
	err := em.LoadEnvironmentVariables()
	if err != nil {
		return nil, err
	}

	return em, nil
}

// LoadEnvironmentVariables loads environment variables from the appropriate .env file based on the environment
func (em *EnvManager) LoadEnvironmentVariables() diabuddyErrors.ApiErrors {
	envFilepath, apiError := em.getEnvFilePath()
	if apiError != nil {
		return apiError
	}

	err := godotenv.Load(envFilepath)
	if err != nil {
		return diabuddyErrors.NewApiError(diabuddyErrors.InternalServerErrorType, fmt.Sprintf("failed to load environment variables from: %s file.", envFilepath), diabuddyErrors.WithInternalError(err))
	}

	return nil
}

func (em *EnvManager) getEnvFilePath() (string, diabuddyErrors.ApiErrors) {
	basePath, err := filepath.Abs("./")
	if err != nil {
		return "", diabuddyErrors.NewApiError(diabuddyErrors.InternalServerErrorType, "could not find appconfig root directory", diabuddyErrors.WithInternalError(err))
	}
	envDir, resolveError := em.pathResolver.Resolve(basePath)
	if resolveError != nil {
		return "", resolveError
	}

	envFileName := ".env"
	if em.environment == "test" {
		envFileName = ".env.test"
	}
	return filepath.Join(envDir, envFileName), nil
}

func (em *EnvManager) ReadEnvironmentVariables() (map[string]string, diabuddyErrors.ApiErrors) {
	envFilepath, apiError := em.getEnvFilePath()
	if apiError != nil {
		return nil, apiError
	}

	envMaps, err := godotenv.Read(envFilepath)
	if err != nil {
		return nil, diabuddyErrors.NewApiError(diabuddyErrors.InternalServerErrorType, fmt.Sprintf("failed to read environment variablesfron : %s file.", envFilepath), diabuddyErrors.WithInternalError(err))
	}

	return envMaps, nil

}

// Get retrieves an environment variable value. If it's not set, it will use the default value if enabled.
func (em *EnvManager) Get(key string, defaultValue ...string) string {
	// First, attempt to retrieve from cache
	if val, ok := em.getFromCache(key); ok {
		return val
	}

	// Retrieve from environment
	val := os.Getenv(key)
	if val == "" && len(defaultValue) > 0 {
		val = defaultValue[0]
	}

	// Retrieve from defaults if necessary
	if val == "" && em.useDefaults {
		if defVal, ok := em.defaults[key]; ok {
			val = defVal
		}
	}

	// Store in cache for future reference
	em.storeInCache(key, val)
	return val
}

// Defaults provides default values for environment variables.
func defaultValues() map[string]string {
	return map[string]string{
		AppNameKey:           "default_app",
		AppEnvKey:            "local",
		AppEncryptionKey:     "",
		AppDebugKey:          "false",
		AppUrlKey:            "http://localhost",
		AppTimezoneKey:       "UTC",
		AppLocaleKey:         "en",
		AppFallbackLocaleKey: "en",
		AppCipherKey:         "AES-256-CBC",
		AuthSecretKey:        "my_default_secret",
		DbHostKey:            "127.0.0.1",
		DbPortKey:            "5432",
		DbDatabaseKey:        "default_db",
		DbUsernameKey:        "default_user",
		DbPasswordKey:        "default_pass",
		DbSslModeKey:         "disable",
	}
}

// Private method to get value from cache
func (em *EnvManager) getFromCache(key string) (string, bool) {
	if em.useCache {
		if val, ok := em.cache.Load(key); ok {
			return val.(string), true
		}
	}
	return "", false
}

// Private method to store value in cache
func (em *EnvManager) storeInCache(key, value string) {
	if em.useCache {
		em.cache.Store(key, value)
	}
}

// ClearCache clears all the cached environment variables.
func (em *EnvManager) ClearCache() {
	if em.useCache {
		em.cache = sync.Map{}
	}
}
