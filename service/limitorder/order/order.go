package order

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/swapbotAA/relayer/blockclient"
	uniswapV3Factory "github.com/swapbotAA/relayer/contracts/uniswapV3/factory"
	uniswapV3Pool "github.com/swapbotAA/relayer/contracts/uniswapV3/pool"
	uniswapV3Router "github.com/swapbotAA/relayer/contracts/uniswapV3/router"
	"github.com/swapbotAA/relayer/ethereum"
	"github.com/swapbotAA/relayer/service/limitorder/bo"
	"github.com/swapbotAA/relayer/settings"
	"github.com/swapbotAA/relayer/utils"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	// The interval to check the price in Milliseconds
	CheckIntervalMil time.Duration = 500
	// The deadline for completing the swap in Minutes
	SwapDeadlineMin  time.Duration = 20
	PriceChannelSize int           = 100
)

var (
	ErrConv        = errors.New("could not convert amount to big int")
	ErrInvalidType = errors.New("invalid type value")
	ErrCanContinue = errors.New("can continue execution")
)

type Type int

const (
	Limit Type = iota
	StopLimit
)

// Custom json unmarshal for Type
func (t *Type) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	switch strings.ToLower(s) {
	case "limit":
		*t = Limit
	case "stop_limit":
		*t = StopLimit
	default:
		return ErrInvalidType
	}

	return nil
}

// Custom json unmarshal for Type
func (t Type) MarshalJSON() ([]byte, error) {
	var s string
	switch t {
	case Limit:
		s = "limit"
	case StopLimit:
		s = "stop_limit"
	default:
		return nil, ErrInvalidType
	}

	return json.Marshal(s)
}

type Router struct {
	Address  common.Address
	Instance *uniswapV3Router.Abi
}

type Pool struct {
	Instance *uniswapV3Pool.Abi
}

type Order struct {
	ID               string
	Type             Type
	UserID           string
	TokenIn          string
	TokenOut         string
	Fee              int
	AmountIn         *big.Int
	AmountOut        *big.Int
	AmountOutMinimum *big.Int
	Price            *big.Float
	router           *Router
	pool             *Pool
}

// UniswapRouterConnecter Router合约连接者
type UniswapRouterConnecter struct {
	ctx           context.Context
	client        *blockclient.Client
	UniswapRouter *uniswapV3Router.Abi
}

func NewUniswapRouterConnecter(client *blockclient.Client) *UniswapRouterConnecter {
	uniswapRouter, err := uniswapV3Router.NewAbi(common.HexToAddress(settings.Conf.EthereumConfig.UniswapRouterContractAddr), client)
	if err != nil {
		log.Fatalf("Failed to NewAbi: %v", err)
	}
	return &UniswapRouterConnecter{
		ctx:           context.Background(),
		client:        client,
		UniswapRouter: uniswapRouter,
	}
}

// Initializes the Order
func (o *Order) Init(bo *bo.LimitOrderBO) error {
	o.router = &Router{Address: common.HexToAddress(settings.Conf.EthereumConfig.UniswapRouterContractAddr)}
	var err error
	o.router.Instance, err = uniswapV3Router.NewAbi(o.router.Address, ethereum.Config.Client)
	if err != nil {
		log.Fatalf("Failed to NewAbi: %v", err)
	}
	o.pool.Instance, err = getPool(ethereum.Config.Client, o.router.Instance, common.HexToAddress(bo.TokenIn), common.HexToAddress(bo.TokenOut))
	if err != nil {
		log.Fatalf("Failed to getPool: %v", err)
	}
	token0, _ := o.pool.Instance.Token0(nil)
	if bo.TokenIn == token0.Hex() {
		o.Type = Limit
		o.Price = utils.DivideBigIntsToBigFloat(bo.AmountIn, bo.AmountOut)
	} else {
		o.Type = StopLimit
		o.Price = utils.DivideBigIntsToBigFloat(bo.AmountOut, bo.AmountIn)
	}
	o.TokenIn = bo.TokenIn
	o.TokenOut = bo.TokenOut
	o.Fee = bo.Fee
	o.AmountIn = bo.AmountIn
	o.AmountOut = bo.AmountOut
	o.AmountOutMinimum = bo.AmountOutMinimum
	return nil
}

