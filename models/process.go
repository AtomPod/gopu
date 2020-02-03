package models

//Process is processes information
type Process struct {
	Pid           int32   `json:"pid"`
	Name          string  `json:"name"`
	Background    bool    `json:"background"`
	CPUPercent    float64 `json:"cpuPercent"`
	Cmdline       string  `json:"cmdline"`
	CreateTime    int64   `json:"createTime"`
	Cwd           string  `json:"cwd"`
	IsRunning     bool    `json:"isRunning"`
	MemoryPercent float32 `json:"memoryPercent"`
	NumOfThreads  int32   `json:"numOfThreads"`
	NumOfFDs      int32   `json:"numOfFDs"`
}
