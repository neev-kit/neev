package migration

// SourceType represents the type of project being migrated.
type SourceType string

const (
	SourceTypeOpenSpec SourceType = "openspec"
	SourceTypeSpecKit  SourceType = "speckit"
	SourceTypeAuto     SourceType = "auto"
)

// MigrationConfig holds migration parameters.
type MigrationConfig struct {
	RootDir    string
	SourceType SourceType
	DryRun     bool
	BackupOld  bool
}

// MigrationResult holds the results of a migration operation.
type MigrationResult struct {
	Success          bool
	SourceType       SourceType
	FilesMovedCount  int
	DirsCreatedCount int
	Messages         []string
	Errors           []string
	BackupDir        string // Path to backup if created
}
