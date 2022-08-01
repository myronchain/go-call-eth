package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"go-call-eth/contracts"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	// client
	rawUrl          = "https://polygon-mumbai.g.alchemy.com/v2/GmqREz-eRvN6r_lVW8aTY3VipXvn3qKL"
	contractAddress = "0x4dede4c8ce79681cae880e54bb670e32673883a3"
	// 0x75965652aC872E3A9fd3c9ad8E290473c23d763c
	hexKeyPrivateKey = "fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19"
	// Mumbai 测试网
	chainId = 80001
)

func main() {
	client, err := ethclient.Dial(rawUrl)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA(hexKeyPrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	transactOpts, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(int64(chainId)))
	if err != nil {
		log.Fatal(err)
	}
	transactOpts.Nonce = big.NewInt(int64(nonce))
	transactOpts.Value = big.NewInt(0)
	transactOpts.GasLimit = uint64(300000)
	transactOpts.GasPrice = gasPrice

	address := common.HexToAddress(contractAddress)
	instance, err := contracts.NewStorage(address, client)
	if err != nil {
		log.Fatal(err)
	}

	storeValue := big.NewInt(10)

	tx, err := instance.Store(transactOpts, storeValue)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s\n", tx.Hash().Hex())

	result, err := instance.Retrieve(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result.Int64())

	fmt.Println(storeValue.Int64() == result.Int64())
}
