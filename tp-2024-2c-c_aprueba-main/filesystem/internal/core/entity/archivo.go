package entity

type Archivo struct {
	Pid       uint32 `json:"pid"`
	Tid       uint32 `json:"tid"`
	Tamanio   int    `json:"size"`
	Contenido []byte `json:"contenido"`
}

type Metadata struct { // para crear archivo metadata
	Index_Block string `json:"index_block"`
	Size        int    `json:"size"`
}
