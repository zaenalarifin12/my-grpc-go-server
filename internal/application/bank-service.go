package application

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/zaenalarifin12/my-grpc-go-server/internal/adapter/database"
	"github.com/zaenalarifin12/my-grpc-go-server/internal/domain/bank"
	"github.com/zaenalarifin12/my-grpc-go-server/internal/port"
	"log"
	"time"
)

type BankService struct {
	db port.BankDatabasePort
}

func NewBankService(dbPort port.BankDatabasePort) *BankService {
	return &BankService{db: dbPort}
}
func (s *BankService) FindCurrentBalance(acct string) (float64, error) {
	bankAccount, err := s.db.GetBankAccountByNumber(acct)

	if err != nil {
		log.Println("Error on find current balance :", err)
		return 0, err
	}

	return bankAccount.CurrentBalance, nil
}

func (s *BankService) CreateExchangeRate(r bank.ExchangeRate) (uuid.UUID, error) {
	newUuid := uuid.New()
	now := time.Now()

	exchangeRateOrm := database.BankExcangeRateOrm{
		ExchangeRateUuid:   newUuid,
		FromCurrency:       r.FromCurrency,
		ToCurrency:         r.ToCurrency,
		Rate:               r.Rate,
		ValidFromTimestamp: r.ValidFromTimestamp,
		ValidToTimestamp:   r.ValidToTimestamp,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	return s.db.CreateExchangeRate(exchangeRateOrm)
}
func (s *BankService) FindExchangeRate(fromCur string, toCur string, ts time.Time) (float64, error) {
	exchangeRate, err := s.db.GetExchangeRateAtTimestamp(fromCur, toCur, ts)

	if err != nil {
		return 0, err
	}

	return float64(exchangeRate.Rate), nil
}

func (s *BankService) CreateTransaction(acct string, t bank.Transaction) (uuid.UUID, error) {
	newUuid := uuid.New()
	now := time.Now()

	bankAccountOrm, err := s.db.GetBankAccountByNumber(acct)

	if err != nil {
		log.Printf("Can't create transaction for %v : %v ", acct, err)
		return uuid.Nil, fmt.Errorf("can't find account number %v : %v", acct, err.Error())
	}

	if t.TransactionType == bank.TransactionTypeOut && bankAccountOrm.CurrentBalance < t.Amount {
		return bankAccountOrm.AccountUuid, fmt.Errorf("insufficient account balance %v for [out] transaction amount %v", bankAccountOrm.CurrentBalance, t.Amount)
	}

	transactionOrm := database.BankTransactionOrm{
		TransactionUuid:      newUuid,
		AccountUuid:          bankAccountOrm.AccountUuid,
		TransactionTimestamp: now,
		Amount:               t.Amount,
		TransactionType:      t.TransactionType,
		Notes:                t.Notes,
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	savedUuid, err := s.db.CreateTransaction(bankAccountOrm, transactionOrm)

	return savedUuid, nil
}

func (s *BankService) CalculateTransactionSummary(tcur *bank.TransactionSummary, trans bank.Transaction) error {
	switch trans.TransactionType {
	case bank.TransactionTypeIn:
		tcur.SumIn += trans.Amount
	case bank.TransactionTypeOut:
		tcur.SumOut += trans.Amount
	default:
		return fmt.Errorf("unknown transaction type %v", trans.TransactionType)
	}

	tcur.SumTotal = tcur.SumIn - tcur.SumOut

	return nil
}

func (s *BankService) Transfer(tt bank.TransferTransaction) (uuid.UUID, bool, error) {
	now := time.Now()

	fromAccountOrm, err := s.db.GetBankAccountByNumber(tt.FromAccountNumber)

	if err != nil {
		log.Printf("Can't find transfer from account %v : %v\n", tt.FromAccountNumber, err)
		return uuid.Nil, false, bank.ErrTransferSourceAccountNotFound
	}

	if fromAccountOrm.CurrentBalance < tt.Amount {
		return uuid.Nil, false, bank.ErrTransferTransactionPair
	}

	toAccountOrm, err := s.db.GetBankAccountByNumber(tt.ToAccountNumber)

	if err != nil {
		log.Printf("Can't find transfer to account %v : %v\n", tt.ToAccountNumber, err)
		return uuid.Nil, false, bank.ErrTransferDestinationAccountNotFound
	}

	fromTransactionOrm := database.BankTransactionOrm{
		TransactionUuid:      uuid.New(),
		AccountUuid:          fromAccountOrm.AccountUuid,
		TransactionTimestamp: now,
		Amount:               tt.Amount,
		TransactionType:      bank.TransactionTypeOut,
		Notes:                "Transfer out to " + tt.ToAccountNumber,
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	toTransactionOrm := database.BankTransactionOrm{
		TransactionUuid:      uuid.New(),
		AccountUuid:          toAccountOrm.AccountUuid,
		TransactionTimestamp: now,
		Amount:               tt.Amount,
		TransactionType:      bank.TransactionTypeOut,
		Notes:                "Transfer in from " + tt.FromAccountNumber,
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	// create transfer request
	newTransferUuid := uuid.New()

	transferOrm := database.BankTransferOrm{
		TransferUuid:      newTransferUuid,
		FromAccountUuid:   fromAccountOrm.AccountUuid,
		ToAccountUuid:     toAccountOrm.AccountUuid,
		Currency:          tt.Currency,
		Amount:            tt.Amount,
		TransferTimestamp: now,
		TransferSuccess:   false,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if _, err := s.db.CreateTransfer(transferOrm); err != nil {
		log.Printf("Can't create transfer from %v to %v : %v\n", tt.FromAccountNumber, tt.ToAccountNumber, err)
		return uuid.Nil, false, bank.ErrTransferRecordFailed
	}

	if transferPairSuccess, _ := s.db.CreateTransferTransactionPair(fromAccountOrm, toAccountOrm, fromTransactionOrm, toTransactionOrm); transferPairSuccess {
		s.db.UpdateTransferStatus(transferOrm, true)
		return newTransferUuid, true, nil
	} else {
		return newTransferUuid, false, bank.ErrTransferTransactionPair
	}
}
