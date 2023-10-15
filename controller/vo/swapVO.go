package vo

type InstantSwapVO struct {
	TxValue *TxValueVO `json:"tx_value" binding:"required"`
	User    string     `json:"user" binding:"required"`
	V       string     `json:"v" binding:"required"`
	R       string     `json:"r" binding:"required"`
	S       string     `json:"s" binding:"required"`
	Salt    string     `json:"salt" binding:"required"`
}

type TxValueVO struct {
	TokenIn          string `json:"token_in" binding:"required"`
	TokenOut         string `json:"token_out" binding:"required"`
	Fee              int    `json:"fee" binding:"required"`
	Recipient        string `json:"recipient" binding:"required"`
	AmountIn         string `json:"amount_in" binding:"required"`
	AmountOutMinimum string `json:"amount_out_minimum" binding:"required"`
}
