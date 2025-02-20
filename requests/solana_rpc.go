package requests

import (
	"bytes"
	"encoding/json"
	"github.com/charmbracelet/log"
	"io/ioutil"
	"math"
	"net/http"
	"sol_test/types"
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

func RequestAccountInfo(address string) types.Wallet {
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

func RequestTokenAccounts(address string) types.GetTokenAccountsByOwnerResponse {
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

func GetTokenMetadata(address string) types.GetTokenMetaDataResponse {
	data := queryRPC("getAsset", []interface{}{address})
	var response types.GetTokenMetaDataResponse
	err := json.Unmarshal([]byte(data), &response)
	if err != nil {
		log.Error("Error occured", "Stack", err)
	}
	return response
}

func GetTransactions(address string) []types.TransactionResult {
	data := queryRPC("getSignaturesForAddress", []interface{}{address})
	var response types.GetSignaturesForAddressResponse
	err := json.Unmarshal([]byte(data), &response)
	if err != nil {
		log.Error("Error occured", "Stack", err)
	}
	var transactions []types.TransactionResult
	for _, transaction_hash := range response.Result {
		transaction := queryRPC("getTransaction", []interface{}{transaction_hash.Signature, "json"})
		var transaction_parsed types.TransactionResponse
		err := json.Unmarshal([]byte(transaction), &transaction_parsed)
		if err != nil {
			log.Error("Error occured", "Stack", err)
		}
		transactions = append(transactions, transaction_parsed.Result)
	}

	log.Info("Found Transactions", "wallet", address, "TransactionAmount", len(transactions))
	return transactions

}
