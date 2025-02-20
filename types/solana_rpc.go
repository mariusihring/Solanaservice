package types

type GetSignaturesForAddressResponse struct {
	JsonRPC string                          `json:"jsonrpc"`
	Result  []WalletTransactionHashResponse `json:"result"`
	Id      int64                           `json:"id"`
}

type WalletTransactionHashResponse struct {
	Err       string `json:"err"`
	Memo      string `json:"memo"`
	Signature string `json:"signature"`
	Slot      int64  `json:"slot"`
	BlockTime int64  `json:"blockTime"`
}

type Wallet struct {
	AccountInfo GetAccountInfoResponse
	SolAmount   float64
}

type GetWalletResponse struct {
	JsonRPC string          `json:"jsonrpc"`
	Result  GetWalletResult `json:"result"`
	Id      int64           `json:"id"`
}

type GetWalletResult struct {
	Context interface{} `json:"context"`
	Value   int64       `json:"value"`
}

type GetAccountInfoResponse struct {
	JsonRPC string               `json:"jsonrpc"`
	Result  GetAccountInfoResult `json:"result"`
	Id      int16                `json:"id"`
}

type GetAccountInfoResult struct {
	Context GetAccountInfoContext `json:"context"`
	Value   GetAccountInfoValue   `json:"value"`
}

type GetAccountInfoContext struct {
	ApiVersion string `json:"apiVersion"`
	Slot       int64  `json:"slot"`
}
type GetAccountInfoValue struct {
	Data       string `json:"data"`
	Executable bool   `json:"executable"`
	Lamports   int64  `json:"lamports"`
	Owner      string `json:"owner"`
	RentEpoch  uint64 `json:"rentEpoch"`
	Space      int64  `json:"space"`
}

type GetTokenAccountsByOwnerResponse struct {
	JsonRPC string                        `json:"jsonrpc"`
	Result  GetTokenAccountsByOwnerResult `json:"result"`
}

type GetTokenAccountsByOwnerResult struct {
	Context GetTokenAccountsByOwnerContext `json:"context"`
	Value   []TokenAccount                 `json:"value"`
}

type GetTokenAccountsByOwnerContext struct {
	ApiVersion string `json:"apiVersion"`
	Slot       int64  `json:"slot"`
}

type TokenAccount struct {
	Account Account `json:"account"`
	Pubkey  string  `json:"pubkey"`
}

type Account struct {
	Data       AccountData `json:"data"`
	Executable bool        `json:"executable"`
	Lamports   int64       `json:"lamports"`
	Owner      string      `json:"owner"`
	RentEpoch  uint64      `json:"rentEpoch"`
	Space      int         `json:"space"`
}

type AccountData struct {
	Parsed  ParsedData `json:"parsed"`
	Program string     `json:"program"`
	Space   int        `json:"space"`
}

type ParsedData struct {
	Info AccountInfo `json:"info"`
	Type string      `json:"type"`
}

type AccountInfo struct {
	IsNative    bool        `json:"isNative"`
	Mint        string      `json:"mint"`
	Owner       string      `json:"owner"`
	State       string      `json:"state"`
	TokenAmount TokenAmount `json:"tokenAmount"`
}

type TokenAmount struct {
	Amount         string  `json:"amount"`
	Decimals       int     `json:"decimals"`
	UIAmount       float64 `json:"uiAmount"`
	UIAmountString string  `json:"uiAmountString"`
}
type GetTokenMetaDataResponse struct {
	JSONRPC string        `json:"jsonrpc"`
	Result  TokenMetaData `json:"result"`
	ID      int           `json:"id"`
}

type TokenMetaData struct {
	Interface   string        `json:"interface"`
	ID          string        `json:"id"`
	Content     Content       `json:"content"`
	Authorities []Authority   `json:"authorities"`
	Compression Compression   `json:"compression"`
	Grouping    []interface{} `json:"grouping"`
	Royalty     Royalty       `json:"royalty"`
	Creators    []interface{} `json:"creators"`
	Ownership   Ownership     `json:"ownership"`
	Mutable     bool          `json:"mutable"`
	Burnt       bool          `json:"burnt"`
}

type Content struct {
	Schema   string   `json:"$schema"`
	JSONURI  string   `json:"json_uri"`
	Files    []File   `json:"files"`
	Metadata Metadata `json:"metadata"`
	Links    Links    `json:"links"`
}

type File struct {
	URI  string `json:"uri"`
	Mime string `json:"mime"`
}

type Metadata struct {
	Description   string `json:"description"`
	Name          string `json:"name"`
	Symbol        string `json:"symbol"`
	TokenStandard string `json:"token_standard"`
}

type Links struct {
	Image string `json:"image"`
}

type Authority struct {
	Address string   `json:"address"`
	Scopes  []string `json:"scopes"`
}

type Compression struct {
	Eligible    bool   `json:"eligible"`
	Compressed  bool   `json:"compressed"`
	DataHash    string `json:"data_hash"`
	CreatorHash string `json:"creator_hash"`
	AssetHash   string `json:"asset_hash"`
	Tree        string `json:"tree"`
	Seq         int    `json:"seq"`
	LeafID      int    `json:"leaf_id"`
}

