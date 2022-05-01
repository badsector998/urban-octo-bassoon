package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/badsector998/urban-octo-bassoon/util"
	"github.com/stretchr/testify/require"
)

func createNewTransfer(t *testing.T) Transfer {
	//create two accounts as from and to.
	accTrfArg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	fromAccount, err := testQueries.CreateAccount(context.Background(), accTrfArg)
	require.NoError(t, err)
	require.NotEmpty(t, fromAccount)
	require.NotZero(t, fromAccount.ID)
	require.NotZero(t, fromAccount.CreatedAt)

	toAccount, err := testQueries.CreateAccount(context.Background(), accTrfArg)
	require.NoError(t, err)
	require.NotEmpty(t, toAccount)
	require.NotZero(t, toAccount.ID)
	require.NotZero(t, toAccount.CreatedAt)

	//create transfer
	txArgs := NewTransferParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        util.RandomMoney(),
	}

	newTx, err := testQueries.NewTransfer(context.Background(), txArgs)
	require.NoError(t, err)
	require.NotEmpty(t, newTx)
	require.NotEqual(t, newTx.FromAccountID, newTx.ToAccountID)

	return newTx
}

func TestNewTransfer(t *testing.T) {
	createNewTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	tx1 := createNewTransfer(t)

	tx2, err := testQueries.GetTransfer(context.Background(), tx1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, tx2)
	require.Equal(t, tx2.ID, tx1.ID)
	require.WithinDuration(t, tx1.CreatedAt, tx2.CreatedAt, time.Second)
}

func TestListTransfers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createNewTransfer(t)
	}

	args := ListTransfersParams{
		Limit:  5,
		Offset: 5,
	}

	txList, err := testQueries.ListTransfers(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, txList, 5)

	for _, listTx := range txList {
		require.NotEmpty(t, listTx)
	}
}

func TestDeleteTransfer(t *testing.T) {
	tx1 := createNewTransfer(t)

	err := testQueries.DeleteTransfer(context.Background(), tx1.ID)
	require.NoError(t, err)

	tx2, err := testQueries.GetTransfer(context.Background(), tx1.ID)
	require.Error(t, err)
	require.Error(t, err, sql.ErrNoRows.Error())
	require.Empty(t, tx2)
}
