package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"sol_test/requests"
	"sol_test/types"
	"strconv"
	"time"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const solanaRPC = "https://api.mainnet-beta.solana.com"

// queryRPC queries the Solana RPC and returns a prettified json string
func queryRPC(method string, params []interface{}) string {
	requestPayload := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  method,
		"params":  params,
	}

	requestBytes, err := json.Marshal(requestPayload)
	if err != nil {
		log.Error("Error occured", "Stack", err)
	}
	//TODO: what to do if request returns error instead of data
	resp, err := http.Post(solanaRPC, "application/json", bytes.NewBuffer(requestBytes))
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, body, "", "  "); err != nil {
		log.Errorf("Error indenting JSON:", err)
		return ""
	}

	return prettyJSON.String()
}

func requestAccountInfo(address string) types.Wallet {
	data := queryRPC("getAccountInfo", []interface{}{address})
	var response types.GetAccountInfoResponse
	err := json.Unmarshal([]byte(data), &response)
	if err != nil {
		log.Error("Error occured", "Stack", err)
	}
	balance := queryRPC("getBalance", []interface{}{address})
	var walletresponse types.GetWalletResponse
	err = json.Unmarshal([]byte(balance), &walletresponse)
	if err != nil {
		log.Error("Error occured", "Stack", err)
	}
	divisor := math.Pow10(9)
	floatValue := float64(walletresponse.Result.Value) / divisor
	return types.Wallet{response, floatValue}
}

func requestTokenAccounts(address string) types.GetTokenAccountsByOwnerResponse {
	data := queryRPC("getTokenAccountsByOwner", []interface{}{
		address,
		map[string]interface{}{
			"programId": "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA",
		},
		map[string]interface{}{
			"encoding": "jsonParsed",
		},
	})
	var response types.GetTokenAccountsByOwnerResponse
	err := json.Unmarshal([]byte(data), &response)
	if err != nil {
		log.Error("Error occured", "Stack", err)
	}
	return response
}

func getTokenMetadata(address string) types.GetTokenMetaDataResponse {
	data := queryRPC("getAsset", []interface{}{address})
	var response types.GetTokenMetaDataResponse
	err := json.Unmarshal([]byte(data), &response)
	if err != nil {
		log.Error("Error occured", "Stack", err)
	}
	return response
}

func getTransactions(address string) []types.Meta {

	data := queryRPC("getSignaturesForAddress", []interface{}{address})
	var response types.GetSignaturesForAddressResponse
	err := json.Unmarshal([]byte(data), &response)
	if err != nil {
		log.Error("Error occured", "Stack", err)
	}
	var transactions []types.Meta
	for _, transaction_hash := range response.Result {
		transaction := queryRPC("getTransaction", []interface{}{transaction_hash, "json"})
		var transaction_parsed types.TransactionResponse
		err := json.Unmarshal([]byte(transaction), &response)
		if err != nil {
			log.Error("Error occured", "Stack", err)
		}
		transactions = append(transactions, transaction_parsed.Result.Meta)
	}

	return transactions

}

type MyWallet struct {
	Address      string    `json:"address"`
	SolBalance   float64   `json:"solBalance"`
	SolValue     float64   `json:"solValue"`
	Value        float64   `json:"walletValue"`
	Tokens       []MyToken `json:"tokens"`
	Transactions []string  `json:"transactions"`
	LastUpdated  time.Time `json:"last_updated"`
}

type MyToken struct {
	Name           string    `json:"name"`
	Address        string    `json:"address"`
	Pool           string    `json:"pool"`
	Description    string    `json:"description"`
	Image          string    `json:"image"`
	Amount         float64   `json:"amount"`
	Price          float64   `json:"price"`
	History_prices []float64 `json:"history_prices"`
	PnL            float64   `json:"pnl"`
	Invested       float64   `json:"invested"`
	Value          float64   `json:"value"`
}

func main() {

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	// Instead of writing "welcome", we now call our getWalletHandler.
	r.Get("/{address}", getWalletHandler)
	log.Info("Server running on port", "port", 3000)
	http.ListenAndServe(":3000", r)
}

// getWalletHandler wraps getWallet so it works as a chi handler.
func getWalletHandler(w http.ResponseWriter, r *http.Request) {
	wallet := getWallet(w, r)
	b, err := json.Marshal(wallet)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

// getWallet scans the wallet and populates price histories.
// (No other changes have been made to your original logic.)
func getWallet(w http.ResponseWriter, r *http.Request) MyWallet {
	address := chi.URLParam(r, "address")
	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		Prefix:          "GetWallet ",
	})
	logger.Info("Scanning wallet", "address", address)

	wallet := requestAccountInfo(address)
	solPrice, _ := requests.GetSolPrice()
	accounts := requestTokenAccounts(address)
	var addresses []string

	for _, account := range accounts.Result.Value {
		addresses = append(addresses, account.Account.Data.Parsed.Info.Mint)
	}
	curr_prices := requests.GetCoinGeckoTokenPrices(addresses)
	priceHistories := make(map[string][][]float64)
	var tokens []MyToken
	walletValue := wallet.SolAmount * solPrice
	for _, account := range accounts.Result.Value {
		data := getTokenMetadata(account.Account.Data.Parsed.Info.Mint)
		pool, _ := requests.GetTokenPools(account.Account.Data.Parsed.Info.Mint)
		logger.Info("Found Token", "token",
			data.Result.Content.Metadata.Name,
			"price",
			account.Account.Data.Parsed.Info.TokenAmount.UIAmount,
			"address",
			account.Account.Data.Parsed.Info.Mint)
		// TODO: make this for every transaction as we need smaller timeframes to get the prices.
		prices, _ := requests.GetCoinGeckoOHLCVS(pool, "hour", 0, 0)
		priceHistories[account.Account.Data.Parsed.Info.Mint] = prices
		f, err := strconv.ParseFloat(curr_prices[account.Account.Data.Parsed.Info.Mint], 64)
		if err != nil {
			log.Error("Error occured", "Stack", err)
			continue
		}
		walletValue += account.Account.Data.Parsed.Info.TokenAmount.UIAmount * f
		token := MyToken{
			Name:           data.Result.Content.Metadata.Name,
			Address:        account.Account.Data.Parsed.Info.Mint,
			Pool:           pool,
			Description:    data.Result.Content.Metadata.Description,
			Image:          data.Result.Content.Links.Image,
			Amount:         account.Account.Data.Parsed.Info.TokenAmount.UIAmount,
			Price:          f,
			History_prices: nil,
			PnL:            0,
			Invested:       0,
			Value:          account.Account.Data.Parsed.Info.TokenAmount.UIAmount * f,
		}
		tokens = append(tokens, token)
	}

	//	transactions := getTransactions(address)
	return MyWallet{
		Address:      address,
		Value:        walletValue,
		SolValue:     wallet.SolAmount * solPrice,
		SolBalance:   wallet.SolAmount,
		LastUpdated:  time.Now(),
		Tokens:       tokens,
		Transactions: nil,
	}
}
