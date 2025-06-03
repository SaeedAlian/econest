package types

import "time"

// Wallet represents a user's digital wallet
// @model Wallet
type Wallet struct {
	// Wallet ID (private, needs permission)
	Id int `json:"id"        exposure:"private,needPermission"`
	// Current balance in the wallet (private, needs permission)
	Balance float64 `json:"balance"   exposure:"private,needPermission"`
	// When the wallet was created (private, needs permission)
	CreatedAt time.Time `json:"createdAt" exposure:"private,needPermission"`
	// When the wallet was last updated (private, needs permission)
	UpdatedAt time.Time `json:"updatedAt" exposure:"private,needPermission"`
	// ID of the user who owns this wallet (private, needs permission)
	UserId int `json:"userId"    exposure:"private,needPermission"`
}

// WalletTransaction represents a transaction in a user's wallet
// @model WalletTransaction
type WalletTransaction struct {
	// Transaction ID (private, needs permission)
	Id int `json:"id"        exposure:"private,needPermission"`
	// Transaction amount (private, needs permission)
	Amount float64 `json:"amount"    exposure:"private,needPermission"`
	// Type of transaction (credit/debit/etc) (private, needs permission)
	TxType TransactionType `json:"txType"    exposure:"private,needPermission"`
	// Current status of the transaction (private, needs permission)
	Status TransactionStatus `json:"status"    exposure:"private,needPermission"`
	// When the transaction was created (private, needs permission)
	CreatedAt time.Time `json:"createdAt" exposure:"private,needPermission"`
	// When the transaction was last updated (private, needs permission)
	UpdatedAt time.Time `json:"updatedAt" exposure:"private,needPermission"`
	// ID of the wallet this transaction belongs to (private, needs permission)
	WalletId int `json:"walletId"  exposure:"private,needPermission"`
}

// CreateWalletTransactionPayload contains data needed to create a new wallet transaction
// @model CreateWalletTransactionPayload
type CreateWalletTransactionPayload struct {
	// Transaction amount (required)
	Amount float64 `json:"amount"   validate:"required"`
	// Type of transaction
	TxType TransactionType `json:"txType"`
	// ID of the wallet this transaction belongs to
	WalletId int `json:"walletId"`
}

// UpdateWalletTransactionPayload contains data for updating a wallet transaction
// @model UpdateWalletTransactionPayload
type UpdateWalletTransactionPayload struct {
	// New status for the transaction
	Status *TransactionStatus `json:"status"`
}

// UpdateWalletPayload contains data for updating a wallet
// @model UpdateWalletPayload
type UpdateWalletPayload struct {
	// New balance for the wallet
	Balance *float64 `json:"balance"`
}

// WalletTransactionSearchQuery contains parameters for searching wallet transactions
// @model WalletTransactionSearchQuery
type WalletTransactionSearchQuery struct {
	// Filter by transaction status
	Status *TransactionStatus `json:"status"`
	// Filter by transaction type
	TxType *TransactionType `json:"txType"`
	// Filter transactions before this date
	BeforeDate *time.Time `json:"beforeDate"`
	// Filter transactions after this date
	AfterDate *time.Time `json:"afterDate"`
	// Filter by user ID
	UserId *int `json:"userId"`
	// Maximum number of results to return
	Limit *int `json:"limit"`
	// Number of results to skip
	Offset *int `json:"offset"`
}