type Royalty struct {
	RoyaltyModel        string  `json:"royalty_model"`
	Target              *string `json:"target"`
	Percent             float64 `json:"percent"`
	BasisPoints         int     `json:"basis_points"`
	PrimarySaleHappened bool    `json:"primary_sale_happened"`
	Locked              bool    `json:"locked"`
}

type Ownership struct {
	Frozen         bool    `json:"frozen"`
	Delegated      bool    `json:"delegated"`
	Delegate       *string `json:"delegate"`
	OwnershipModel string  `json:"ownership_model"`
	Owner          string  `json:"owner"`
}

// TransactionResponse is the top-level response type.
type TransactionResponse struct {
	JsonRPC string             `json:"jsonrpc"`
	Result  *TransactionResult `json:"result"`
	Error   *SolanaError       `json:"error"`
	ID      int                `json:"id"`
}

type SolanaError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Result contains the main result fields.
type TransactionResult struct {
	BlockTime   int64       `json:"blockTime"`
	Meta        Meta        `json:"meta"`
	Slot        int         `json:"slot"`
	Transaction Transaction `json:"transaction"`
	Version     string      `json:"version"`
}

// Meta holds metadata about the transaction.
type Meta struct {
	ComputeUnitsConsumed int                `json:"computeUnitsConsumed"`
	Err                  interface{}        `json:"err"` // Can be further refined if error details are known.
	Fee                  int                `json:"fee"`
	InnerInstructions    []InnerInstruction `json:"innerInstructions"`
	LoadedAddresses      LoadedAddresses    `json:"loadedAddresses"`
	LogMessages          []string           `json:"logMessages"`
	PostBalances         []int64            `json:"postBalances"`
	PostTokenBalances    []TokenBalance     `json:"postTokenBalances"`
	PreBalances          []int64            `json:"preBalances"`
	PreTokenBalances     []TokenBalance     `json:"preTokenBalances"`
	Rewards              []interface{}      `json:"rewards"` // Adjust type if rewards have a defined structure.
	Status               Status             `json:"status"`
}

// InnerInstruction groups a set of instructions with an index.
type InnerInstruction struct {
	Index        int           `json:"index"`
	Instructions []Instruction `json:"instructions"`
}

// Instruction represents a single instruction.
type Instruction struct {
	Accounts       []int  `json:"accounts"`
	Data           string `json:"data"`
	ProgramIdIndex int    `json:"programIdIndex"`
	// stackHeight can be null so we use a pointer.
	StackHeight *int `json:"stackHeight"`
}

// LoadedAddresses lists read-only and writable addresses.
type LoadedAddresses struct {
	Readonly []string `json:"readonly"`
	Writable []string `json:"writable"`
}

// TokenBalance represents the balance details for a token account.
type TokenBalance struct {
	AccountIndex  int                    `json:"accountIndex"`
	Mint          string                 `json:"mint"`
	Owner         string                 `json:"owner"`
	ProgramId     string                 `json:"programId"`
	UiTokenAmount TransactionTokenAmount `json:"uiTokenAmount"`
}

// TokenAmount holds amount details in different formats.
type TransactionTokenAmount struct {
	Amount         string  `json:"amount"`
	Decimals       int     `json:"decimals"`
	UiAmount       float64 `json:"uiAmount"`
	UiAmountString string  `json:"uiAmountString"`
}

// Status represents the transaction status.
type Status struct {
	Ok interface{} `json:"Ok"`
}

// Transaction contains the signed transaction data.
type Transaction struct {
	Message    Message  `json:"message"`
	Signatures []string `json:"signatures"`
}

// Message contains details of the transaction message.
type Message struct {
	AccountKeys         []string             `json:"accountKeys"`
	AddressTableLookups []AddressTableLookup `json:"addressTableLookups"`
	Header              Header               `json:"header"`
	Instructions        []MessageInstruction `json:"instructions"`
	RecentBlockhash     string               `json:"recentBlockhash"`
}

// AddressTableLookup represents an address table lookup.
type AddressTableLookup struct {
	AccountKey      string `json:"accountKey"`
	ReadonlyIndexes []int  `json:"readonlyIndexes"`
	WritableIndexes []int  `json:"writableIndexes"`
}

// Header defines the header information for the transaction message.
type Header struct {
	NumReadonlySignedAccounts   int `json:"numReadonlySignedAccounts"`
	NumReadonlyUnsignedAccounts int `json:"numReadonlyUnsignedAccounts"`
	NumRequiredSignatures       int `json:"numRequiredSignatures"`
}

// MessageInstruction represents an instruction within the transaction message.
type MessageInstruction struct {
	Accounts       []int  `json:"accounts"`
	Data           string `json:"data"`
	ProgramIdIndex int    `json:"programIdIndex"`
	// stackHeight can be null so we use a pointer.
	StackHeight *int `json:"stackHeight"`
}
