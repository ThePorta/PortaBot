package config

type Token struct {
	Address string
	Name    string
}

type ChainConfig struct {
	ChainName       string
	ChainId         int
	EndpointUrl     string
	SupportedTokens []*Token
}

var ChainConfigs = []*ChainConfig{
	{
		ChainName:   "Polygon",
		ChainId:     137,
		EndpointUrl: "https://polygon.llamarpc.com",
		SupportedTokens: []*Token{
			{
				Address: "0x2791Bca1f2de4661ED88A30C99A7a9449Aa84174",
				Name:    "USDC",
			},
			{
				Address: "0xc2132D05D31c914a87C6611C10748AEb04B58e8F",
				Name:    "USDT",
			},
		},
	},
	{
		ChainName:   "Gnosis Chain",
		ChainId:     100,
		EndpointUrl: "https://rpc.gnosischain.com",
		SupportedTokens: []*Token{
			{
				Address: "0xDDAfbb505ad214D7b80b1f830fcCc89B60fb7A83",
				Name:    "USDC",
			},
			{
				Address: "0x4ECaBa5870353805a9F068101A40E0f32ed605C6",
				Name:    "USDT",
			},
		},
	},
	{
		ChainName:   "Arbitrum One",
		ChainId:     42161,
		EndpointUrl: "https://arbitrum.llamarpc.com",
		SupportedTokens: []*Token{
			{
				Address: "0xFF970A61A04b1cA14834A43f5dE4533eBDDB5CC8",
				Name:    "USDC",
			},
			{
				Address: "0xFd086bC7CD5C481DCC9C85ebE478A1C0b69FCbb9",
				Name:    "USDT",
			},
		},
	},
	{
		ChainName:   "Base",
		ChainId:     8453,
		EndpointUrl: "https://base-mainnet.public.blastapi.io",
		SupportedTokens: []*Token{
			{
				Address: "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913",
				Name:    "USDC",
			},
		},
	},
}
