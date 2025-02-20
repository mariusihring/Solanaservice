package requests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sol_test/types"
	"strings"

	"github.com/charmbracelet/log"
)

func GetCoinGeckoTokenPrices(addresses []string) map[string]string {
	tokens := strings.Join(addresses, ",")
	request_url := fmt.Sprintf("https://api.geckoterminal.com/api/v2/simple/networks/solana/token_price/%s", tokens)
	resp, err := http.Get(request_url)
	if err != nil {
		log.Error("Error occured", "Stack", err)
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("Error occured", "Stack", err)
		return nil
	}
	var response types.CoinGeckoPriceResponse
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {

		log.Error("Error occured", "Stack", err)
		return nil
	}
	return response.Data.Attributes.TokenPrices
}

func GetTokenPools(address string) (string, error) {
	request_url := fmt.Sprintf("https://api.geckoterminal.com/api/v2/networks/solana/tokens/%s/pools?page=1", address)
	resp, err := http.Get(request_url)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var response types.CoinGeckoPoolResponse
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		return "", err
	}
	return response.Data[0].Attributes.Address, nil
}

func GetCoinGeckoOHLCVS(address string, timeframe string, start int64, end int64) ([][]float64, error) {
	request_url := fmt.Sprintf("https://api.geckoterminal.com/api/v2/networks/solana/pools/%s/ohlcv/%s?currency=usd", address, timeframe)
	resp, err := http.Get(request_url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {

		return nil, err
	}
	var response types.CoinGeckoOHLCVSResponse
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		log.Error("Error occured", "Stack", err)
		return nil, err
	}
	return response.Data.Attributes.OHLCVList, nil
}

// GetSolPrice retrieves the current USD price for SOL from CoinGecko.
func GetSolPrice() (float64, error) {
	url := "https://api.coingecko.com/api/v3/simple/price?ids=solana&vs_currencies=usd"
	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to get SOL price: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read SOL price response: %w", err)
	}

	var priceResp map[string]map[string]float64
	if err := json.Unmarshal(body, &priceResp); err != nil {
		return 0, fmt.Errorf("failed to unmarshal SOL price response: %w", err)
	}

	// Expected response: {"solana": {"usd": <price>}}
	price, ok := priceResp["solana"]["usd"]
	if !ok {
		return 0, fmt.Errorf("SOL price not found in response")
	}

	return price, nil
}
