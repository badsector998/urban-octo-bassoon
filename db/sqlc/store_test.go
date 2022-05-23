package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {

	store := NewStore(testDB)
	account1 := accountRandomTest(t)
	Account2 := accountRandomTest(t)
	fmt.Println(">> before : ", account1.Balance, Account2.Balance)

	n := 2
	amount := int64(10)

	errs := make(chan error)
	result_chan := make(chan TransferTxResult)

	// run n cocurent transfer transaction
	for i := 0; i < n; i++ {

		txName := fmt.Sprintf("Tx %d", i+1)

		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   Account2.ID,
				Amount:        amount,
			})

			errs <- err
			result_chan <- result

		}()
	}

	existed := make(map[int]bool)
	// check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-result_chan
		require.NotEmpty(t, result)

		// result store multiple objects.
		// check one by one

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, Account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		// check entry account1
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		// check entry account2
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, Account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, Account2.ID, toAccount.ID)

		// check balance
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - Account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)

		require.NotContains(t, existed, k)
		existed[k] = true

		// check final update balance
		updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
		require.NoError(t, err)

		updatedAccount2, err := store.GetAccount(context.Background(), Account2.ID)
		require.NoError(t, err)

		fmt.Println(">> after : ", updatedAccount1.Balance, updatedAccount2.Balance)

		require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
		require.Equal(t, Account2.Balance+int64(n)*amount, updatedAccount2.Balance)

	}

}
