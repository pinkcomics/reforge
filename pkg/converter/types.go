package converter

// JobType indica o tipo de operação a realizar num arquivo.
type JobType int

const (
	JobConvert JobType = iota // CBR → CBZ
	JobCopy                   // CBZ → cópia direta
)

// Job representa uma unidade de trabalho enviada ao worker pool.
type Job struct {
	Source string
	Dest   string
	Type   JobType
}

// Result é retornado pelo worker após processar um Job.
type Result struct {
	Type    JobType
	Source  string
	Dest    string
	Success bool
	Error   error
}

// Config agrupa as opções de execução da conversão.
type Config struct {
	WorkDir string // Pasta com os arquivos a processar
	Skip    bool   // Se true, pula preview e confirmação
	Replace bool   // Se true, remove CBRs originais após converter
}

// ScanResult resume o que foi encontrado na pasta de trabalho.
type ScanResult struct {
	CBRFiles []string // Caminhos absolutos dos .cbr encontrados
	CBZFiles []string // Caminhos absolutos dos .cbz encontrados
}

// ProgressEvent é emitido em tempo real durante a conversão.
// A GUI consome esses eventos via canal ou WebSocket.
type ProgressEvent struct {
	// Status pode ser: "converting", "done", "error", "kept"
	Status string
	File   string
	Dest   string // Preenchido apenas quando Status == "done"
	Err    error  // Preenchido apenas quando Status == "error"
}
