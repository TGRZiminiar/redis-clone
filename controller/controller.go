package controller

type (
	Status  byte
	Command byte

	ResponsetSet struct {
		Status Status
	}

	ResponsetGet struct {
		Status Status
		Value  []byte
	}
)

const (
	// -------- Status --------
	StatusNone Status = iota
	StatusOk
	StatusError
	StatusNotFound

	// -------- Command --------
	CMDNonce Command = iota
	CMDSet
	CMDGet
	CMDDel
	CMDJoin
)
