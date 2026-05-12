package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestStoreCommand(t *testing.T) {
	tempDir := setupTestEnvironment(t)

	filename := "teste1.txt"
	f, filePath := createTestFile(t, filename, tempDir)
	defer f.Close()

	content := "1"
	size := 9000
	storeContentTestFile(t, f, content, size)

	_, _, err := executeCommandC(rootCmd, "store", filePath)
	assert.NoError(t, err)
}

func TestGetCommand(t *testing.T) {
	tempDir := setupTestEnvironment(t)

	filename := "teste1.txt"
	f, filePath := createTestFile(t, filename, tempDir)
	defer f.Close()

	content := "1"
	size := 9000
	storeContentTestFile(t, f, content, size)

	_, _, err := executeCommandC(rootCmd, "store", filePath)
	assert.NoError(t, err)

	destDir := filepath.Join(tempDir, "getDir")
	_, _, err = executeCommandC(rootCmd, "get", filename, destDir)
	assert.NoError(t, err)
}

func setupTestEnvironment(t *testing.T) string {
	t.Helper()
	tempDir := t.TempDir()
	volumePath = filepath.Join(tempDir, "volume.dat")
	dbDirPath = filepath.Join(tempDir, "badger-data")
	return tempDir
}

func createTestFile(t *testing.T, fileName, tempDir string) (*os.File, string) {
	t.Helper()
	filePath := filepath.Join(tempDir, fileName)
	f, err := os.Create(filePath)
	assert.NoError(t, err)
	assert.NotNil(t, f)
	return f, filePath
}

func storeContentTestFile(t *testing.T, f *os.File, content string, size int) {
	t.Helper()
	byteContent := bytes.Repeat([]byte(content), size)
	n, err := f.Write(byteContent)
	assert.NoError(t, err)
	assert.Equal(t, size, n)
}

func executeCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err
}
