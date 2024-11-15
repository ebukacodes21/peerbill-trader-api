package utils

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
)

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
