package types

type DirectoryConfig struct {
	Name           string
	Paths          []string
	IgnorePrefixes []string
	Commands       []string
}

type Config struct {
	Delay int32
	Watch []DirectoryConfig
}
