package entity

type IDs struct {
	PID uint32
	TID uint32
}

func CrearIDs(pid uint32, tid uint32) *IDs {
	return &IDs{
		PID: pid,
		TID: tid,
	}
}
