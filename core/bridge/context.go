package bridge

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// BuildContext aggregates context from foundation and blueprints.
func BuildContext(focus string) (string, error) {
	var contextBuilder strings.Builder
	contextBuilder.WriteString("# Project Foundation\n")

	// Read foundation files
	foundationPath := ".neev/foundation"
	if err := readFilesInDir(foundationPath, &contextBuilder, focus); err != nil {
		return "", err
	}

	// Read blueprint files
	blueprintsPath := ".neev/blueprints"
	files, err := ioutil.ReadDir(blueprintsPath)
	if err != nil {
		return "", fmt.Errorf("failed to read blueprints directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			blueprintDir := filepath.Join(blueprintsPath, file.Name())
			if err := readFilesInDir(blueprintDir, &contextBuilder, focus); err != nil {
				return "", err
			}
		}
	}

	return contextBuilder.String(), nil
}

// BuildRemoteContext aggregates context from synced remote foundations
func BuildRemoteContext() (string, error) {
	remotesPath := ".neev/remotes"
	
	// Check if remotes directory exists
	if _, err := os.Stat(remotesPath); os.IsNotExist(err) {
		return "", nil // No remotes synced
	}

	var contextBuilder strings.Builder
	contextBuilder.WriteString("# Remote Foundations\n\n")

	// Read each remote directory
	remotes, err := ioutil.ReadDir(remotesPath)
	if err != nil {
		return "", fmt.Errorf("failed to read remotes directory: %w", err)
	}

	if len(remotes) == 0 {
		return "", nil
	}

	for _, remote := range remotes {
		if !remote.IsDir() {
			continue
		}

		remoteName := remote.Name()
		contextBuilder.WriteString(fmt.Sprintf("## Remote: %s\n\n", remoteName))

		remoteDir := filepath.Join(remotesPath, remoteName)
		if err := readFilesInDir(remoteDir, &contextBuilder, ""); err != nil {
			return "", fmt.Errorf("failed to read remote %s: %w", remoteName, err)
		}

		contextBuilder.WriteString("\n")
	}

	return contextBuilder.String(), nil
}

func readFilesInDir(dir string, builder *strings.Builder, focus string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", dir, err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".md" {
			content, err := ioutil.ReadFile(filepath.Join(dir, file.Name()))
			if err != nil {
				return fmt.Errorf("failed to read file %s: %w", file.Name(), err)
			}
			if focus == "" || strings.Contains(string(content), focus) {
				builder.WriteString(fmt.Sprintf("## File: %s\n%s\n", file.Name(), content))
			}
		}
	}

	return nil
}
