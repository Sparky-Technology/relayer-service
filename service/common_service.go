package service

import (
	"context"
	goethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/swapbotAA/relayer/blockclient"
	entrypoint "github.com/swapbotAA/relayer/contracts/swapbotAA/entrypoint"
	"github.com/swapbotAA/relayer/ethereum"
	"github.com/swapbotAA/relayer/service/bo"
	"github.com/swapbotAA/relayer/settings"
	"github.com/swapbotAA/relayer/utils"
	"go.uber.org/zap"
	"log"
	"math/big"
	"strings"
)

type EntryPointConnecter struct {
	ctx        context.Context
	client     *blockclient.Client
	entryPoint *entrypoint.Abi
}

func NewEntryPointConnecter(client *blockclient.Client) *EntryPointConnecter {
	entryPoint, err := entrypoint.NewAbi(common.HexToAddress(settings.Conf.EthereumConfig.EntryPointContractAddr), client)
	if err != nil {
		log.Fatalf("Failed to NewAbi: %v", err)
	}
	return &EntryPointConnecter{
		ctx:        context.Background(),
		client:     client,
		entryPoint: entryPoint,
	}
}

// sendExactInputSingleSwapTx uniswapRouter转账
func (c *EntryPointConnecter) handleOp(fromAuth *bind.TransactOpts, userOp *bo.UserOperationBO, beneficiary common.Address) (string, error) {
	userOps := []entrypoint.UserOperation{
		convertBOToUserOperation(userOp),
	}
	tx, err := c.entryPoint.HandleOps(fromAuth, userOps, beneficiary)
	if err != nil {
		return "", err
	}
	//// 等待执行
	//receipt, err := bind.WaitMined(c.ctx, c.client, tx)
	//if err != nil {
	//	log.Fatalf("Wait Transfer Mined error: %v", err)
	//	return false
	//}
	zap.L().Info("[HandleOp] Send transaction success!", zap.String("tx", tx.Hash().String()))
	return tx.Hash().String(), nil
}

func convertBOToUserOperation(bo *bo.UserOperationBO) entrypoint.UserOperation {
	userOp := entrypoint.UserOperation{}
	userOp.Sender = bo.Sender
	userOp.Nonce = bo.Nonce
	userOp.InitCode = bo.InitCode
	userOp.CallData = bo.CallData
	userOp.CallGasLimit = bo.CallGasLimit
	userOp.VerificationGasLimit = bo.VerificationGasLimit
	userOp.PreVerificationGas = bo.PreVerificationGas
	userOp.MaxFeePerGas = bo.MaxFeePerGas
	userOp.MaxPriorityFeePerGas = bo.MaxPriorityFeePerGas
	userOp.PaymasterAndData = bo.PaymasterAndData
	userOp.Signature = bo.Signature
	return userOp
}

func HandleOp(bo *bo.HandleOpBO) (string, error) {
	client := ethereum.Config.Client
	// abi
	ABI, _ := abi.JSON(strings.NewReader(entrypoint.AbiMetaData.ABI))
	userOps := []entrypoint.UserOperation{
		convertBOToUserOperation(bo.UserOp),
	}
	out, err := ABI.Pack("handleOps", userOps, bo.UserOp.Sender)
	if err != nil {
		log.Fatal(err)
	}

	// fail guard
	gasPrice, err := failGuard("HandleOp", client, ethereum.Config.Wallet.PublicKey, common.HexToAddress(settings.Conf.EntryPointContractAddr), out)
	if err != nil {
		zap.L().Error("[HandleOp] Not pass FailGuard check!", zap.Error(err))
		return "", utils.NewError(int64(utils.CodeNotPassFailGuardCheck), err.Error())
	}

	conn := NewEntryPointConnecter(client)
	auth := ethereum.Config.Auth
	auth.GasPrice = gasPrice

	// send tx
	txHash, err := conn.handleOp(auth, bo.UserOp, bo.Beneficiary)
	if err != nil {
		zap.L().Error("[HandleOp] Send transaction fail!",
			zap.String("sender", bo.UserOp.Sender.Hex()),
			zap.String("nonce", bo.UserOp.Nonce.String()),
			zap.Error(err),
		)
		return "", utils.NewError(int64(utils.CodeSendTxFail), err.Error())
	}
	return txHash, err
}

func failGuard(tag string, client *blockclient.Client, from common.Address, to common.Address, out []byte) (*big.Int, error) {
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	zap.L().Info("["+tag+"] Query SuggestGasPrice finish.", zap.String("gasPrice", gasPrice.String()))

	estimateGas, err := client.EstimateGas(context.Background(), goethereum.CallMsg{
		From:     from,
		To:       &to,
		GasPrice: gasPrice,
		Value:    big.NewInt(0),
		Data:     out,
	})
	if err != nil {
		return nil, err
	}
	zap.L().Info("["+tag+"] EstimateGas finish.", zap.Uint64("estimateGas", estimateGas))

	return gasPrice, nil
}
