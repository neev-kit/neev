package remotes

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Sync synchronizes remote foundations to the local .neev/remotes directory
func Sync(rootDir string, remotes []Remote) (*SyncResult, error) {
	result := &SyncResult{
		Success:       true,
		SyncedRemotes: []string{},
		Errors:        make(map[string]string),
		FilesCopied:   0,
	}

	remotesDir := filepath.Join(rootDir, ".neev", "remotes")

	// Ensure remotes directory exists
	if err := os.MkdirAll(remotesDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create remotes directory: %w", err)
	}

	for _, remote := range remotes {
		if err := syncRemote(rootDir, remotesDir, remote, result); err != nil {
			result.Success = false
			result.Errors[remote.Name] = err.Error()
		} else {
			result.SyncedRemotes = append(result.SyncedRemotes, remote.Name)
		}
	}

	return result, nil
}

// syncRemote syncs a single remote foundation
func syncRemote(rootDir, remotesDir string, remote Remote, result *SyncResult) error {
	// Resolve the remote path (handle relative paths)
	remotePath := remote.Path
	if !filepath.IsAbs(remotePath) {
		remotePath = filepath.Join(rootDir, remotePath)
	}

	// Check if remote path exists
	if _, err := os.Stat(remotePath); os.IsNotExist(err) {
		return fmt.Errorf("remote path does not exist: %s", remotePath)
	}

	// Create destination directory for this remote
	destDir := filepath.Join(remotesDir, remote.Name)
	if err := os.RemoveAll(destDir); err != nil {
		return fmt.Errorf("failed to clean destination directory: %w", err)
	}
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Copy files from remote to destination
	filesCopied := 0
	err := filepath.Walk(remotePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			// Skip archive subdirectory
			if info.Name() == "archive" {
				return filepath.SkipDir
			}
			return nil
		}

		// Only copy markdown files for now
		if !strings.HasSuffix(info.Name(), ".md") {
			return nil
		}

		// If public_only is set, skip files that start with underscore
		if remote.PublicOnly && strings.HasPrefix(info.Name(), "_") {
			return nil
		}

		// Calculate relative path
		relPath, err := filepath.Rel(remotePath, path)
		if err != nil {
			return err
		}

		// Copy file
		destPath := filepath.Join(destDir, relPath)
		
		// Validate destPath is within destDir to prevent path traversal
		cleanDest := filepath.Clean(destPath)
		cleanDestDir := filepath.Clean(destDir)
		if !strings.HasPrefix(cleanDest, cleanDestDir+string(os.PathSeparator)) && cleanDest != cleanDestDir {
			return fmt.Errorf("invalid destination path: %s", destPath)
		}
		
		if err := copyFile(path, destPath); err != nil {
			return fmt.Errorf("failed to copy %s: %w", path, err)
		}

		filesCopied++
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to sync remote %s: %w", remote.Name, err)
	}

	result.FilesCopied += filesCopied
	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	// Ensure destination directory exists
	destDir := filepath.Dir(dst)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	// Open source file
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create destination file
	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Copy contents
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	// Sync to disk
	return destFile.Sync()
}

// GetRemoteInfo returns information about a synced remote
func GetRemoteInfo(rootDir, remoteName string) (*RemoteInfo, error) {
	// Validate remoteName to prevent path traversal
	if remoteName == "" {
		return nil, fmt.Errorf("remote name cannot be empty")
	}
	// Ensure remoteName is a simple name without path separators or traversal sequences
	if remoteName != filepath.Base(remoteName) ||
		strings.Contains(remoteName, "..") ||
		strings.ContainsAny(remoteName, `/\`) {
		return nil, fmt.Errorf("invalid remote name '%s'", remoteName)
	}
	
	remoteDir := filepath.Join(rootDir, ".neev", "remotes", remoteName)

	if _, err := os.Stat(remoteDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("remote '%s' not found (not synced?)", remoteName)
	}

	info := &RemoteInfo{
		Name:  remoteName,
		Path:  remoteDir,
		Files: []string{},
	}

	// Count files and gather names
	var lastModTime time.Time
	err := filepath.Walk(remoteDir, func(path string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !fileInfo.IsDir() {
			relPath, relErr := filepath.Rel(remoteDir, path)
			if relErr != nil {
				return relErr
			}
			info.Files = append(info.Files, relPath)
			info.FileCount++

			// Update last modified if this file is newer
			if lastModTime.IsZero() || fileInfo.ModTime().After(lastModTime) {
				lastModTime = fileInfo.ModTime()
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to scan remote directory: %w", err)
	}

	if !lastModTime.IsZero() {
		info.LastModified = lastModTime.Format("2006-01-02 15:04:05")
	}

	return info, nil
}

// ListRemotes lists all synced remotes
func ListRemotes(rootDir string) ([]string, error) {
	remotesDir := filepath.Join(rootDir, ".neev", "remotes")

	if _, err := os.Stat(remotesDir); os.IsNotExist(err) {
		return []string{}, nil // No remotes directory yet
	}

	entries, err := os.ReadDir(remotesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read remotes directory: %w", err)
	}

	remotes := []string{}
	for _, entry := range entries {
		if entry.IsDir() {
			remotes = append(remotes, entry.Name())
		}
	}

	return remotes, nil
}
