package models

//Memory memory information
type Memory struct {
	Total       uint64  `json:"total"`
	Available   uint64  `json:"available"`
	Used        uint64  `json:"used"`
	Freed       uint64  `json:"freed"`
	UsedPercent float64 `json:"usedPercent"`
}
