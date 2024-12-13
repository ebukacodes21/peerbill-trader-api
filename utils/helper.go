package utils

import (
	"context"
	"crypto/ecdsa"

	"fmt"
	"log"
	"math"
	"math/big"
	"math/rand"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
)

const ERC20_ABI = `[{
	"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"payable":false,"stateMutability":"view","type":"function"
},{
	"constant":true,"inputs":[{"name":"account","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"
},{
	"constant":false,"inputs":[{"name":"recipient","type":"address"},{"name":"amount","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"
},{
	"constant":false,"inputs":[{"name":"sender","type":"address"},{"name":"recipient","type":"address"},{"name":"amount","type":"uint256"}],"name":"transferFrom","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"
},{
	"constant":true,"inputs":[{"name":"spender","type":"address"},{"name":"amount","type":"uint256"}],"name":"approve","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"
}]`

const alpha = "abcdefghijklmnopqrstuvwxyz"

var randSource = rand.New(rand.NewSource(time.Now().UnixNano()))

func init() {
	rand.Seed(time.Now().UnixMicro())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder

	for i := 0; i < n; i++ {
		char := alpha[rand.Intn(len(alpha))]
		sb.WriteByte(char)
	}

	return sb.String()
}

func RandomPhone() string {
	var sb strings.Builder

	for i := 0; i < 10; i++ {
		sb.WriteByte(byte(randSource.Intn(10) + '0'))
	}

	return sb.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomFloat(min, max float64, precision int) float64 {
	factor := math.Pow(10, float64(precision))
	randValue := min + randSource.Float64()*(max-min)
	return math.Round(randValue*factor) / factor
}

func RandomFiat() string {
	currencies := []string{"GBP", "NGN", "USD", "CAD", "AUD"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

func RandomCrypto() string {
	currencies := []string{"BTC", "ETH", "BNB", "USDT", "USDC"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}

func CheckBalance(ctx context.Context, crypto string, address string) string {
	var balance big.Int
	var formattedBalance string

	switch crypto {
	case "ETH":
		balance = getBalance(ctx, "https://eth-mainnet.g.alchemy.com/v2/xtngPwzqpjqcWBvKoVKLiBwYo1kTbWxe", address)
		formattedBalance = weiToEther(balance, 18)
	case "BNB":
		balance = getBalance(ctx, "https://bsc-dataseed.binance.org/", address)
		formattedBalance = weiToEther(balance, 18)
	case "USDT", "USDC":
		formattedBalance = getTokenBalance(ctx, "https://bsc-dataseed.binance.org/", crypto, address)
	default:
		log.Printf("Unsupported crypto: %s", crypto)
		return "Unsupported crypto"
	}

	return formattedBalance
}

// get balance of native currencies
func getBalance(ctx context.Context, provider string, address string) big.Int {
	client, err := ethclient.Dial(provider)
	if err != nil {
		log.Fatal(err)
	}

	account := common.HexToAddress(address)
	balance, err := client.BalanceAt(ctx, account, nil)
	if err != nil {
		log.Fatal(err)
		return big.Int{}
	}

	return *balance
}

// get balance of token currencies
func getTokenBalance(ctx context.Context, provider string, crypto string, address string) string {
	var contractAddr string
	var decimals uint8 = 6
	var balance *big.Int

	// USDT and USDC smart contract addresses on BNB
	switch crypto {
	case "USDT":
		contractAddr = "0x55d398326f99059fF775485246999027B3197955"
	case "USDC":
		contractAddr = "0x8ac76a51cc950d9822d68b83fe1ad97b32cd580d"
	default:
		log.Printf("Unsupported token: %s", crypto)
		return "Unsupported token"
	}

	client, err := ethclient.Dial(provider)
	if err != nil {
		log.Printf("Failed to connect to provider %v", err)
		return "Error connecting to provider"
	}

	addr := common.HexToAddress(contractAddr)
	contractABI, err := abi.JSON(strings.NewReader(ERC20_ABI))
	if err != nil {
		log.Printf("Failed to parse contract ABI: %v", err)
		return "Error parsing contract ABI"
	}

	// call balanceOf func - building a contract object
	callData, err := contractABI.Pack("balanceOf", common.HexToAddress(address))
	if err != nil {
		log.Printf("Failed to call balanceOf method: %v", err)
		return "Error fetching balance"
	}
	result, err := client.CallContract(ctx, ethereum.CallMsg{To: &addr, Data: callData}, nil)
	if err != nil {
		log.Printf("Failed to call balanceOf method: %v", err)
		return "Error fetching balance"
	}

	// Decode the balance
	err = contractABI.UnpackIntoInterface(&balance, "balanceOf", result)
	if err != nil {
		log.Printf("Failed to decode balance result: %v", err)
		return "Error decoding balance"
	}

	// Convert the balance from the smallest unit to the token's standard unit - USDT, USDC is 6 decimal places
	decimalFactor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	realBalance := new(big.Rat).SetFrac(balance, decimalFactor)

	// Return the balance as a string with the appropriate decimal places
	return realBalance.FloatString(int(decimals))
}

// Converts balance from Wei to Ether (18 decimals)
func weiToEther(balance big.Int, decimals int) string {
	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	etherBalance := new(big.Rat).SetFrac(&balance, divisor)
	return etherBalance.FloatString(decimals)
}

func SendNative(provider string, privateKeyString string, receiverAddress string, amount big.Int) {
	client, err := ethclient.Dial(provider)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA(privateKeyString)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal(err)
	}

	sender := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), sender)
	if err != nil {
		log.Fatal(err)
	}

	gasLimit := uint64(21000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	receiver := common.HexToAddress(receiverAddress)

	// unsigned tx ===
	tx := types.NewTransaction(nonce, receiver, &amount, gasLimit, gasPrice, nil)
	// sign tx
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
}

func SendTokens(provider string, privateKeyString string, receiverAddr string, contractAddr string, tokenAmount big.Int) {
	client, err := ethclient.Dial(provider)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA(privateKeyString)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal(err)
	}

	sender := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), sender)
	if err != nil {
		log.Fatal(err)
	}

	// gasLimit := uint64(21000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	value := big.NewInt(0)
	toAddress := common.HexToAddress(receiverAddr)
	tokenAddress := common.HexToAddress(contractAddr)

	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	fmt.Println(hexutil.Encode(methodID))

	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	amount := tokenAmount
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &toAddress,
		Data: data,
	})
	if err != nil {
		log.Fatal(err)
	}

	tx := types.NewTransaction(nonce, tokenAddress, value, gasLimit, gasPrice, data)
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
}
