package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"sync"
	"math"
	"math/big"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	// "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
)

const (
	coinGeckoAPIURL = "https://api.coingecko.com/api/v3"
	updateInterval = 5 * time.Minute
)

var (
	exchangeRatesMutex      sync.Mutex
	exchangeRateHistoryMutex sync.Mutex
	exchangeRateHistory map[string][]ExchangeRateHistory
	exchangeRates map[string]map[string]float64
)

type ExchangeRateHistory struct {
	Timestamp int64
	Rate	  float64
}

type EthereumClient struct {
	client *rpc.Client
}

func main() {
	exchangeRates = make(map[string]map[string]float64)
	exchangeRateHistory = make(map[string][]ExchangeRateHistory)

	go updateExchangeRates()

	r := gin.Default() // create gin.Engine 
	r.GET("/rates/:cryptocurrency/:fiat", getExchangeRate)
	r.GET("/rates/:cryptocurrency", getExchangeRates)
	r.GET("/rates", getAllExchangeRates)
	r.GET("/rates/history/:cryptocurrency/:fiat", getExchangeRateHistory)
	r.GET("/balance/:address", getBalanceHandler)

	r.Run(":8080")
}

func updateExchangeRates() {
	for {
		exchangeRatesMutex.Lock()
		fetchExchangeRates()
		exchangeRatesMutex.Unlock()

		// Update exchane rate history
		exchangeRateHistoryMutex.Lock()
		for cryptocurrency, rates := range exchangeRates {
			for fiat, rate := range rates {
				key := fmt.Sprintf("%s-%s", cryptocurrency, fiat)
				history, ok := exchangeRateHistory[key]
				if !ok {
					exchangeRateHistory[key] = []ExchangeRateHistory{}
				}

				exchangeRateHistory[key] = append(history, ExchangeRateHistory{
					Timestamp: time.Now().Unix(),
					Rate: rate,
				})
			}
		}
		exchangeRateHistoryMutex.Unlock()

		log.Println("Exchange rates updated")

		time.Sleep(updateInterval)
	}
}

func fetchExchangeRates() {
	client := resty.New()
	res, err := client.R().Get(coinGeckoAPIURL + "/simple/price?ids=bitcoin,ethereum,litecoin&vs_currencies=usd,eur,gbp")
	if err != nil {
		log.Println("Failed to fetch exchange rates:", err)
		return
	}
	if res.StatusCode() != 200 {
		log.Println("Failed to fetch exchange rates. Status code:", res.StatusCode())
		return
	}

	data := res.Body()

	// Define a struct to match the JSON structure of the response
	type ExchangeRateData struct {
		Bitcoin  map[string]float64 `json:"bitcoin"`
		Ethereum map[string]float64 `json:"ethereum"`
		Litecoin map[string]float64 `json:"litecoin"`
	}

	// Create an instance of the struct to hold the parsed data
	exchangeRateData := ExchangeRateData{}

	// Parse the JSON data into the struct
	err = json.Unmarshal(data, &exchangeRateData)
	if err != nil {
		log.Println("Failed to parse exchange rate data:", err)
		return
	}

	// Update the exchangeRates variable with the parsed data
	exchangeRates["bitcoin"] = exchangeRateData.Bitcoin
	exchangeRates["ethereum"] = exchangeRateData.Ethereum
	exchangeRates["litecoin"] = exchangeRateData.Litecoin


	// log.Println("Response body:", string(data))
	log.Println("Exchange rates updated:", exchangeRates)
}

func getExchangeRate(c *gin.Context) {
	cryptocurrency := c.Param("cryptocurrency")
	fiat := c.Param("fiat")

	// Check if the provided cryptocurrency is supported
	exchangeRate, ok := exchangeRates[cryptocurrency]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cryptocurrency not found"})
		return
	}

	// Check if the provided fiat currency is supported
	rate, ok := exchangeRate[fiat]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Fiat currency not found"})
		return
	}

	res := gin.H{
		"cryptocurrency": cryptocurrency,
		"fiat": fiat,
		"rate": rate,
	}

	c.JSON(http.StatusOK, res)
}

func getExchangeRates(c *gin.Context) {
	cryptocurrency := c.Param("cryptocurrency")

	// Check if the provided cryptocurrency is supported
	exchangeRate, ok := exchangeRates[cryptocurrency]

	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cryptocurrency not found"})
		return
	}

	res := gin.H{
		"cryptocurrency": cryptocurrency,
		"rates": exchangeRate,
	}

	c.JSON(http.StatusOK, res)
}

func getAllExchangeRates(c *gin.Context) {
	res := gin.H{
		"rates": exchangeRates,
	}

	c.JSON(http.StatusOK, res)
}

func getExchangeRateHistory(c *gin.Context) {
	cryptocurrency := c.Param("cryptocurrency")
	fiat := c.Param("fiat")

	// Check if the provided cryptocurrency is supported
	exchangeRate, ok := exchangeRates[cryptocurrency]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cryptocurrency not found"})
		return
	}

	// Check if the provided fiat currency is supported
	rate, ok := exchangeRate[fiat]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Fiat currency not found"})
		return
	}

	// Get the exchange rate history for the past 24 hours
	history, ok := exchangeRateHistory[fmt.Sprintf("%s-%s", cryptocurrency, fiat)]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exchange rate history not available"})
		return
	}

	// Filter the history for the past 24 hours
	currentTime := time.Now().Unix()
	filteredHistory := []ExchangeRateHistory{}
	for _, entry := range history {
		if entry.Timestamp > currentTime-24*60*60 {
			filteredHistory = append(filteredHistory, entry)
		}
	}

	res := gin.H{
		"cryptocurrency": cryptocurrency,
		"fiat": fiat,
		"rate": rate,
		"history": filteredHistory,
	}

	c.JSON(http.StatusOK, res)
}

// function to create Ethereum Client
func NewEthereumClient(endpoint string) (*EthereumClient, error) {
	client, err := rpc.Dial(endpoint)
	if err != nil {
		return nil, err
	}

	return &EthereumClient{
		client: client,
	}, nil
}

// GetBalance retrieves the current balance of the specified Ethereum address.
func (ec *EthereumClient) GetBalance(address string) (float64, error) {
	var result hexutil.Big
	err := ec.client.Call(&result, "eth_getBalance", common.HexToAddress(address), "latest")
	if err != nil {
		return 0, nil
	}

	balance := new(big.Float).SetInt(result.ToInt())
	ethValue := new(big.Float).Quo(balance, big.NewFloat(math.Pow10(18)))

	balanceFloat64, _ := ethValue.Float64() // Extract the float64 value

	return balanceFloat64, nil
}

func getBalanceHandler(c *gin.Context) {
	address := c.Param("address")

	url := "https://mainnet.infura.io/v3/b2ca6feb9d1a471d97c766f4a165d272"

	client, err := NewEthereumClient(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to connect to the Ethereum network",
		})
		return
	}

	balance, err := client.GetBalance(address)
	// balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve balance",
		})
		return
	}

	// Convert the balance to Ether value
	etherBalance := new(big.Float).Quo(new(big.Float).SetFloat64(balance), big.NewFloat(math.Pow10(18)))
	// etherBalance := new(big.Float).SetFloat64(balance)

	c.JSON(http.StatusOK, gin.H{
		"address": address,
		"balance": etherBalance,
	})
}

// address: 0x00000000219ab540356cbb839cbe05303d7705fa





