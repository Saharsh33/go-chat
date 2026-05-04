package testutil

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockStore_UserOperations(t *testing.T) {
	store := NewMockStore()
	ctx := context.Background()

	store.CreateUserIfNotExists(ctx, "alice")
	store.CreateUserIfNotExists(ctx, "bob")

	// duplicate insert should not increase count
	store.CreateUserIfNotExists(ctx, "alice")
	assert.Len(t, store.Users, 2)

	id, err := store.GetUserByName(ctx, "alice")
	assert.NoError(t, err)
	assert.Greater(t, id, 0)

	_, err = store.GetUserByName(ctx, "unknown")
	assert.Error(t, err)
}

func TestMockStore_ConversationOperations(t *testing.T) {
	store := NewMockStore()
	ctx := context.Background()

	// create conversation
	c1, err := store.GetOrCreateConversation(ctx, "bob", "alice")
	require.NoError(t, err)
	// canonical: alice < bob
	assert.Equal(t, "alice", c1.UserOne)
	assert.Equal(t, "bob", c1.UserTwo)

	// get same conversation (reversed args)
	c2, err := store.GetOrCreateConversation(ctx, "alice", "bob")
	require.NoError(t, err)
	assert.Equal(t, c1.ID, c2.ID, "should return same conversation regardless of order")

	// only one conversation exists
	assert.Len(t, store.Conversations, 1)
}

func TestMockStore_DirectMessageWithConversation(t *testing.T) {
	store := NewMockStore()
	ctx := context.Background()

	convo, _ := store.GetOrCreateConversation(ctx, "alice", "bob")

	err := store.SendDirectMessage(ctx, "hello", convo.ID, "alice", "bob")
	assert.NoError(t, err)

	err = store.SendDirectMessage(ctx, "hi back", convo.ID, "bob", "alice")
	assert.NoError(t, err)

	msgs, err := store.GetRecentDirectMessages(ctx, convo.ID, 20, 0)
	assert.NoError(t, err)
	assert.Len(t, msgs, 2)
	// newest first
	assert.Equal(t, "hi back", msgs[0].Content)
	assert.Equal(t, "hello", msgs[1].Content)
}

func TestMockStore_DirectMessage_Pagination(t *testing.T) {
	store := NewMockStore()
	ctx := context.Background()

	convo, _ := store.GetOrCreateConversation(ctx, "alice", "bob")

	// Send 5 messages
	for i := 0; i < 5; i++ {
		store.SendDirectMessage(ctx, "msg", convo.ID, "alice", "bob")
	}

	// Get messages with ID < 4 (should return IDs 3, 2, 1)
	msgs, err := store.GetRecentDirectMessages(ctx, convo.ID, 20, 4)
	assert.NoError(t, err)
	assert.Len(t, msgs, 3)
}

func TestMockStore_ConversationsOfUser(t *testing.T) {
	store := NewMockStore()
	ctx := context.Background()

	store.GetOrCreateConversation(ctx, "alice", "bob")
	store.GetOrCreateConversation(ctx, "alice", "charlie")
	store.GetOrCreateConversation(ctx, "bob", "charlie") // alice not involved

	convos, err := store.GetConversationsOfUser(ctx, "alice")
	assert.NoError(t, err)
	assert.Len(t, convos, 2, "alice should have 2 conversations")
}

func TestMockStore_RoomOperations(t *testing.T) {
	store := NewMockStore()
	ctx := context.Background()

	room, err := store.CreateRoom(ctx, "general", "alice")
	require.NoError(t, err)
	assert.Equal(t, "general", room.Name)

	found, err := store.GetRoomByName(ctx, "general")
	assert.NoError(t, err)
	assert.Equal(t, room.ID, found.ID)

	found, err = store.GetRoomById(ctx, room.ID)
	assert.NoError(t, err)
	assert.Equal(t, "general", found.Name)
}

func TestMockStore_RoomMessages(t *testing.T) {
	store := NewMockStore()
	ctx := context.Background()

	store.SaveMessage(ctx, "hello", 1, "alice")
	store.SaveMessage(ctx, "world", 1, "bob")
	store.SaveMessage(ctx, "other room", 2, "charlie")

	msgs, err := store.GetRecentMessages(ctx, 1, 20, 0)
	assert.NoError(t, err)
	assert.Len(t, msgs, 2, "should only get room 1 messages")
}
