package config

import (
	diabuddyErrors "github.com/hbttundar/diabuddy-errors"
)

type Config interface {
	Get(key string, defaultValue ...string) string
	Validate() diabuddyErrors.ApiErrors
}
