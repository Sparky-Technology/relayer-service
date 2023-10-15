package bo

import (
	"encoding/json"
	"errors"
	"github.com/swapbotAA/relayer/controller/vo"
	"github.com/swapbotAA/relayer/utils"
	"math/big"
)

type LimitOrderBO struct {
	User             string
	TokenIn          string
	TokenOut         string
	Fee              int
	AmountIn         *big.Int
	AmountOut        *big.Int
	AmountOutMinimum *big.Int
}

func (bo *LimitOrderBO) ToString() string {
	jsonStr, _ := json.Marshal(bo)
	return string(jsonStr)
}

func (bo *LimitOrderBO) UnmarshalJSON(data []byte) (err error) {
	printRequestParam("[LimitOrder] Request param:", string(data))

	vo := vo.LimitOrderVO{}
	err = json.Unmarshal(data, &vo)

	if err != nil {
		return err
	} else if len(vo.User) == 0 {
		err = errors.New("缺少必填字段 user")
	} else if len(vo.TokenIn) == 0 {
		err = errors.New("缺少必填字段 token_in")
	} else if len(vo.TokenOut) == 0 {
		err = errors.New("缺少必填字段 token_out")
	} else if vo.Fee == 0 {
		err = errors.New("缺少必填字段 fee")
	} else if len(vo.AmountIn) == 0 {
		err = errors.New("缺少必填字段 amount_in")
	} else if len(vo.AmountOut) == 0 {
		err = errors.New("缺少必填字段 amount_out")
	} else if len(vo.AmountOutMinimum) == 0 {
		err = errors.New("缺少必填字段 amount_out_minimum")
	} else {
		bo.User = vo.User
		bo.TokenIn = vo.TokenIn
		bo.TokenOut = vo.TokenOut
		bo.Fee = vo.Fee
		bo.AmountIn, err = utils.ConvertStringToBigInt(vo.AmountIn)
		if err != nil {
			err = errors.New("amountIn 无法转换为大整数")
			return err
		}
		bo.AmountOut, err = utils.ConvertStringToBigInt(vo.AmountOut)
		if err != nil {
			err = errors.New("amountOut 无法转换为大整数")
			return err
		}
		bo.AmountOutMinimum, err = utils.ConvertStringToBigInt(vo.AmountOutMinimum)
		if err != nil {
			err = errors.New("amountOutMinimum 无法转换为大整数")
			return err
		}
	}
	return
}
