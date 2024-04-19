package db

import (
	"context"
	"testing"

	"github.com/lordofthemind/backendMaster/utils"
	"github.com/stretchr/testify/require"
)

func CreateRandomEntry(t *testing.T, account Account) Entry {
	rg := utils.NewRandomGenerator()
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    rg.RandomMoney(),
	}
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	account := CreateRandomAccount(t)
	CreateRandomEntry(t, account)
}

func TestGetEntry(t *testing.T) {
	account := CreateRandomAccount(t)
	entry1 := CreateRandomEntry(t, account)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.NotZero(t, entry2.ID)
	require.NotZero(t, entry2.CreatedAt)
}

func TestListEntries(t *testing.T) {
	account := CreateRandomAccount(t)
	n := 10
	for i := 0; i < n; i++ {
		CreateRandomEntry(t, account)
	}
	arg := ListEntriesParams{
		AccountID: account.ID, // Use the ID of the created account
		Limit:     5,
		Offset:    5,
	}
	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)
	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
