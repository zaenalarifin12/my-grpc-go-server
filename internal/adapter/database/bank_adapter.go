package database

import (
	"github.com/google/uuid"
	"github.com/zaenalarifin12/my-grpc-go-server/internal/domain/bank"
	"log"
	"time"
)

func (a *DatabaseAdapter) GetBankAccountByNumber(acct string) (BankAccountOrm, error) {
	var bankAccountOrm BankAccountOrm

	if err := a.db.First(&bankAccountOrm, "account_number = ?", acct).Error; err != nil {
		log.Printf("Can't find bank account number %v : %v\n", acct, err)
		return bankAccountOrm, err
	}

	return bankAccountOrm, nil
}

func (a *DatabaseAdapter) CreateExchangeRate(r BankExcangeRateOrm) (uuid.UUID, error) {
	if err := a.db.Create(r).Error; err != nil {
		return uuid.Nil, err
	}

	return r.ExchangeRateUuid, nil
}

func (a *DatabaseAdapter) GetExchangeRateAtTimestamp(fromCur string, toCur string, ts time.Time) (BankExcangeRateOrm, error) {
	var exchangeRateOrm BankExcangeRateOrm

	err := a.db.First(&exchangeRateOrm, "from_currency = ? "+" AND to_currency = ? "+"AND (? BETWEEN valid_from_timestamp AND valid_TO_timestamp)", fromCur, toCur, ts).Error

	return exchangeRateOrm, err
}

func (a *DatabaseAdapter) CreateTransaction(acct BankAccountOrm, t BankTransactionOrm) (uuid.UUID, error) {
	tx := a.db.Begin()

	if err := a.db.Create(t).Error; err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	// recalculate current balance
	newAmount := t.Amount

	if t.TransactionType == bank.TransactionTypeOut {
		newAmount = -1 * t.Amount
	}

	newCurrentBalance := acct.CurrentBalance + newAmount

	if err := tx.Model(&acct).Updates(&map[string]interface{}{
		"current_balance": newCurrentBalance,
		"updated_at":      time.Now(),
	}).Error; err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	tx.Commit()

	return t.TransactionUuid, nil
}

func (a *DatabaseAdapter) CreateTransfer(transfer BankTransferOrm) (uuid.UUID, error) {
	if err := a.db.Create(transfer).Error; err != nil {
		return uuid.Nil, err
	}

	return transfer.TransferUuid, nil
}

func (a *DatabaseAdapter) CreateTransferTransactionPair(fromAccountOrm, toAccountOrm BankAccountOrm, fromTransactionOrm, toTransactionOrm BankTransactionOrm) (bool, error) {
	tx := a.db.Begin()

	if err := tx.Create(fromTransactionOrm).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	// recalculate current balance (fromAccount)
	fromAccountBalanceNew := fromAccountOrm.CurrentBalance - fromTransactionOrm.Amount

	if err := tx.Model(&fromAccountOrm).Updates(
		map[string]interface{}{
			"current_balance": fromAccountBalanceNew,
			"updated_at":      time.Now(),
		},
	).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	// recalculate current balance (toAccount)
	toAccountBalanceNew := toAccountOrm.CurrentBalance + toTransactionOrm.Amount

	if err := tx.Model(&toAccountOrm).Updates(
		map[string]interface{}{
			"current_balance": toAccountBalanceNew,
			"updated_at":      time.Now(),
		},
	).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	tx.Commit()

	return true, nil

}

func (a *DatabaseAdapter) UpdateTransferStatus(transfer BankTransferOrm, status bool) error {
	if err := a.db.Model(&transfer).Updates(
		map[string]interface{}{
			"transfer_success": status,
			"updated_at":       time.Now(),
		}).Error; err != nil {
		return err
	}

	return nil
}
