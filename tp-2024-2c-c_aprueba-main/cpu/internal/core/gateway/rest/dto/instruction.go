package dto

type Instruction struct {
	Code string   `json:"code"` // Code del proceso
	Args []string `json:"args"` // Args
}
