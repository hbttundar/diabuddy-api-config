package rootpath_test

import (
	testmain "github.com/hbttundar/diabuddy-api-config/test"
	"github.com/hbttundar/diabuddy-api-config/util/resolver/rootpath"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestRootPathResolver_Resolve(t *testing.T) {

	tests := []struct {
		name           string
		path           string
		SetupGoMod     bool
		expectedError  bool
		useNestedSetup bool
	}{
		{
			name:           "Resolve root path successfully with go.mod",
			path:           "./..",
			SetupGoMod:     false,
			expectedError:  false,
			useNestedSetup: false,
		},
		{
			name:           "Fail to resolve root path when go.mod is missing",
			path:           "./..",
			SetupGoMod:     true,
			expectedError:  false,
			useNestedSetup: false,
		},
		{
			name:           "Fail to resolve root path when the current director doesn't exist",
			path:           "./none_exist_dir",
			SetupGoMod:     true,
			expectedError:  false,
			useNestedSetup: false,
		},
		{
			name:           "Resolve root path from deeply nested directory",
			path:           "./..",
			SetupGoMod:     true,
			expectedError:  false,
			useNestedSetup: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.useNestedSetup {
				setupNestedDir(t)
			}
			if tt.SetupGoMod {
				testmain.SetupGoMod()
				defer testmain.TeardownGoMod()
			}
			path, _ := filepath.Abs(tt.path)
			resolver := rootpath.NewRootPathResolver()
			rootPath, err := resolver.Resolve(path)
			if !tt.expectedError && err == nil {
				assert.NoError(t, err, "expected no error while resolving root path")
				assert.NotEmpty(t, rootPath, "expected a valid root path to be resolved")
				goModPath := filepath.Join(rootPath, "go.mod")
				if _, err := os.Stat(goModPath); os.IsNotExist(err) {
					t.Fatalf("expected go.mod to exist in the resolved root path, but it was not found")
				}
			} else {
				assert.Error(t, err, "expected an error when go.mod is missing")
				assert.Empty(t, rootPath, "expected no root path to be resolved when go.mod is missing")
			}
		})
	}

}

func setupNestedDir(t *testing.T) {
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
}
