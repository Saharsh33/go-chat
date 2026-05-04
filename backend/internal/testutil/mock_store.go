package testutil

import (
	"chat-server/internal/models"
	"context"
	"sync"
	"time"
)

// MockStore is an in-memory implementation of storage.StorageInterface for testing.
type MockStore struct {
	mu sync.Mutex

	Users         map[string]int // username -> id
	Rooms         map[int]*models.StoredRoom
	RoomMembers   map[int]map[string]time.Time // roomID -> username -> joinedAt
	RoomMessages  []models.Message
	Conversations map[int]*models.Conversation
	DirectMsgs    []models.Message

	nextUserID   int
	nextRoomID   int
	nextMsgID    int
	nextConvoID  int
	nextDmID     int
}

func NewMockStore() *MockStore {
	return &MockStore{
		Users:         make(map[string]int),
		Rooms:         make(map[int]*models.StoredRoom),
		RoomMembers:   make(map[int]map[string]time.Time),
		RoomMessages:  []models.Message{},
		Conversations: make(map[int]*models.Conversation),
		DirectMsgs:    []models.Message{},
		nextUserID:    1,
		nextRoomID:    1,
		nextMsgID:     1,
		nextConvoID:   1,
		nextDmID:      1,
	}
}

// ----- User operations -----

func (m *MockStore) CreateUserIfNotExists(_ context.Context, user string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.Users[user]; !ok {
		m.Users[user] = m.nextUserID
		m.nextUserID++
	}
}

func (m *MockStore) GetUserByName(_ context.Context, user string) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	id, ok := m.Users[user]
	if !ok {
		return 0, context.DeadlineExceeded // simulate "not found"
	}
	return id, nil
}

// ----- Room operations -----

func (m *MockStore) CreateRoom(_ context.Context, room string, creator string) (*models.StoredRoom, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	r := &models.StoredRoom{ID: m.nextRoomID, Name: room}
	m.Rooms[m.nextRoomID] = r
	// also add user to room
	m.RoomMembers[m.nextRoomID] = map[string]time.Time{creator: time.Now()}
	m.nextRoomID++
	return r, nil
}

func (m *MockStore) GetRoomByName(_ context.Context, name string) (*models.StoredRoom, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, r := range m.Rooms {
		if r.Name == name {
			return r, nil
		}
	}
	// Match postgres behavior: when no rows, QueryContext returns no error,
	// the loop doesn't execute, and a zero-value StoredRoom is returned.
	return &models.StoredRoom{}, nil
}

func (m *MockStore) GetRoomById(_ context.Context, id int) (*models.StoredRoom, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	r, ok := m.Rooms[id]
	if !ok {
		// Match postgres behavior: when no rows, returns zero-value + nil
		return &models.StoredRoom{}, nil
	}
	return r, nil
}

func (m *MockStore) GetAllRooms(_ context.Context) ([]*models.StoredRoom, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var rooms []*models.StoredRoom
	for _, r := range m.Rooms {
		rooms = append(rooms, r)
	}
	return rooms, nil
}

// ----- Room member operations -----

func (m *MockStore) AddUserToRoom(_ context.Context, roomId int, username string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.RoomMembers[roomId] == nil {
		m.RoomMembers[roomId] = make(map[string]time.Time)
	}
	m.RoomMembers[roomId][username] = time.Now()
	return nil
}

func (m *MockStore) RemoveUserFromRoom(_ context.Context, roomId int, username string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.RoomMembers[roomId], username)
	return nil
}

func (m *MockStore) GetUsersInRoom(_ context.Context, roomId int) ([]*models.RoomMember, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var members []*models.RoomMember
	for uname, jat := range m.RoomMembers[roomId] {
		members = append(members, &models.RoomMember{RoomID: roomId, Username: uname, JoinedAt: jat})
	}
	return members, nil
}

