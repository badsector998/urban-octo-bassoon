package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/badsector998/urban-octo-bassoon/util"
	"github.com/stretchr/testify/require"
)

func createNewEntry(t *testing.T) Entry {

	//create account needed because of the foreign constraints.
	//the user needs to be already existed in accounts table.
	accEntryParam := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	createEntryAcc, err := testQueries.CreateAccount(context.Background(), accEntryParam)
	require.NoError(t, err)
	require.NotEmpty(t, createEntryAcc)

	//create Entry
	args := NewEntryParams{
		AccountID: createEntryAcc.ID,
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.NewEntry(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, entry.AccountID, args.AccountID)
	require.Equal(t, entry.Amount, args.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestNewEntry(t *testing.T) {
	createNewEntry(t)
}

func TestGetEntry(t *testing.T) {
	entry := createNewEntry(t)

	entryExcute, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entryExcute)
}

func TestListEntries(t *testing.T) {
	for i := 0; i < 10; i++ {
		createNewEntry(t)
	}

	argsListEntry := ListEntriesParams{
		Limit:  5,
		Offset: 5,
	}

	entryExcute, err := testQueries.ListEntries(context.Background(), argsListEntry)
	require.NoError(t, err)
	require.Len(t, entryExcute, 5)

	for _, each_entry := range entryExcute {
		require.NotEmpty(t, each_entry)
	}
}

func TestDeleteEntry(t *testing.T) {
	newEntry := createNewEntry(t)

	err := testQueries.DeleteEntry(context.Background(), newEntry.ID)
	require.NoError(t, err)

	delEntry, err := testQueries.GetEntry(context.Background(), newEntry.ID)
	require.Error(t, err)
	require.Error(t, err, sql.ErrNoRows.Error())
	require.Empty(t, delEntry)
}
