package tc

import(
	"sync"
	"crypto/rand"
	"encoding/hex"
	"time"
)

type Session struct {
	Mutex    sync.Mutex
	Username string
	LastSeen time.Time
	Data     map[string]any
}

type SessionManager struct {
	Mutex    sync.RWMutex
	Sessions map[string]*Session
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		Sessions: make(map[string]*Session),
	}
}

// Check if session exists
func (sm *SessionManager) CheckSessionUsername(username string) bool {
	sm.Mutex.RLock()
	defer sm.Mutex.RUnlock()

	for _, session := range sm.Sessions {
		if session.Username == username {
			return true
		}
	}
	return false
}


// Add a session
func (sm *SessionManager) AddSession(id string, username string) {
	sm.Mutex.Lock()
	defer sm.Mutex.Unlock()

	sm.Sessions[id] = &Session{
		Username: username,
		LastSeen: time.Now(),
		Data:     make(map[string]any),
	}
}

// get a session
func (sm *SessionManager) GetSession(id string) (*Session, bool) {
	sm.Mutex.RLock()
	defer sm.Mutex.RUnlock()

	s, ok := sm.Sessions[id]
	return s, ok
}

// loop and check for an empty session
func (sm *SessionManager) GetUniqueID() string {
	for {
		id := make([]byte, 16) // 128-bit ID
		_, _ = rand.Read(id)
		uid := hex.EncodeToString(id)

		sm.Mutex.RLock()
		_, exists := sm.Sessions[uid]
		sm.Mutex.RUnlock()

		if !exists {
			return uid
		}
	}
}
