package model

type Config struct {
	Tag           string
	LocalFile     bool
	NoCache       bool
	Log           bool
	LocalFilePath string
	Dev           bool
	Repo          string
	Volume        string
	JumpBuild     bool
}
