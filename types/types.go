package types

//go:generate msgp

type AccountAndInputData struct {
	AccountAddress string
	InputData      []byte
}
