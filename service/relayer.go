package service

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/swapbotAA/relayer/blockclient"
	"github.com/swapbotAA/relayer/contracts/swapbotAA/router"
	"github.com/swapbotAA/relayer/contracts/swapbotAA/wallet"
	"github.com/swapbotAA/relayer/ethereum"
	"github.com/swapbotAA/relayer/settings"
	"github.com/swapbotAA/relayer/utils"
	"go.uber.org/zap"
	"log"
	"math/big"
	"strings"

	bo "github.com/swapbotAA/relayer/service/bo"
)

// WalletConnecter Wallet合约连接者
type WalletConnecter struct {
	ctx    context.Context
	client *blockclient.Client
	wallet *wallet.Wallet
}

func NewWalletConnecter(client *blockclient.Client) *WalletConnecter {
	wallet, err := wallet.NewWallet(common.HexToAddress(settings.Conf.EthereumConfig.WalletContractAddr), client)
	if err != nil {
		log.Fatalf("Failed to NewWallet: %v", err)
	}
	return &WalletConnecter{
		ctx:    context.Background(),
		client: client,
		wallet: wallet,
	}
}

// SwapbotAARouterConnecter Router合约连接者
type SwapbotAARouterConnecter struct {
	ctx           context.Context
	client        *blockclient.Client
	uniswapRouter *router.Abi
}

func NewSwapbotAARouterConnecter(client *blockclient.Client) *SwapbotAARouterConnecter {
	swapbotAARouter, err := router.NewAbi(common.HexToAddress(settings.Conf.EthereumConfig.UniswapRouterContractAddr), client)
	if err != nil {
		log.Fatalf("Failed to NewAbi: %v", err)
	}
	return &SwapbotAARouterConnecter{
		ctx:           context.Background(),
		client:        client,
		uniswapRouter: swapbotAARouter,
	}
}

// sendExactInputSingleSwapTx uniswapRouter转账
func (c *SwapbotAARouterConnecter) sendExactInputSingleSwapTx(fromAuth *bind.TransactOpts, param router.IV3SwapRouterExactInputSingleParams, swapBO *bo.InstantSwapBO) (string, error) {
	tx, err := c.uniswapRouter.ExactInputSingle(fromAuth, param, swapBO.Salt, common.HexToAddress(swapBO.User), swapBO.V, swapBO.R, swapBO.S)
	if err != nil {
		return "", err
	}
	//// 等待执行
	//receipt, err := bind.WaitMined(c.ctx, c.client, tx)
	//if err != nil {
	//	log.Fatalf("Wait Transfer Mined error: %v", err)
	//	return false
	//}
	zap.L().Info("[InstantSwap] Send transaction success!", zap.String("tx", tx.Hash().String()))
	return tx.Hash().String(), nil
}

func HandleInstantSwap(swapBO *bo.InstantSwapBO) (string, error) {
	client := ethereum.Config.Client
	// abi
	ABI, _ := abi.JSON(strings.NewReader(router.AbiMetaData.ABI))
	param, err := AssembleExactInputSingleParams(swapBO)
	if err != nil {
		zap.L().Error("[InstantSwap] Failed to assemble ExactInputSingleParams!", zap.Error(err))
		return "", err
	}
	out, err := ABI.Pack("exactInputSingle", param, swapBO.Salt, common.HexToAddress(swapBO.User), swapBO.V, swapBO.R, swapBO.S)
	if err != nil {
		log.Fatal(err)
	}

	// fail guard
	gasPrice, err := failGuard("InstantSwap", client, ethereum.Config.Wallet.PublicKey, common.HexToAddress(settings.Conf.UniswapRouterContractAddr), out)
	if err != nil {
		zap.L().Error("[InstantSwap] Not pass FailGuard check!", zap.Error(err))
		return "", utils.NewError(int64(utils.CodeNotPassFailGuardCheck), err.Error())
	}

	conn := NewSwapbotAARouterConnecter(client)
	auth := ethereum.Config.Auth
	auth.GasPrice = gasPrice

	// send tx
	txHash, err := conn.sendExactInputSingleSwapTx(auth, param, swapBO)
	if err != nil {
		zap.L().Error("[InstantSwap] Send transaction fail!",
			zap.String("user", swapBO.User),
			zap.String("token_out", swapBO.TxValue.TokenOut),
			zap.String("token_in", swapBO.TxValue.TokenIn),
			zap.String("amount", swapBO.TxValue.AmountIn.String()),
			zap.Error(err),
		)
		return "", utils.NewError(int64(utils.CodeSendTxFail), err.Error())
	}
	return txHash, err
}

func AssembleExactInputSingleParams(swapBO *bo.InstantSwapBO) (router.IV3SwapRouterExactInputSingleParams, error) {
	params := router.IV3SwapRouterExactInputSingleParams{
		TokenIn:           common.HexToAddress(swapBO.TxValue.TokenIn),
		TokenOut:          common.HexToAddress(swapBO.TxValue.TokenOut),
		Fee:               big.NewInt(int64(swapBO.TxValue.Fee)),
		Recipient:         common.HexToAddress(swapBO.TxValue.Recipient),
		AmountIn:          swapBO.TxValue.AmountIn,
		AmountOutMinimum:  swapBO.TxValue.AmountOutMinimum,
		SqrtPriceLimitX96: big.NewInt(0),
	}
	return params, nil
}
