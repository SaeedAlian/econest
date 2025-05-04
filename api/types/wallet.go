package types

import "time"

type Wallet struct {
	Id        int       `json:"id"        exposure:"private,needPermission"`
	Balance   float64   `json:"balance"   exposure:"private,needPermission"`
	CreatedAt time.Time `json:"createdAt" exposure:"private,needPermission"`
	UpdatedAt time.Time `json:"updatedAt" exposure:"private,needPermission"`
	UserId    int       `json:"userId"    exposure:"private,needPermission"`
}

type WalletTransaction struct {
	Id        int               `json:"id"        exposure:"private,needPermission"`
	Amount    float64           `json:"amount"    exposure:"private,needPermission"`
	TxType    TransactionType   `json:"txType"    exposure:"private,needPermission"`
	Status    TransactionStatus `json:"status"    exposure:"private,needPermission"`
	CreatedAt time.Time         `json:"createdAt" exposure:"private,needPermission"`
	WalletId  int               `json:"walletId"  exposure:"private,needPermission"`
}

type CreateWalletTransactionPayload struct {
	Amount   float64         `json:"amount"   validate:"required"`
	TxType   TransactionType `json:"txType"   validate:"required"`
	WalletId int             `json:"walletId" validate:"required"`
}

type UpdateWalletTransactionPayload struct {
	Status *TransactionStatus `json:"status"`
}

type UpdateWalletPayload struct {
	Balance *float64 `json:"balance"`
}

type WalletTransactionSearchQuery struct {
	Status     *TransactionStatus `json:"status"`
	TxType     *TransactionType   `json:"txType"`
	BeforeDate *time.Time         `json:"beforeDate"`
	AfterDate  *time.Time         `json:"afterDate"`
	UserId     *int               `json:"userId"`
	Limit      *int               `json:"limit"`
	Offset     *int               `json:"offset"`
}
