package version

var (
	Version      string = "unknown"
	VersionMajor string = "unknown"
	VersionMinor string = "unknown"
	VersionPatch string = "unknown"
	SHA          string = "unknown"
	BuildDate    string = "unknown"

	FullVersionInfoAsString string = "unknown"
	VersionInfoAsString     string = "unknown"
	Semantic                string = "unknown"
)

func init() {
	Semantic = VersionMajor + "." + VersionMinor + "." + VersionPatch
	VersionInfoAsString = "crypto-rw/" + Semantic
	FullVersionInfoAsString = VersionInfoAsString + "/" + SHA + "/" + BuildDate

}
