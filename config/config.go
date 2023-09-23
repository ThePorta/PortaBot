package config

type Token struct {
	Address string
	Name    string
}

type ChainConfig struct {
	ChainId                 int
	EndpointUrl             string
	SupportedTokenAddresses []*Token
}

var ChainConfigs = []*ChainConfig{
	{
		ChainId:     137,
		EndpointUrl: "https://polygon.llamarpc.com",
		SupportedTokenAddresses: []*Token{
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
		ChainId:     100,
		EndpointUrl: "https://rpc.gnosischain.com",
		SupportedTokenAddresses: []*Token{
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
		ChainId:     42161,
		EndpointUrl: "https://arbitrum.llamarpc.com",
		SupportedTokenAddresses: []*Token{
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
		ChainId:     8453,
		EndpointUrl: "https://base-mainnet.public.blastapi.io",
		SupportedTokenAddresses: []*Token{
			{
				Address: "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913",
				Name:    "USDC",
			},
		},
	},
}
