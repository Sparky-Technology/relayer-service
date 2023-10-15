package vo

type LimitOrderVO struct {
	User             string `json:"user" binding:"required"`
	TokenIn          string `json:"token_in" binding:"required"`
	TokenOut         string `json:"token_out" binding:"required"`
	Fee              int    `json:"fee" binding:"required"`
	AmountIn         string `json:"amount_in" binding:"required"`
	AmountOut        string `json:"amount_out" binding:"required"`
	AmountOutMinimum string `json:"amount_out_minimum" binding:"required"`
}
