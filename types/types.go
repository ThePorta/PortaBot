package types

//go:generate msgp

type AccountInfo struct {
	WalletConnectSession string
	ChatId               int64
}

type Revoke struct {
	AccountAddress   string `json:"accountAddress"`
	TokenAddress     string `json:"tokenAddress"`
	MaliciousAddress string `json:"maliciousAddress"`
	ChatId           string `json:"chatId"`
}