func (m *MockStore) GetRoomsOfUser(_ context.Context, username string) ([]*models.StoredRoom, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var rooms []*models.StoredRoom
	for roomID, members := range m.RoomMembers {
		if _, ok := members[username]; ok {
			if r, exists := m.Rooms[roomID]; exists {
				rooms = append(rooms, r)
			}
		}
	}
	return rooms, nil
}

// ----- Message operations -----

func (m *MockStore) SaveMessage(_ context.Context, msg string, roomId int, userName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.RoomMessages = append(m.RoomMessages, models.Message{
		ID:        m.nextMsgID,
		Room:      roomId,
		User:      userName,
		Content:   msg,
		CreatedAt: time.Now(),
	})
	m.nextMsgID++
	return nil
}

func (m *MockStore) GetRecentMessages(_ context.Context, roomId int, limit int, lastid int) ([]models.Message, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var result []models.Message
	// Iterate in reverse (newest first)
	for i := len(m.RoomMessages) - 1; i >= 0; i-- {
		msg := m.RoomMessages[i]
		if msg.Room != roomId {
			continue
		}
		if lastid > 0 && msg.ID >= lastid {
			continue
		}
		result = append(result, msg)
		if len(result) >= limit {
			break
		}
	}
	return result, nil
}

// ----- Conversation operations -----

func (m *MockStore) GetOrCreateConversation(_ context.Context, userA string, userB string) (*models.Conversation, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// canonical ordering
	one, two := userA, userB
	if one > two {
		one, two = two, one
	}

	for _, c := range m.Conversations {
		if c.UserOne == one && c.UserTwo == two {
			c.UpdatedAt = time.Now()
			return c, nil
		}
	}
	c := &models.Conversation{
		ID:        m.nextConvoID,
		UserOne:   one,
		UserTwo:   two,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	m.Conversations[m.nextConvoID] = c
	m.nextConvoID++
	return c, nil
}

func (m *MockStore) SendDirectMessage(_ context.Context, msg string, conversationID int, sender string, receiver string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.DirectMsgs = append(m.DirectMsgs, models.Message{
		ID:             m.nextDmID,
		User:           sender,
		Receiver:       receiver,
		Content:        msg,
		ConversationID: conversationID,
		CreatedAt:      time.Now(),
	})
	m.nextDmID++
	return nil
}

func (m *MockStore) GetRecentDirectMessages(_ context.Context, conversationID int, limit int, lastid int) ([]models.Message, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var result []models.Message
	for i := len(m.DirectMsgs) - 1; i >= 0; i-- {
		dm := m.DirectMsgs[i]
		if dm.ConversationID != conversationID {
			continue
		}
		if lastid > 0 && dm.ID >= lastid {
			continue
		}
		result = append(result, dm)
		if len(result) >= limit {
			break
		}
	}
	return result, nil
}

func (m *MockStore) GetConversationsOfUser(_ context.Context, username string) ([]models.Conversation, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var result []models.Conversation
	for _, c := range m.Conversations {
		if c.UserOne == username || c.UserTwo == username {
			result = append(result, *c)
		}
	}
	return result, nil
}

// ----- Thread-safe accessors for test assertions -----

func (m *MockStore) GetDirectMsgsCopy() []models.Message {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := make([]models.Message, len(m.DirectMsgs))
	copy(cp, m.DirectMsgs)
	return cp
}

func (m *MockStore) GetConversationsCopy() map[int]*models.Conversation {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := make(map[int]*models.Conversation, len(m.Conversations))
	for k, v := range m.Conversations {
		conv := *v
		cp[k] = &conv
	}
	return cp
}

func (m *MockStore) GetRoomMessagesCopy() []models.Message {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := make([]models.Message, len(m.RoomMessages))
	copy(cp, m.RoomMessages)
	return cp
}

func (m *MockStore) ConversationCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.Conversations)
}

func (m *MockStore) DirectMsgCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.DirectMsgs)
}

func (m *MockStore) RoomMessageCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.RoomMessages)
}
