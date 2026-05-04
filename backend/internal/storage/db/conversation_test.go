package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetOrCreateConversation(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	store := &Store{db: db}

	now := time.Now()
	mock.ExpectQuery("INSERT INTO conversations").
		WithArgs("alice", "bob").
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_one", "user_two", "created_at", "updated_at"}).
			AddRow(1, "alice", "bob", now, now))

	convo, err := store.GetOrCreateConversation(context.Background(), "alice", "bob")
	assert.NoError(t, err)
	require.NotNil(t, convo)
	assert.Equal(t, 1, convo.ID)
	assert.Equal(t, "alice", convo.UserOne)
	assert.Equal(t, "bob", convo.UserTwo)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetConversationsOfUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	store := &Store{db: db}

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "user_one", "user_two", "created_at", "updated_at"}).
		AddRow(1, "alice", "bob", now, now).
		AddRow(2, "alice", "charlie", now, now)

	mock.ExpectQuery("SELECT .+ FROM conversations").
		WithArgs("alice").
		WillReturnRows(rows)

	convos, err := store.GetConversationsOfUser(context.Background(), "alice")
	assert.NoError(t, err)
	assert.Len(t, convos, 2)
	assert.Equal(t, "bob", convos[0].UserTwo)
	assert.Equal(t, "charlie", convos[1].UserTwo)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetConversationsOfUser_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	store := &Store{db: db}

	mock.ExpectQuery("SELECT .+ FROM conversations").
		WithArgs("newuser").
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_one", "user_two", "created_at", "updated_at"}))

	convos, err := store.GetConversationsOfUser(context.Background(), "newuser")
	assert.NoError(t, err)
	assert.Nil(t, convos)
	assert.NoError(t, mock.ExpectationsWereMet())
}
