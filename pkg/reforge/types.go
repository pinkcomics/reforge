package reforge

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

type Settings struct {
	Replace bool `json:"replace"`
}

type Summary struct {
	Total       int    `json:"total"`
	Converted   int    `json:"converted"`
	Kept        int    `json:"kept"`
	Failed      int    `json:"failed"`
	Cancelled   bool   `json:"cancelled"`
	ElapsedMs   int64  `json:"elapsedMs"`
	ElapsedText string `json:"elapsedText"`
}
