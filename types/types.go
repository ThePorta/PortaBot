package types

//go:generate msgp

type AccountAndInputData struct {
	AccountAddress string `json:"accountAddress"`
	TargetContract string `json:"targetContract"`
	InputData      []byte `json:"inputData"`
	ChainId        int    `json:"chainId"`
	ChainName      string `json:"chainName"`
}
