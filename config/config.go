package config

type ChainConfig struct {
	ChainName       string
	ChainId         int
	EndpointUrl     string
	SupportedTokens map[string]string
}

const URL = "localhost:3000"

var ChainConfigs = []*ChainConfig{
	{
		ChainName:   "Polygon",
		ChainId:     137,
		EndpointUrl: "https://polygon.llamarpc.com",
		SupportedTokens: map[string]string{
			"USDC": "0x2791Bca1f2de4661ED88A30C99A7a9449Aa84174",
			"USDT": "0xc2132D05D31c914a87C6611C10748AEb04B58e8F",
			"WETH": "0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270",
		},
	},
	{
		ChainName:   "Gnosis Chain",
		ChainId:     100,
		EndpointUrl: "https://rpc.gnosischain.com",
		SupportedTokens: map[string]string{
			"USDC": "0xDDAfbb505ad214D7b80b1f830fcCc89B60fb7A83",
			"USDT": "0x4ECaBa5870353805a9F068101A40E0f32ed605C6",
			"WETH": "0x6A023CCd1ff6F2045C3309768eAd9E68F978f6e1",
		},
	},
	{
		ChainName:   "Arbitrum One",
		ChainId:     42161,
		EndpointUrl: "https://arbitrum.llamarpc.com",
		SupportedTokens: map[string]string{
			"USDC": "0xFF970A61A04b1cA14834A43f5dE4533eBDDB5CC8",
			"USDT": "0xFd086bC7CD5C481DCC9C85ebE478A1C0b69FCbb9",
			"WETH": "0x82aF49447D8a07e3bd95BD0d56f35241523fBab1",
		},
	},
	{
		ChainName:   "Base",
		ChainId:     8453,
		EndpointUrl: "https://base-mainnet.public.blastapi.io",
		SupportedTokens: map[string]string{
			"USDC": "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913",
		},
	},
}
