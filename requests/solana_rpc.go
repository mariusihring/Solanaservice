package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"sol_test/types"
	"strconv"
	"time"

	"github.com/charmbracelet/log"
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
func queryRPCWithRetry(method string, params []interface{}) (string, error) {
	requestPayload := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  method,
		"params":  params,
	}

	requestBytes, err := json.Marshal(requestPayload)
	if err != nil {
		log.Error("Error marshalling request", "Stack", err)
		return "", err
	}

	resp, err := http.Post(solanaRPC, "application/json", bytes.NewBuffer(requestBytes))
	if err != nil {
		log.Error("HTTP Post error", "Stack", err)
		return "", err
	}
	defer resp.Body.Close()

	// Check for rate limiting or other errors via status code.
	if resp.StatusCode != http.StatusOK {
		retryAfterHeader := resp.Header.Get("Retry-After")
		if retryAfterHeader != "" {
			if delaySec, err := strconv.Atoi(retryAfterHeader); err == nil {
				return "", fmt.Errorf("retry after %d seconds", delaySec)
			}
		}
		return "", fmt.Errorf("received non-200 status: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("Error reading response body", "Stack", err)
		return "", err
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, body, "", "  "); err != nil {
		log.Errorf("Error indenting JSON", err)
		return "", err
	}

	return prettyJSON.String(), nil
}

// GetTransactions fetches the transaction signatures, then uses a queue to ensure that all transaction data is fetched.
// It will respect the Retry-After header if the RPC returns a rate limiting response.
func GetTransactions(address string) []types.TransactionResponse {
	// First, get the signatures.
	data := queryRPC("getSignaturesForAddress", []interface{}{address})
	var sigResponse types.GetSignaturesForAddressResponse
	if err := json.Unmarshal([]byte(data), &sigResponse); err != nil {
		log.Error("Error unmarshalling getSignaturesForAddress response", "Stack", err)
		return nil
	}

	// Create a queue of signatures.
	var queue []string
	for _, sig := range sigResponse.Result {
		queue = append(queue, sig.Signature)
	}

	var transactions []types.TransactionResponse

	// Process the queue.
	for len(queue) > 0 {
		// Dequeue the first signature.
		signature := queue[0]
		queue = queue[1:]

		params := []interface{}{
			signature,
			map[string]interface{}{
				"encoding":                       "json",
				"maxSupportedTransactionVersion": 0,
			},
		}

		// Attempt to fetch the transaction data.
		result, err := queryRPCWithRetry("getTransaction", params)
		if err != nil {
			// If the error contains a retry delay, parse it.
			var delaySeconds int
			if n, _ := fmt.Sscanf(err.Error(), "retry after %d seconds", &delaySeconds); n == 1 {
				log.Info("Rate limited. Retrying after delay", "delaySeconds", delaySeconds, "signature", signature)
				time.Sleep(time.Duration(delaySeconds) * time.Second)
			} else {
				// For other errors, log and briefly wait before requeuing.
				log.Error("Error fetching transaction", "signature", signature, "error", err)
				time.Sleep(1 * time.Second)
			}
			// Requeue the signature for a retry.
			queue = append(queue, signature)
			continue
		}

		var txResponse types.TransactionResponse
		if err = json.Unmarshal([]byte(result), &txResponse); err != nil {
			log.Error("Error unmarshalling transaction", "signature", signature, "Stack", err)
			// Optionally requeue for another attempt.
			queue = append(queue, signature)
			continue
		}

		transactions = append(transactions, txResponse)
	}

	log.Info("Found Transactions", "wallet", address, "TransactionAmount", len(transactions))
	return transactions
}
