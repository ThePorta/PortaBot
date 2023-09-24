package types

//go:generate msgp

type AccountAndInputData struct {
	AccountAddress string `json:"accountAddress"`
	TargetContract string `json:"targetContract"`
	InputData      string `json:"inputData"`
	ChainId        int    `json:"chainId"`
	ChainName      string `json:"chainName"`
}

type SetChatIdRequest struct {
	Opt     int    `json:"opt"`
	Address string `json:"address"`
}
