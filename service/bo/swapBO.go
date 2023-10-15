package bo

import (
	"encoding/json"
	"errors"
	"github.com/swapbotAA/relayer/controller/vo"
	"github.com/swapbotAA/relayer/utils"
	"math/big"
)

type InstantSwapBO struct {
	TxValue *TxValueBO
	User    string
	V       uint8
	R       [32]byte
	S       [32]byte
	Salt    [32]byte
}

func (bo *InstantSwapBO) ToString() string {
	jsonStr, _ := json.Marshal(bo)
	return string(jsonStr)
}

func (bo *InstantSwapBO) UnmarshalJSON(data []byte) (err error) {
	printRequestParam("[InstantSwap] Request param:", string(data))

	vo := vo.InstantSwapVO{}
	err = json.Unmarshal(data, &vo)

	if err != nil {
		return err
	} else if vo.TxValue == nil {
		err = errors.New("缺少必填字段 tx_value")
	} else if len(vo.User) == 0 {
		err = errors.New("缺少必填字段 user")
	} else if len(vo.V) == 0 {
		err = errors.New("缺少必填字段 v")
	} else if len(vo.R) == 0 {
		err = errors.New("缺少必填字段 r")
	} else if len(vo.S) == 0 {
		err = errors.New("缺少必填字段 s")
	} else if len(vo.Salt) == 0 {
		err = errors.New("缺少必填字段 salt")
	} else {
		bo.TxValue, err = convertTxValueVOToBO(vo.TxValue)
		if err != nil {
			return err
		}
		bo.User = vo.User
		bo.V, err = utils.ConvertHexStringToUint8(vo.V)
		if err != nil {
			err = errors.New("v 解析失败")
			return err
		}
		bo.R, err = utils.ConvertHexStringToByte32(vo.R)
		if err != nil {
			err = errors.New("r 解析失败")
			return err
		}
		bo.S, err = utils.ConvertHexStringToByte32(vo.S)
		if err != nil {
			err = errors.New("s 解析失败")
			return err
		}
		bo.Salt, err = utils.ConvertHexStringToByte32(vo.Salt)
		if err != nil {
			err = errors.New("salt 解析失败")
			return err
		}
	}
	return
}

type TxValueBO struct {
	TokenIn          string
	TokenOut         string
	Fee              int
	Recipient        string
	AmountIn         *big.Int
	AmountOutMinimum *big.Int
}

func convertTxValueVOToBO(vo *vo.TxValueVO) (bo *TxValueBO, err error) {
	bo = &TxValueBO{}
	bo.TokenIn = vo.TokenIn
	bo.TokenOut = vo.TokenOut
	bo.Fee = vo.Fee
	bo.Recipient = vo.Recipient
	bo.AmountIn, err = utils.ConvertStringToBigInt(vo.AmountIn)
	if err != nil {
		err = errors.New("amountIn 无法转换为大整数")
		return nil, err
	}
	bo.AmountOutMinimum, err = utils.ConvertStringToBigInt(vo.AmountOutMinimum)
	if err != nil {
		err = errors.New("amountOutMinimum 无法转换为大整数")
		return nil, err
	}
	return bo, nil
}
