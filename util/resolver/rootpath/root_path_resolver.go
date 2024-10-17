package rootpath

import (
	diabuddyErrors "github.com/hbttundar/diabuddy-errors"
	"os"
	"path/filepath"
)

type RootPathResolver struct{}

func NewRootPathResolver() *RootPathResolver {
	return &RootPathResolver{}
}

func (pr *RootPathResolver) Resolve(path string) (string, diabuddyErrors.ApiErrors) {
	basePath, err := filepath.Abs(path)
	if err != nil {
		return "", diabuddyErrors.NewApiError(diabuddyErrors.InternalServerErrorType, "could not find app root directory", diabuddyErrors.WithInternalError(err))
	}

	rootPath, err := pr.findRootDir(basePath)
	if err != nil {
		return "", diabuddyErrors.NewApiError(diabuddyErrors.InternalServerErrorType, "could not find app root directory", diabuddyErrors.WithInternalError(err))
	}
	return rootPath, nil
}

func (pr *RootPathResolver) findRootDir(dir string) (string, diabuddyErrors.ApiErrors) {
	if dir == "" {
		return "", diabuddyErrors.NewApiError(diabuddyErrors.InternalServerErrorType, "directory path is not set")
	}
	dir = filepath.Clean(dir)

	// Traverse upwards until we find `go.mod` or reach the system root directory.
	for {
		goModPath := filepath.Join(dir, "go.mod")
		file, err := os.Open(goModPath)
		if err == nil {
			defer file.Close()
			if fi, err := file.Stat(); err == nil && !fi.IsDir() {
				return dir, nil
			}
		}

		// Move to the parent directory.
		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			// Reached root directory, stop the search.
			break
		}
		dir = parentDir
	}

	return "", diabuddyErrors.NewApiError(diabuddyErrors.InternalServerErrorType, "could not find app root directory; go.mod not found")
}
