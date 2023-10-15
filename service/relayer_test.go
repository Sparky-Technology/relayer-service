package service

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/swapbotAA/relayer/contracts/swapbotAA/router"
	"github.com/swapbotAA/relayer/ethereum"
	"github.com/swapbotAA/relayer/models"
	"github.com/swapbotAA/relayer/settings"
	"github.com/swapbotAA/relayer/utils"
	"log"
	"math/big"
	"strings"
	"testing"
)

var (
	swapVO = models.InstantSwapVO{
		TxValue: &models.TxValue{
			TokenIn:          "0xfFf9976782d46CC05630D1f6eBAb18b2324d6B14",
			TokenOut:         "0x1f9840a85d5aF5bf1D1762F925BDADdC4201F984",
			Fee:              3000,
			Recipient:        "0xb67e30aDb44c83E516681392FA16cD933B93b7ad",
			AmountIn:         "20000000000000",
			AmountOutMinimum: "7512000000000000",
		},
		User: "0xB178e99e401cBbd7F1a9bdafaa7D2D027B42d80A",
		R:    "0x92380612a7cd9f1a31771fabc250e0f165db105de44e676ae369062104bba50d",
		S:    "0x68fd217e8fff65b148be8367214173c7462171f24adc84ee4f14d07c28843361",
		V:    "0x1b",
	}
	salt = "0x4d3463635161376142546e52746b4d6972596447727a664b514a6d3858433747"
)

func TestInstantSwap(t *testing.T) {
	settings.Conf.EthereumConfig = &settings.EthereumConfig{}
	settings.Conf.EthereumConfig.RpcEndpointAddr = "wss://eth-sepolia.g.alchemy.com/v2/6Y_wOo59vIqBe8RD0NlOOZJwQtsvcaIf"
	settings.Conf.EthereumConfig.SystemAccountPrivateKey = "55bc18e822225d4a701ca9e4ead6f9696d794d1bfce3a8a18ac8b306b8cf3be3"
	settings.Conf.EthereumConfig.WalletContractAddr = "0x5c0B9D48f40d46634d1AA383CB15987708Ac39E6"
	settings.Conf.EthereumConfig.UniswapRouterContractAddr = "0xb67e30aDb44c83E516681392FA16cD933B93b7ad"

	instantSwap, err := initTestParam()
	if err != nil {
		t.Fatalf("Failed to init test params: %v", err)
	}
	t.Logf("init test params finish: instantSwap%v \n", instantSwap)

	// init ethereum config holder
	ethereum.Init()
	t.Logf("init ethereum config holder finish: EthereumConfig%v \n", ethereum.Config)

	client := ethereum.Config.Client
	if err != nil {
		t.Fatalf("Failed to connect to the Ethereum client: %v", err)
		return
	}
	ABI, _ := abi.JSON(strings.NewReader(router.AbiMetaData.ABI))
	params, err := AssembleExactInputSingleParams(instantSwap)
	if err != nil {
		t.Fatalf("failed to assemble ExactInputSingleParams: %s \n", err)
		return
	}
	t.Logf("assemble ExactInputSingleParams finish: ExactInputSingleParams%v \n", params)

	out, err := ABI.Pack("exactInputSingle", params, instantSwap.Salt, common.HexToAddress(instantSwap.User), instantSwap.V, instantSwap.R, instantSwap.S)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("abi pack finish: %v \n", out)

	// check balance
	checkBalance(instantSwap)

	// fail guard
	//err = failGuard(client, &uniswapRouterContractAddr, out)
	//if err != nil {
	//	log.Printf("connot pass fail guard check: %v \n", err)
	//	return
	//}
	conn := NewSwapbotAARouterConnecter(client)
	auth := ethereum.Config.Auth
	t.Logf("-------sendExactInputSingleSwapTx begin-------")
	// send tx
	conn.sendExactInputSingleSwapTx(auth, params, instantSwap)
	t.Logf("-------sendExactInputSingleSwapTx end-------")
}

func checkBalance(instantSwap *models.InstantSwapBO) {
	// 调用以太坊节点的eth_getBalance方法查询余额
	// 注意：余额是以wei为单位的，需要进行进一步的转换
	client := ethereum.Config.Client
	balance, err := client.BalanceAt(context.Background(), common.HexToAddress(instantSwap.User), nil)
	if err != nil {
		log.Fatalf("查询余额出错: %v", err)
	}

	// 将余额从wei转换为ether（以太币）
	etherBalance := new(big.Float).SetInt(balance)
	etherBalance = etherBalance.Quo(etherBalance, big.NewFloat(1e18))

	// 打印账户余额
	fmt.Printf("账户余额: %s ETH\n", etherBalance.String())
}

func initTestParam() (*models.InstantSwapBO, error) {
	swapBO := models.InstantSwapBO{
		TxValue: swapVO.TxValue,
		User:    swapVO.User,
	}
	var err error
	swapBO.Salt, err = utils.ConvertHexStringToByte32(salt)
	if err != nil {
		return &models.InstantSwapBO{}, err
	}
	swapBO.R, err = utils.ConvertHexStringToByte32(swapVO.R)
	if err != nil {
		return &models.InstantSwapBO{}, err
	}
	swapBO.S, err = utils.ConvertHexStringToByte32(swapVO.S)
	if err != nil {
		return &models.InstantSwapBO{}, err
	}
	swapBO.V, err = utils.ConvertHexStringToUint8(swapVO.V)
	if err != nil {
		return &models.InstantSwapBO{}, err
	}
	return &swapBO, nil
}
