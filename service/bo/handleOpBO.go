package bo

import (
	"encoding/json"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/swapbotAA/relayer/controller/vo"
	"github.com/swapbotAA/relayer/utils"
	"math/big"
)

type HandleOpBO struct {
	UserOp      *UserOperationBO
	Beneficiary common.Address
}

func (bo *HandleOpBO) ToString() string {
	jsonStr, _ := json.Marshal(bo)
	return string(jsonStr)
}

func (bo *HandleOpBO) UnmarshalJSON(data []byte) (err error) {
	printRequestParam("[HandleOp] Request param:", string(data))

	vo := vo.HandleOpVO{}
	err = json.Unmarshal(data, &vo)

	if err != nil {
		return err
	} else if vo.UserOp == nil {
		err = errors.New("缺少必填字段 user_op")
	} else if len(vo.BeneficiaryAddr) == 0 {
		err = errors.New("缺少必填字段 beneficiary_addr")
	} else {
		bo.UserOp, err = convertUserOperationVOToBO(vo.UserOp)
		if err != nil {
			return err
		}
		bo.Beneficiary = common.HexToAddress(vo.BeneficiaryAddr)
	}
	return
}

type UserOperationBO struct {
	Sender               common.Address
	Nonce                *big.Int
	InitCode             []byte
	CallData             []byte
	CallGasLimit         *big.Int
	VerificationGasLimit *big.Int
	PreVerificationGas   *big.Int
	MaxFeePerGas         *big.Int
	MaxPriorityFeePerGas *big.Int
	PaymasterAndData     []byte
	Signature            []byte
}

func convertUserOperationVOToBO(vo *vo.UserOperationVO) (bo *UserOperationBO, err error) {
	bo = &UserOperationBO{}
	bo.Sender = common.HexToAddress(vo.Sender)
	bo.Nonce, err = utils.ConvertStringToBigInt(vo.Nonce)
	if err != nil {
		err = errors.New("nonce 无法转换为大整数")
		return nil, err
	}
	bo.InitCode, err = utils.ConvertHexStringToBytes(vo.InitCode)
	if err != nil {
		err = errors.New("initCode 无法转换为字节数组")
		return nil, err
	}
	bo.CallData, err = utils.ConvertHexStringToBytes(vo.CallData)
	if err != nil {
		err = errors.New("callData 无法转换为字节数组")
		return nil, err
	}
	bo.CallGasLimit, err = utils.ConvertStringToBigInt(vo.CallGasLimit)
	if err != nil {
		err = errors.New("callGasLimit 无法转换为大整数")
		return nil, err
	}
	bo.VerificationGasLimit, err = utils.ConvertStringToBigInt(vo.VerificationGasLimit)
	if err != nil {
		err = errors.New("verificationGasLimit 无法转换为大整数")
		return nil, err
	}
	bo.PreVerificationGas, err = utils.ConvertStringToBigInt(vo.PreVerificationGas)
	if err != nil {
		err = errors.New("preVerificationGas 无法转换为大整数")
		return nil, err
	}
	bo.MaxFeePerGas, err = utils.ConvertStringToBigInt(vo.MaxFeePerGas)
	if err != nil {
		err = errors.New("maxFeePerGas 无法转换为大整数")
		return nil, err
	}
	bo.MaxPriorityFeePerGas, err = utils.ConvertStringToBigInt(vo.MaxPriorityFeePerGas)
	if err != nil {
		err = errors.New("maxPriorityFeePerGas 无法转换为大整数")
		return nil, err
	}
	bo.PaymasterAndData, err = utils.ConvertHexStringToBytes(vo.PaymasterAndData)
	if err != nil {
		err = errors.New("paymasterAndData 无法转换为字节数组")
		return nil, err
	}
	bo.Signature, err = utils.ConvertHexStringToBytes(vo.Signature)
	if err != nil {
		err = errors.New("signature 无法转换为字节数组")
		return nil, err
	}
	return bo, nil
}
