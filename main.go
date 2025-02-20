package main

import (
	"encoding/json"
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
func getWallet(w http.ResponseWriter, r *http.Request) types.MyWallet {
	address := chi.URLParam(r, "address")
	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		Prefix:          "GetWallet ",
	})
	logger.Info("Scanning wallet", "address", address)

	wallet := requests.RequestAccountInfo(address)
	solPrice, _ := requests.GetSolPrice()
	accounts := requests.RequestTokenAccounts(address)
	var addresses []string

	for _, account := range accounts.Result.Value {
		addresses = append(addresses, account.Account.Data.Parsed.Info.Mint)
	}
	curr_prices := requests.GetCoinGeckoTokenPrices(addresses)
	priceHistories := make(map[string][][]float64)
	var tokens []types.MyToken
	walletValue := wallet.SolAmount * solPrice
	for _, account := range accounts.Result.Value {
		data := requests.GetTokenMetadata(account.Account.Data.Parsed.Info.Mint)
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
		token := types.MyToken{
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

	transactions := requests.GetTransactions(address)
	return types.MyWallet{
		Address:      address,
		Value:        walletValue,
		SolValue:     wallet.SolAmount * solPrice,
		SolBalance:   wallet.SolAmount,
		LastUpdated:  time.Now(),
		Tokens:       tokens,
		Transactions: transactions,
	}
}
