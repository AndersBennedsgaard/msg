package version

var (
	// Version is set at build time using -ldflags.
	Version = "dev"
	// Commit is the Git commit SHA, also set at build time.
	Commit = "unknown"
	// Date is the build date, optional but useful.
	Date = "unknown"
)
