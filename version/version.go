package version

import (
	"runtime/debug"
	"time"
)

var (
	Version    string = "0.0.0"
	Revision   string = "development"
	LastCommit time.Time
	DirtyBuild bool = true
)

func init() {
	// attempt to read vcs information if available
	info, ok := debug.ReadBuildInfo()
	if ok {
		for _, kv := range info.Settings {
			switch kv.Key {
			case "vcs.revision":
				SetRevision(kv.Value)
			case "vcs.time":
				lc, err := time.Parse(time.RFC3339, kv.Value)
				if err == nil && !lc.IsZero() {
					SetLastCommit(lc)
				}
			case "vcs.modified":
				SetDirtyBuild(kv.Value == "true")
			}
		}
	}
}

func SetVersion(version string) {
	Version = version
}

func SetRevision(revision string) {
	Revision = revision
}

func SetLastCommit(lastCommit time.Time) {
	LastCommit = lastCommit
}

func SetDirtyBuild(dirtyBuild bool) {
	DirtyBuild = dirtyBuild
}
