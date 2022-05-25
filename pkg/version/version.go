package version

import (
	"encoding/json"
	"fmt"
	"runtime"
)

var (
	gitVersion   = "v0.0.0-master+$Format:%H$"
	gitCommit    = "$Format:%H$" // sha1 from git, output of $(git rev-parse HEAD)
	gitTreeState = ""            // state of git tree, either "clean" or "dirty"
	gitBranch    = ""

	buildDate = "1970-01-01T00:00:00Z" // build date in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')

	version = &Version{
		Version:      gitVersion,
		GitCommit:    gitCommit,
		GitTreeState: gitTreeState,
		GitBranch:    gitBranch,
		BuildDate:    buildDate,
		GoVersion:    runtime.Version(),
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
)

type Version struct {
	Version      string `json:"version"`
	GitCommit    string `json:"gitCommit"`
	GitTreeState string `json:"gitTreeState"`
	GitBranch    string `json:"gitBranch"`
	BuildDate    string `json:"buildDate"`
	GoVersion    string `json:"goVersion"`
	Platform     string `json:"platform"`
}

func (v *Version) String() string {
	data, _ := json.Marshal(v)
	return string(data)
}

func Print() {
	data, _ := json.MarshalIndent(version, "", "    ")
	fmt.Println(string(data))
}

func Get() *Version {
	return version
}
