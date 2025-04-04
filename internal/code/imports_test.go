package code

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractModuleName(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	modDir := filepath.Join(wd, "..", "..", "testdata")
	content, err := os.ReadFile(filepath.Join(modDir, "gomod-with-leading-comments.mod"))
	require.NoError(t, err)
	assert.Equal(t, "github.com/99designs/gqlgen", extractModuleName(content))
}

func TestImportPathForDir(t *testing.T) {
	wd, err := os.Getwd()

	require.NoError(t, err)

	assert.Equal(t, "github.com/99designs/gqlgen/internal/code", ImportPathForDir(wd))

	assert.Equal(t, "github.com/99designs/gqlgen/internal/code", ImportPathForDir(wd))
	assert.Equal(t, "github.com/99designs/gqlgen/api", ImportPathForDir(filepath.Join(wd, "..", "..", "api")))

	// doesnt contain go code, but should still give a valid import path
	assert.Equal(t, "github.com/99designs/gqlgen/docs", ImportPathForDir(filepath.Join(wd, "..", "..", "docs")))

	// directory does not exist
	assert.Equal(t, "github.com/99designs/gqlgen/dos", ImportPathForDir(filepath.Join(wd, "..", "..", "dos")))

	// out of module
	assert.Empty(t, ImportPathForDir(filepath.Join(wd, "..", "..", "..")))

	if runtime.GOOS == "windows" {
		assert.Empty(t, ImportPathForDir("C:/doesnotexist"))
	} else {
		assert.Empty(t, ImportPathForDir("/doesnotexist"))
	}
}

func TestNameForDir(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	assert.Equal(t, "tmp", NameForDir("/tmp"))
	assert.Equal(t, "code", NameForDir(wd))
	assert.Equal(t, "docs", NameForDir(wd+"/../../docs"))
	assert.Equal(t, "main", NameForDir(wd+"/../.."))
}
