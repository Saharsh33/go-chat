package websocket

import (
	"chat-server/internal/testutil"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// helper: create a hub with mock store, start it, and return both
func setupHub(t *testing.T) (*Hub, *testutil.MockStore) {
	t.Helper()
	store := testutil.NewMockStore()
	hub := NewHub(store)
	go hub.Run()
	return hub, store
}

// helper: create a fake client and register it with the hub
func registerClient(t *testing.T, hub *Hub, username string) *Client {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	client := &Client{
		Conn:     nil, // not needed for hub tests
		Username: username,
		Send:     make(chan Message, 100), // buffered so tests don't block
		Ctx:      ctx,
		Cancel:   cancel,
	}
	hub.Register <- client
	// Wait for the registration to be processed
	time.Sleep(50 * time.Millisecond)
	return client
}

// drain all pending messages from client channel
func drainMessages(client *Client) []Message {
	var msgs []Message
	for {
		select {
		case m, ok := <-client.Send:
			if !ok {
				return msgs // channel closed
			}
			msgs = append(msgs, m)
		default:
			return msgs
		}
	}
}

// hasMsg checks if the messages contain a specific type+content
func hasMsg(msgs []Message, typ MessageType, content string) bool {
	for _, m := range msgs {
		if m.Type == typ && m.Content == content {
			return true
		}
	}
	return false
}

// hasMsgType checks if the messages contain a specific type
func hasMsgType(msgs []Message, typ MessageType) bool {
	for _, m := range msgs {
		if m.Type == typ {
			return true
		}
	}
	return false
}

// ===================== Tests =====================

func TestNewHub(t *testing.T) {
	store := testutil.NewMockStore()
	hub := NewHub(store)

	assert.NotNil(t, hub)
	assert.NotNil(t, hub.Clients)
	assert.NotNil(t, hub.UsersOfRoom)
	assert.NotNil(t, hub.RoomsOfUser)
	assert.NotNil(t, hub.JoinRoom)
	assert.NotNil(t, hub.LeaveRoom)
	assert.NotNil(t, hub.CreateRoom)
	assert.NotNil(t, hub.Register)
	assert.NotNil(t, hub.Unregister)
	assert.NotNil(t, hub.SendMessage)
}

func TestRegisterClient(t *testing.T) {
	hub, store := setupHub(t)

	client := registerClient(t, hub, "alice")

	// Verify user was created in store
	_, err := store.GetUserByName(context.Background(), "alice")
	assert.NoError(t, err)

	// Should receive a system registration message
	msgs := drainMessages(client)
	assert.True(t, hasMsgType(msgs, MsgSystem), "expected system registration message")
}

func TestUnregisterClient(t *testing.T) {
	hub, _ := setupHub(t)

	client := registerClient(t, hub, "bob")
	drainMessages(client) // discard registration msgs

	hub.Unregister <- client
	time.Sleep(50 * time.Millisecond)

	// The hub sends a system message then closes the channel.
	// Drain any remaining messages first.
	drainMessages(client)

	// Now the channel should be closed — next receive returns zero value + false.
	_, ok := <-client.Send
	assert.False(t, ok, "Send channel should be closed after unregister")
}

func TestCreateRoom(t *testing.T) {
	hub, store := setupHub(t)

	client := registerClient(t, hub, "alice")
	drainMessages(client)

	hub.CreateRoom <- RoomOps{clientDetails: client, roomName: "general"}
	time.Sleep(50 * time.Millisecond)

	// Verify room was created in store
	room, err := store.GetRoomByName(context.Background(), "general")
	require.NoError(t, err)
	assert.Equal(t, "general", room.Name)

	// Should receive a system confirmation message
	msgs := drainMessages(client)
	assert.True(t, hasMsgType(msgs, MsgSystem), "expected system message after room creation")
}

func TestJoinRoom(t *testing.T) {
	hub, store := setupHub(t)

	// alice creates a room
	alice := registerClient(t, hub, "alice")
	drainMessages(alice)

	hub.CreateRoom <- RoomOps{clientDetails: alice, roomName: "dev"}
	time.Sleep(50 * time.Millisecond)
	drainMessages(alice)

	// Get room ID from store (thread-safe)
	room, err := store.GetRoomByName(context.Background(), "dev")
	require.NoError(t, err)
	roomID := room.ID

	// bob joins the room
	bob := registerClient(t, hub, "bob")
	drainMessages(bob)

	hub.JoinRoom <- RoomOps{clientDetails: bob, roomDetails: roomID}
	time.Sleep(50 * time.Millisecond)

	// Verify via system message that join was successful
	msgs := drainMessages(bob)
	assert.True(t, hasMsgType(msgs, MsgSystem), "bob should receive join confirmation")

	// Verify bob can now send a message to the room (functional test)
	hub.SendMessage <- Message{Type: MsgRoomMessage, User: "bob", Room: roomID, Content: "hi from bob"}
	time.Sleep(50 * time.Millisecond)

	assert.Equal(t, 1, store.RoomMessageCount(), "message from bob should be saved")
}

func TestLeaveRoom(t *testing.T) {
	hub, store := setupHub(t)

	alice := registerClient(t, hub, "alice")
	drainMessages(alice)

	hub.CreateRoom <- RoomOps{clientDetails: alice, roomName: "temp"}
	time.Sleep(50 * time.Millisecond)
	drainMessages(alice)

	room, _ := store.GetRoomByName(context.Background(), "temp")
	roomID := room.ID

	bob := registerClient(t, hub, "bob")
	drainMessages(bob)

	hub.JoinRoom <- RoomOps{clientDetails: bob, roomDetails: roomID}
	time.Sleep(50 * time.Millisecond)
	drainMessages(bob)

	// bob leaves
	hub.LeaveRoom <- RoomOps{clientDetails: bob, roomDetails: roomID}
	time.Sleep(50 * time.Millisecond)

	// Verify via system message
	msgs := drainMessages(bob)
	assert.True(t, hasMsgType(msgs, MsgSystem), "bob should receive leave confirmation")

	// Verify bob can NO LONGER send a message to the room
	drainMessages(bob)
	hub.SendMessage <- Message{Type: MsgRoomMessage, User: "bob", Room: roomID, Content: "should fail"}
	time.Sleep(50 * time.Millisecond)

	bobMsgs := drainMessages(bob)
	assert.True(t, hasMsg(bobMsgs, MsgSystem, "User not part of room"), "bob should be rejected")
}

func TestRoomMessage(t *testing.T) {
	hub, store := setupHub(t)

	alice := registerClient(t, hub, "alice")
	drainMessages(alice)

	hub.CreateRoom <- RoomOps{clientDetails: alice, roomName: "chat"}
	time.Sleep(50 * time.Millisecond)
	drainMessages(alice)

	room, _ := store.GetRoomByName(context.Background(), "chat")
	roomID := room.ID

	bob := registerClient(t, hub, "bob")
	drainMessages(bob)

	hub.JoinRoom <- RoomOps{clientDetails: bob, roomDetails: roomID}
	time.Sleep(50 * time.Millisecond)
	drainMessages(bob)
	drainMessages(alice)

	// alice sends a room message
	hub.SendMessage <- Message{Type: MsgRoomMessage, User: "alice", Room: roomID, Content: "hello room!"}
	time.Sleep(50 * time.Millisecond)

	// bob should receive the message
	bobMsgs := drainMessages(bob)
	assert.True(t, hasMsg(bobMsgs, MsgRoomMessage, "hello room!"), "bob should receive alice's room message")

	// message should be saved in store
	roomMsgs := store.GetRoomMessagesCopy()
	assert.Len(t, roomMsgs, 1)
	assert.Equal(t, "hello room!", roomMsgs[0].Content)
}

func TestRoomMessageNonMember(t *testing.T) {
	hub, store := setupHub(t)

	alice := registerClient(t, hub, "alice")
	drainMessages(alice)

	hub.CreateRoom <- RoomOps{clientDetails: alice, roomName: "private"}
	time.Sleep(50 * time.Millisecond)
	drainMessages(alice)

	room, _ := store.GetRoomByName(context.Background(), "private")
	roomID := room.ID

	// bob is NOT in the room
	bob := registerClient(t, hub, "bob")
	drainMessages(bob)

	// bob tries to send message to the room
	hub.SendMessage <- Message{Type: MsgRoomMessage, User: "bob", Room: roomID, Content: "sneaky msg"}
	time.Sleep(50 * time.Millisecond)

	// Message should NOT be saved
	assert.Equal(t, 0, store.RoomMessageCount())

	// bob should receive "User not part of room" system message
	msgs := drainMessages(bob)
	assert.True(t, hasMsg(msgs, MsgSystem, "User not part of room"), "bob should get 'User not part of room' error")
}

func TestDirectMessage_CreatesConversation(t *testing.T) {
	hub, store := setupHub(t)

	alice := registerClient(t, hub, "alice")
	bob := registerClient(t, hub, "bob")
	drainMessages(alice)
	drainMessages(bob)

	hub.SendMessage <- Message{Type: MsgDirectMessage, User: "alice", Receiver: "bob", Content: "hey bob!"}
	time.Sleep(50 * time.Millisecond)

	// A conversation should be created
	assert.Equal(t, 1, store.ConversationCount())

	// Get the conversation
	convos := store.GetConversationsCopy()
	var convoID int
	var convoUserOne, convoUserTwo string
	for _, c := range convos {
		convoID = c.ID
		convoUserOne = c.UserOne
		convoUserTwo = c.UserTwo
	}
	require.NotZero(t, convoID)
	// canonical: alice < bob
	assert.Equal(t, "alice", convoUserOne)
	assert.Equal(t, "bob", convoUserTwo)

	// DM should be saved
	dms := store.GetDirectMsgsCopy()
	assert.Len(t, dms, 1)
	assert.Equal(t, "hey bob!", dms[0].Content)
	assert.Equal(t, convoID, dms[0].ConversationID)

	// Both alice and bob should receive the DM with conversation ID
	aliceMsgs := drainMessages(alice)
	bobMsgs := drainMessages(bob)

	aliceDM := findMsg(aliceMsgs, MsgDirectMessage, "hey bob!")
	bobDM := findMsg(bobMsgs, MsgDirectMessage, "hey bob!")

	require.NotNil(t, aliceDM, "alice should receive the DM echo")
	require.NotNil(t, bobDM, "bob should receive the DM")
	assert.Equal(t, convoID, aliceDM.ConversationID)
	assert.Equal(t, convoID, bobDM.ConversationID)
}

func TestDirectMessage_ReusesConversation(t *testing.T) {
	hub, store := setupHub(t)

	alice := registerClient(t, hub, "alice")
	bob := registerClient(t, hub, "bob")
	drainMessages(alice)
	drainMessages(bob)

	// First DM
	hub.SendMessage <- Message{Type: MsgDirectMessage, User: "alice", Receiver: "bob", Content: "msg 1"}
	time.Sleep(50 * time.Millisecond)

	// Second DM in reverse direction
	hub.SendMessage <- Message{Type: MsgDirectMessage, User: "bob", Receiver: "alice", Content: "msg 2"}
	time.Sleep(50 * time.Millisecond)

	// Still only one conversation
	assert.Equal(t, 1, store.ConversationCount())

	// Both DMs should share the same conversation ID
	dms := store.GetDirectMsgsCopy()
	assert.Equal(t, dms[0].ConversationID, dms[1].ConversationID)
}

func TestDirectMessage_ToOfflineUser(t *testing.T) {
	hub, store := setupHub(t)

	// Create both users so alice is "known"
	alice := registerClient(t, hub, "alice")
	drainMessages(alice)

	// alice sends to "charlie" who exists but is offline
	store.CreateUserIfNotExists(context.Background(), "charlie")

	hub.SendMessage <- Message{Type: MsgDirectMessage, User: "alice", Receiver: "charlie", Content: "hello charlie!"}
	time.Sleep(50 * time.Millisecond)

	// DM should be saved even if charlie is offline
	assert.Equal(t, 1, store.DirectMsgCount())

	// alice should receive the echo
	aliceMsgs := drainMessages(alice)
	dm := findMsg(aliceMsgs, MsgDirectMessage, "hello charlie!")
	require.NotNil(t, dm)
}

func TestGetConversations(t *testing.T) {
	hub, _ := setupHub(t)

	alice := registerClient(t, hub, "alice")
	bob := registerClient(t, hub, "bob")
	drainMessages(alice)
	drainMessages(bob)

	// Create a conversation via DM
	hub.SendMessage <- Message{Type: MsgDirectMessage, User: "alice", Receiver: "bob", Content: "hi"}
	time.Sleep(50 * time.Millisecond)
	drainMessages(alice)
	drainMessages(bob)

	// Request conversations
	hub.SendMessage <- Message{Type: MsgGetConversations, User: "alice"}
	time.Sleep(50 * time.Millisecond)

	msgs := drainMessages(alice)
	found := false
	for _, m := range msgs {
		if m.Type == MsgConversationsList {
			found = true
			assert.Contains(t, m.Content, "alice")
			assert.Contains(t, m.Content, "bob")
		}
	}
	assert.True(t, found, "alice should receive conversations list")
}

func TestNextDirectMessages_Pagination(t *testing.T) {
	hub, store := setupHub(t)

	alice := registerClient(t, hub, "alice")
	bob := registerClient(t, hub, "bob")
	drainMessages(alice)
	drainMessages(bob)

	// Send 3 DMs to create a conversation
	for i := 0; i < 3; i++ {
		hub.SendMessage <- Message{Type: MsgDirectMessage, User: "alice", Receiver: "bob", Content: "msg"}
		time.Sleep(30 * time.Millisecond)
	}
	drainMessages(alice)
	drainMessages(bob)

	// Get the conversation ID
	convos := store.GetConversationsCopy()
	var convoID int
	for id := range convos {
		convoID = id
	}

	// Request page with lastid = 3 (should return messages with id < 3)
	hub.SendMessage <- Message{Type: MsgNextDirectMessages, User: "alice", Content: "3", ConversationID: convoID}
	time.Sleep(50 * time.Millisecond)

	msgs := drainMessages(alice)
	var dmMsgs []Message
	for _, m := range msgs {
		if m.Type == MsgDirectMessage {
			dmMsgs = append(dmMsgs, m)
		}
	}
	// Should get messages with ID 1 and 2
	assert.Equal(t, 2, len(dmMsgs))
}

func TestBroadcastMessage(t *testing.T) {
	hub, store := setupHub(t)

	alice := registerClient(t, hub, "alice")
	bob := registerClient(t, hub, "bob")
	drainMessages(alice)
	drainMessages(bob)

	hub.SendMessage <- Message{Type: MsgBroadcast, User: "alice", Content: "hello everyone!"}
	time.Sleep(50 * time.Millisecond)

	// Both should receive it
	bobMsgs := drainMessages(bob)
	assert.True(t, hasMsg(bobMsgs, MsgBroadcast, "hello everyone!"), "bob should receive broadcast")

	// Message saved in store
	assert.Equal(t, 1, store.RoomMessageCount())
}

func TestMultipleRegistrations(t *testing.T) {
	hub, _ := setupHub(t)

	c1 := registerClient(t, hub, "alice")
	drainMessages(c1)

	// Register another client with same name (simulates reconnect)
	c2 := registerClient(t, hub, "alice")
	drainMessages(c2)

	// Verify the new client works by sending a message
	// c2 should be the active client now
	hub.SendMessage <- Message{Type: MsgGetConversations, User: "alice"}
	time.Sleep(50 * time.Millisecond)

	msgs := drainMessages(c2)
	assert.True(t, hasMsgType(msgs, MsgConversationsList), "c2 should receive response")
}

// ===================== Helpers =====================

func findMsg(msgs []Message, typ MessageType, content string) *Message {
	for _, m := range msgs {
		if m.Type == typ && m.Content == content {
			return &m
		}
	}
	return nil
}
