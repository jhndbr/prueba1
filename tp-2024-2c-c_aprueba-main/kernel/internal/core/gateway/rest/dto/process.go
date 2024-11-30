package dto

// + agregue tamaño
type Process struct {
	PID      uint32 `json:"pid"`
	FilePath string `json:"file_path"`
	Size     uint32 `json:"size"`
}
