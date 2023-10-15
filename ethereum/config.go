package ethereum

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/swapbotAA/relayer/blockclient"
	"github.com/swapbotAA/relayer/settings"
	"log"
)

var Config ConfigHolder

type ConfigHolder struct {
	Auth   *bind.TransactOpts
	Client *blockclient.Client
	Wallet *Wallet

	PrivateClient *blockclient.Client
}

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  common.Address
}

func Init() error {
	rpcEndpointAddr := settings.Conf.EthereumConfig.RpcEndpointAddr
	client, err := blockclient.Dial(rpcEndpointAddr)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
		return err
	}

	// 解码私钥
	privateKey, err := crypto.HexToECDSA(settings.Conf.EthereumConfig.SystemAccountPrivateKey)
	if err != nil {
		log.Fatalf("failed to decode private key: %v", err)
		return err
	}

	// 通过私钥创建账户对象
	auth := bind.NewKeyedTransactor(privateKey)

	// 生成公钥 组装钱包对象
	wallet := &Wallet{
		PrivateKey: privateKey,
		PublicKey:  crypto.PubkeyToAddress(privateKey.PublicKey),
	}

	Config = ConfigHolder{Auth: auth, Client: client, Wallet: wallet}
	return nil
}
