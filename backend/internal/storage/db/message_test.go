package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveMessage(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	store := &Store{db: db}

	mock.ExpectExec("INSERT INTO roommessages").
		WithArgs(1, "alice", "hello").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = store.SaveMessage(context.Background(), "hello", 1, "alice")
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetRecentMessages_NoLastID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	store := &Store{db: db}

	rows := sqlmock.NewRows([]string{"id", "room_id", "username", "content", "created_at"}).
		AddRow(2, 1, "bob", "world", time.Now()).
		AddRow(1, 1, "alice", "hello", time.Now())

	mock.ExpectQuery("SELECT .+ FROM roommessages").
		WithArgs(1, 20).
		WillReturnRows(rows)

	msgs, err := store.GetRecentMessages(context.Background(), 1, 20, 0)
	assert.NoError(t, err)
	assert.Len(t, msgs, 2)
	assert.Equal(t, "world", msgs[0].Content)
	assert.Equal(t, "hello", msgs[1].Content)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetRecentMessages_WithLastID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	store := &Store{db: db}

	rows := sqlmock.NewRows([]string{"id", "room_id", "username", "content", "created_at"}).
		AddRow(1, 1, "alice", "hello", time.Now())

	mock.ExpectQuery("SELECT .+ FROM roommessages").
		WithArgs(1, 20, 5).
		WillReturnRows(rows)

	msgs, err := store.GetRecentMessages(context.Background(), 1, 20, 5)
	assert.NoError(t, err)
	assert.Len(t, msgs, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSendDirectMessage(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	store := &Store{db: db}

	mock.ExpectExec("INSERT INTO directmessages").
		WithArgs(1, "alice", "bob", "hi bob").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = store.SendDirectMessage(context.Background(), "hi bob", 1, "alice", "bob")
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetRecentDirectMessages_NoLastID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	store := &Store{db: db}

	rows := sqlmock.NewRows([]string{"id", "conversation_id", "receiver", "sender", "content", "created_at"}).
		AddRow(2, 1, "bob", "alice", "msg2", time.Now()).
		AddRow(1, 1, "bob", "alice", "msg1", time.Now())

	mock.ExpectQuery("SELECT .+ FROM directmessages").
		WithArgs(1, 20).
		WillReturnRows(rows)

	msgs, err := store.GetRecentDirectMessages(context.Background(), 1, 20, 0)
	assert.NoError(t, err)
	assert.Len(t, msgs, 2)
	assert.Equal(t, 1, msgs[0].ConversationID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetRecentDirectMessages_WithLastID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	store := &Store{db: db}

	rows := sqlmock.NewRows([]string{"id", "conversation_id", "receiver", "sender", "content", "created_at"}).
		AddRow(1, 1, "bob", "alice", "msg1", time.Now())

	mock.ExpectQuery("SELECT .+ FROM directmessages").
		WithArgs(1, 20, 5).
		WillReturnRows(rows)

	msgs, err := store.GetRecentDirectMessages(context.Background(), 1, 20, 5)
	assert.NoError(t, err)
	assert.Len(t, msgs, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}
