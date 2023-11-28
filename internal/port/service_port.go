package port

import (
	"github.com/google/uuid"
	"github.com/zaenalarifin12/my-grpc-go-server/internal/domain/bank"
	"time"
)

type HelloServicePort interface {
	GenerateHello(name string) string
}

type BankServicePort interface {
	FindCurrentBalance(acc string) (float64, error)
	CreateExchangeRate(r bank.ExchangeRate) (uuid.UUID, error)
	FindExchangeRate(fromCur string, toCur string, ts time.Time) (float64, error)
	CreateTransaction(acct string, t bank.Transaction) (uuid.UUID, error)
	CalculateTransactionSummary(tcur *bank.TransactionSummary, trans bank.Transaction) error
	Transfer(tt bank.TransferTransaction) (uuid.UUID, bool, error)
}

type ResiliencyServicePort interface {
	GenerateResiliency(minDelaySecond int32, maxDelaySecond int32, statusCode []uint32) (string, uint32)
}
