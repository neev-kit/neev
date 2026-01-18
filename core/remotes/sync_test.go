package remotes

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSync_NoRemotes(t *testing.T) {
	tmpDir := t.TempDir()

	result, err := Sync(tmpDir, []Remote{})
	if err != nil {
		t.Fatalf("Sync failed: %v", err)
	}

	if !result.Success {
		t.Errorf("Expected success with no remotes")
	}

	if len(result.SyncedRemotes) != 0 {
		t.Errorf("Expected 0 synced remotes, got %d", len(result.SyncedRemotes))
	}
}

func TestSync_SingleRemote(t *testing.T) {
	// Create temporary directories
	tmpDir := t.TempDir()
	remoteDir := filepath.Join(tmpDir, "remote-foundation")
	os.MkdirAll(remoteDir, 0755)

	// Create some markdown files in the remote
	os.WriteFile(filepath.Join(remoteDir, "auth.md"), []byte("# Auth"), 0644)
	os.WriteFile(filepath.Join(remoteDir, "api.md"), []byte("# API"), 0644)

	// Create remote configuration
	remote := Remote{
		Name:       "test-remote",
		Path:       remoteDir,
		PublicOnly: false,
	}

	// Sync
	result, err := Sync(tmpDir, []Remote{remote})
	if err != nil {
		t.Fatalf("Sync failed: %v", err)
	}

	if !result.Success {
		t.Errorf("Expected success")
	}

	if len(result.SyncedRemotes) != 1 {
		t.Errorf("Expected 1 synced remote, got %d", len(result.SyncedRemotes))
	}

	if result.FilesCopied != 2 {
		t.Errorf("Expected 2 files copied, got %d", result.FilesCopied)
	}

	// Verify files were copied
	destDir := filepath.Join(tmpDir, ".neev", "remotes", "test-remote")
	if _, err := os.Stat(filepath.Join(destDir, "auth.md")); os.IsNotExist(err) {
		t.Errorf("auth.md was not copied")
	}
	if _, err := os.Stat(filepath.Join(destDir, "api.md")); os.IsNotExist(err) {
		t.Errorf("api.md was not copied")
	}
}

func TestSync_PublicOnly(t *testing.T) {
	// Create temporary directories
	tmpDir := t.TempDir()
	remoteDir := filepath.Join(tmpDir, "remote-foundation")
	os.MkdirAll(remoteDir, 0755)

	// Create public and private files
	os.WriteFile(filepath.Join(remoteDir, "public.md"), []byte("# Public"), 0644)
	os.WriteFile(filepath.Join(remoteDir, "_private.md"), []byte("# Private"), 0644)

	// Create remote with public_only flag
	remote := Remote{
		Name:       "test-remote",
		Path:       remoteDir,
		PublicOnly: true,
	}

	// Sync
	result, err := Sync(tmpDir, []Remote{remote})
	if err != nil {
		t.Fatalf("Sync failed: %v", err)
	}

	if result.FilesCopied != 1 {
		t.Errorf("Expected 1 file copied (public only), got %d", result.FilesCopied)
	}

	// Verify only public file was copied
	destDir := filepath.Join(tmpDir, ".neev", "remotes", "test-remote")
	if _, err := os.Stat(filepath.Join(destDir, "public.md")); os.IsNotExist(err) {
		t.Errorf("public.md was not copied")
	}
	if _, err := os.Stat(filepath.Join(destDir, "_private.md")); err == nil {
		t.Errorf("_private.md should not have been copied")
	}
}

func TestSync_NonExistentRemote(t *testing.T) {
	tmpDir := t.TempDir()

	remote := Remote{
		Name:       "missing",
		Path:       "/nonexistent/path",
		PublicOnly: false,
	}

	result, err := Sync(tmpDir, []Remote{remote})
	if err != nil {
		t.Fatalf("Sync failed: %v", err)
	}

	if result.Success {
		t.Errorf("Expected failure for non-existent remote")
	}

	if len(result.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(result.Errors))
	}

	if _, exists := result.Errors["missing"]; !exists {
		t.Errorf("Expected error for 'missing' remote")
	}
}

func TestSync_ArchiveSkipped(t *testing.T) {
	// Create temporary directories
	tmpDir := t.TempDir()
	remoteDir := filepath.Join(tmpDir, "remote-foundation")
	os.MkdirAll(remoteDir, 0755)
	archiveDir := filepath.Join(remoteDir, "archive")
	os.MkdirAll(archiveDir, 0755)

	// Create files in main and archive
	os.WriteFile(filepath.Join(remoteDir, "main.md"), []byte("# Main"), 0644)
	os.WriteFile(filepath.Join(archiveDir, "archived.md"), []byte("# Archived"), 0644)

	remote := Remote{
		Name:       "test-remote",
		Path:       remoteDir,
		PublicOnly: false,
	}

	result, err := Sync(tmpDir, []Remote{remote})
	if err != nil {
		t.Fatalf("Sync failed: %v", err)
	}

	if result.FilesCopied != 1 {
		t.Errorf("Expected 1 file (archive should be skipped), got %d", result.FilesCopied)
	}

	// Verify archive was not copied
	destArchive := filepath.Join(tmpDir, ".neev", "remotes", "test-remote", "archive")
	if _, err := os.Stat(destArchive); err == nil {
		t.Errorf("archive directory should not have been copied")
	}
}

func TestGetRemoteInfo(t *testing.T) {
	// Create temporary directory with synced remote
	tmpDir := t.TempDir()
	remoteDir := filepath.Join(tmpDir, ".neev", "remotes", "test-remote")
	os.MkdirAll(remoteDir, 0755)

	// Create some files
	os.WriteFile(filepath.Join(remoteDir, "file1.md"), []byte("Content 1"), 0644)
	os.WriteFile(filepath.Join(remoteDir, "file2.md"), []byte("Content 2"), 0644)

	info, err := GetRemoteInfo(tmpDir, "test-remote")
	if err != nil {
		t.Fatalf("GetRemoteInfo failed: %v", err)
	}

	if info.Name != "test-remote" {
		t.Errorf("Expected name 'test-remote', got '%s'", info.Name)
	}

	if info.FileCount != 2 {
		t.Errorf("Expected 2 files, got %d", info.FileCount)
	}

	if len(info.Files) != 2 {
		t.Errorf("Expected 2 file names, got %d", len(info.Files))
	}
}

func TestGetRemoteInfo_NotFound(t *testing.T) {
	tmpDir := t.TempDir()

	_, err := GetRemoteInfo(tmpDir, "nonexistent")
	if err == nil {
		t.Errorf("Expected error for non-existent remote")
	}
}

func TestListRemotes(t *testing.T) {
	tmpDir := t.TempDir()
	remotesDir := filepath.Join(tmpDir, ".neev", "remotes")
	os.MkdirAll(filepath.Join(remotesDir, "remote1"), 0755)
	os.MkdirAll(filepath.Join(remotesDir, "remote2"), 0755)

	remotes, err := ListRemotes(tmpDir)
	if err != nil {
		t.Fatalf("ListRemotes failed: %v", err)
	}

	if len(remotes) != 2 {
		t.Errorf("Expected 2 remotes, got %d", len(remotes))
	}
}

func TestListRemotes_NoRemotesDir(t *testing.T) {
	tmpDir := t.TempDir()

	remotes, err := ListRemotes(tmpDir)
	if err != nil {
		t.Fatalf("ListRemotes failed: %v", err)
	}

	if len(remotes) != 0 {
		t.Errorf("Expected 0 remotes, got %d", len(remotes))
	}
}
