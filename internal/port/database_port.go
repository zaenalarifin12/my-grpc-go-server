package port

import (
	"github.com/google/uuid"
	"github.com/zaenalarifin12/my-grpc-go-server/internal/adapter/database"
	"time"
)

type DummyDatabasePort interface {
	Save(data *database.DummyOrm) (uuid.UUID, error)
	GetByUUID(uuid uuid.UUID) (*database.DummyOrm, error)
}

type BankDatabasePort interface {
	GetBankAccountByNumber(acc string) (database.BankAccountOrm, error)
	CreateExchangeRate(r database.BankExcangeRateOrm) (uuid.UUID, error)
	GetExchangeRateAtTimestamp(fromCur string, toCur string, ts time.Time) (database.BankExcangeRateOrm, error)
	CreateTransaction(acct database.BankAccountOrm, t database.BankTransactionOrm) (uuid.UUID, error)
	CreateTransfer(transfer database.BankTransferOrm) (uuid.UUID, error)
	CreateTransferTransactionPair(fromAccountOrm, toAccountOrm database.BankAccountOrm, fromTransactionOrm, toTransactionOrm database.BankTransactionOrm) (bool, error)
	UpdateTransferStatus(transfer database.BankTransferOrm, status bool) error
}
