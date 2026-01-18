package bridge

import (
	"fmt"
	"io/ioutil"
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
























































}	return nil	}		}			}				builder.WriteString(fmt.Sprintf("## File: %s\n%s\n", file.Name(), content))			if focus == "" || strings.Contains(string(content), focus) {			}				return fmt.Errorf("failed to read file %s: %w", file.Name(), err)			if err != nil {			content, err := ioutil.ReadFile(filepath.Join(dir, file.Name()))		if filepath.Ext(file.Name()) == ".md" {	for _, file := range files {	}		return fmt.Errorf("failed to read directory %s: %w", dir, err)	if err != nil {	files, err := ioutil.ReadDir(dir)func readFilesInDir(dir string, builder *strings.Builder, focus string) error {}	return contextBuilder.String(), nil	}		}			}				return "", err			if err := readFilesInDir(blueprintDir, &contextBuilder, focus); err != nil {			blueprintDir := filepath.Join(blueprintsPath, file.Name())		if file.IsDir() {	for _, file := range files {	}		return "", fmt.Errorf("failed to read blueprints directory: %w", err)	if err != nil {	files, err := ioutil.ReadDir(blueprintsPath)	blueprintsPath := ".neev/blueprints"	// Read blueprint files	}		return "", err	if err := readFilesInDir(foundationPath, &contextBuilder, focus); err != nil {	foundationPath := ".neev/foundation"	// Read foundation files	contextBuilder.WriteString("# Project Foundation\n")	var contextBuilder strings.Builderfunc BuildContext(focus string) (string, error) {// BuildContext aggregates context from foundation and blueprints.

// BuildContext aggregates context from foundation and blueprints.)	"strings"	"path/filepath"	"io/ioutil"	"fmt"import (