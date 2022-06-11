package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(q *Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("err : %v, err rollback: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// type TransferTxParams struct {
// 	ToAccountID   int64 `json:"to_account_id"`
// 	FromAccountID int64 `json:"from_account_id"`
// 	Amount        int64 `json:"amount"`
// }

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	ToEntry     Entry    `json:"to_entry"`
	FromEntry   Entry    `json:"from_entry"`
	ToAccount   Account  `json:"to_account"`
	FromAccount Account  `json:"from_account"`
}

//Transfer money from one account to another
//Create transfer data
//Create entry data with positive value for ToAccountID and entry data with negative value for FromAccountID
func (store *Store) TransferTx(ctx context.Context, arg CreateTransferParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			ToAccountID:   arg.ToAccountID,
			FromAccountID: arg.FromAccountID,
			Amount:        arg.Amount,
		})

		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    -arg.Amount,
		})

		if err != nil {
			return err
		}

		//TODO: update Accounts' balance

		return err
	})
	return result, err
}
