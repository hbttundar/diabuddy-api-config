package rootpath_test

import (
	"github.com/hbttundar/diabuddy-api-config/util/resolver/rootpath"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestRootPathResolver(t *testing.T) {
	// Mock different scenarios for root path resolution
	t.Run("Resolve root path successfully with go.mod", func(t *testing.T) {
		// Setup: Ensure a go.mod file is present in an ancestor directory for testing
		path, _ := filepath.Abs("./..")
		resolver := rootpath.NewRootPathResolver()
		rootPath, err := resolver.Resolve(path)

		assert.NoError(t, err, "expected no error while resolving root path")
		assert.NotEmpty(t, rootPath, "expected a valid root path to be resolved")

		// Validate that the resolved path contains a `go.mod` file
		goModPath := filepath.Join(rootPath, "go.mod")
		if _, err := os.Stat(goModPath); os.IsNotExist(err) {
			t.Fatalf("expected go.mod to exist in the resolved root path, but it was not found")
		}
	})

	t.Run("Fail to resolve root path when go.mod is missing", func(t *testing.T) {
		// Simulate an environment without a go.mod file.
		path, _ := filepath.Abs("./..")
		resolver := rootpath.NewRootPathResolver()
		rootPath, err := resolver.Resolve(path)
		// Temporarily rename or remove go.mod if possible for this test
		originalGoModPath := filepath.Join(rootPath, "./go.mod")
		temporaryGoModPath := filepath.Join(rootPath, "./go.mod.temp")

		if _, err := os.Stat(originalGoModPath); err == nil {
			_ = os.Rename(originalGoModPath, temporaryGoModPath)   // Temporarily move go.mod
			defer os.Rename(temporaryGoModPath, originalGoModPath) // Restore after test
		}

		rootPath, err = resolver.Resolve(path)

		assert.Error(t, err, "expected an error when go.mod is missing")
		assert.Empty(t, rootPath, "expected no root path to be resolved when go.mod is missing")
	})

	t.Run("Resolve root path from deeply nested directory", func(t *testing.T) {
		// Navigate to a nested directory to simulate a deep structure
		currentDir, err := os.Getwd()
		assert.NoError(t, err, "expected no error while getting current directory")

		// Create nested directory structure
		nestedDir := filepath.Join(currentDir, "test/nested/structure")
		err = os.MkdirAll(nestedDir, 0755)
		assert.NoError(t, err, "expected no error while creating nested directory structure")

		defer os.RemoveAll(filepath.Join(currentDir, "test")) // Cleanup nested structure

		// Change directory to nested structure
		err = os.Chdir(nestedDir)
		assert.NoError(t, err, "expected no error while changing to nested directory")

		defer os.Chdir(currentDir) // Change back to original directory after test
		path, _ := filepath.Abs("./..")
		resolver := rootpath.NewRootPathResolver()
		rootPath, err := resolver.Resolve(path)

		assert.NoError(t, err, "expected no error while resolving root path from nested directory")
		assert.NotEmpty(t, rootPath, "expected a valid root path to be resolved from nested directory")
	})
}
