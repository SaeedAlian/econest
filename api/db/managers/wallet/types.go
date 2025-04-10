package wallet_db_manager

import "time"

type TransactionType string

const (
	TransactionTypeDeposit  TransactionType = "deposit"
	TransactionTypeWithdraw TransactionType = "withdraw"
	TransactionTypePurchase TransactionType = "purchase"
	TransactionTypeSale     TransactionType = "sale"
)

var ValidTransactionTypes = []TransactionType{
	TransactionTypeDeposit,
	TransactionTypeWithdraw,
	TransactionTypePurchase,
	TransactionTypeSale,
}

func (t TransactionType) IsValid() bool {
	for _, v := range ValidTransactionTypes {
		if t == v {
			return true
		}
	}

	return false
}

type TransactionStatus string

const (
	TransactionStatusPending    TransactionStatus = "pending"
	TransactionStatusSuccessful TransactionStatus = "successful"
	TransactionStatusFailed     TransactionStatus = "failed"
)

var ValidTransactionStatuses = []TransactionStatus{
	TransactionStatusPending,
	TransactionStatusSuccessful,
	TransactionStatusFailed,
}

func (s TransactionStatus) IsValid() bool {
	for _, v := range ValidTransactionStatuses {
		if s == v {
			return true
		}
	}

	return false
}

type Wallet struct {
	Id        int
	Balance   float32
	CreatedAt time.Time
	UpdatedAt time.Time
}

type WalletTransaction struct {
	Id        int
	amount    float32
	TxType    TransactionType
	Status    TransactionStatus
	CreatedAt time.Time
	WalletId  int
}
