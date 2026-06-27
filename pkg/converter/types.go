package converter

type JobType int

const (
	JobConvert JobType = iota
	JobCopy
)

type Job struct {
	Source string
	Dest   string
	Type   JobType
}

type Result struct {
	Type    JobType
	Source  string
	Dest    string
	Success bool
	Error   error
}

type Config struct {
	WorkDir string
	Skip    bool
	Replace bool
}

type ScanResult struct {
	CBRFiles []string
	CBZFiles []string
}

type ProgressEvent struct {
	Status string
	File   string
	Dest   string
	Err    error
}
