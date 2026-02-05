package buildinfo

import "fmt"

// These values can be overridden at build time with:
// -ldflags "-X github.com/fixi2/InfraTrack/internal/buildinfo.Version=v0.2.0 -X github.com/fixi2/InfraTrack/internal/buildinfo.Commit=<sha>"
var (
	Version = "dev"
	Commit  = "none"
)

func String() string {
	return fmt.Sprintf("%s (%s)", Version, Commit)
}