// Executes the approve transaction with the exact order sell amount
func (o *Order) ApproveAmount(c *blockclient.Client) (*types.Receipt, error) {
	// 调用 entryPoint 合约
	return nil, nil
}

// Compares the given amount with the wanted amount and returns the transaction or nil
// if no transaction has been made
func (o *Order) CompAndSwap(c *blockclient.Client, currentPrice *big.Float) (*types.Receipt, error) {

	if o.Type == Limit && currentPrice.Cmp(o.Price) != -1 {
		return o.Swap(c, currentPrice)
	}

	if o.Type == StopLimit && currentPrice.Cmp(o.Price) != 1 {
		return o.Swap(c, currentPrice)
	}

	return nil, nil
}

// Executes the swap transaction
func (o *Order) Swap(c *blockclient.Client, currentPrice *big.Float) (*types.Receipt, error) {
	conn := NewUniswapRouterConnecter(c)
	auth := ethereum.Config.Auth
	var err error
	auth.GasPrice, err = c.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	// send tx
	param, err := assembleExactInputSingleParams(o)
	tx, err := conn.UniswapRouter.ExactInputSingle(auth, param)
	if err != nil {
		return nil, err
	}

	return c.GetPendingTxReceipt(context.Background(), tx)
}

func assembleExactInputSingleParams(o *Order) (uniswapV3Router.ISwapRouterExactInputSingleParams, error) {
	params := uniswapV3Router.ISwapRouterExactInputSingleParams{
		TokenIn:           common.HexToAddress(o.TokenIn),
		TokenOut:          common.HexToAddress(o.TokenOut),
		Fee:               big.NewInt(int64(o.Fee)),
		Recipient:         ethereum.Config.Wallet.PublicKey,
		AmountIn:          o.AmountIn,
		AmountOutMinimum:  o.AmountOutMinimum,
		SqrtPriceLimitX96: big.NewInt(0),
	}
	return params, nil
}

func sqrtPriceX96ToPrice(sqrtPriceX96 *big.Int) *big.Float {
	bigFloat := new(big.Float).SetInt(sqrtPriceX96)
	price := bigFloat.Quo(bigFloat, new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(2), big.NewInt(96), nil)))
	return price
}

// Returns a price channel for receiving a stream of prices and a channel for error.
// Stops only when it receives from the done channel
func (o *Order) GetPriceStream() (chan *big.Float, chan error) {
	priceChan, errChan := make(chan *big.Float, PriceChannelSize), make(chan error)

	go func() {
		for {
			time.Sleep(time.Second)
			slot0Result, err := o.pool.Instance.Slot0(nil)
			if err != nil {
				errChan <- err
				return
			}
			price := sqrtPriceX96ToPrice(slot0Result.SqrtPriceX96)

			priceChan <- price
		}
	}()

	return priceChan, errChan
}

// Returns uniswapV3 pool contract instance for the tokens with the given addresses
func getPool(c bind.ContractBackend, r *uniswapV3Router.Abi, addr1, addr2 common.Address) (*uniswapV3Pool.Abi, error) {
	factoryAddress, err := r.Factory(nil)
	if err != nil {
		return nil, err
	}

	factoryInstance, err := uniswapV3Factory.NewAbi(factoryAddress, c)
	if err != nil {
		return nil, err
	}

	poolAddress, err := factoryInstance.GetPool(nil, addr1, addr2, big.NewInt(500))
	if err != nil {
		return nil, err
	}

	return uniswapV3Pool.NewAbi(poolAddress, c)
}
