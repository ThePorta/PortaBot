package types

//go:generate msgp

type AccountInfo struct {
	WalletConnectSession string
	ChatId               string
}
