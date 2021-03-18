package version

// Dynamically overridden at build time. See `ldflags` in .goreleaser.yml.
var (
	Version  = "dev"
	Revision = "dev"
)
