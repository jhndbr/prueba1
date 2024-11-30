package dto

import (
	"github.com/sisoputnfrba/tp-golang/kernel/internal/core/entity"
	"github.com/sisoputnfrba/tp-golang/utils/list"
)

type ProcesoDTO struct {
	Process *entity.PCB `json:"processId"`
}

type ProcesosDTO struct {
	List list.ArrayList[*entity.PCB] `json:"process"`
}
