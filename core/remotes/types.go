package remotes

// Remote represents a remote foundation source
type Remote struct {
	Name       string `yaml:"name" json:"name"`
	Path       string `yaml:"path" json:"path"`
	PublicOnly bool   `yaml:"public_only" json:"public_only"`
}

// SyncResult contains the result of a sync operation
type SyncResult struct {
	Success       bool              `json:"success"`
	SyncedRemotes []string          `json:"synced_remotes"`
	Errors        map[string]string `json:"errors"` // remote name -> error message
	FilesCopied   int               `json:"files_copied"`
}

// RemoteInfo provides information about a synced remote
type RemoteInfo struct {
	Name         string   `json:"name"`
	Path         string   `json:"path"`
	FileCount    int      `json:"file_count"`
	LastModified string   `json:"last_modified"`
	Files        []string `json:"files"`
}
