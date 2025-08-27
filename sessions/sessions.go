package sessions

// ===================================== //
//            SESSION MANAGER            //
// ------------------------------------- //

import (
    /* session manager */
    "crypto/rand"
    "encoding/hex"
    "sync"
    "time"

    /* websockets */
    "github.com/gorilla/websocket"
)

// Session represents a single connected user
type Session struct {
    Username     string
    SessionToken string
    Conn         *websocket.Conn
    ConnectedAt  time.Time
}

// SessionManager manages all active sessions
type SessionManager struct {
    Sessions map[string]*Session // key: session token
    Mu       sync.RWMutex
}

// NewSessionManager creates a new session manager
func NewSessionManager() *SessionManager {
    return &SessionManager{
        Sessions: make(map[string]*Session),
    }
}

// Add adds a new session
func (sm *SessionManager) Add(session *Session) {
    sm.Mu.Lock()
    defer sm.Mu.Unlock()
    sm.Sessions[session.SessionToken] = session
}

// Get retrieves a session by token
func (sm *SessionManager) Get(token string) (*Session, bool) {
    sm.Mu.RLock()
    defer sm.Mu.RUnlock()
    s, ok := sm.Sessions[token]
    return s, ok
}

// Remove deletes a session
func (sm *SessionManager) Remove(token string) {
    sm.Mu.Lock()
    defer sm.Mu.Unlock()
    delete(sm.Sessions, token)
}

// SendToClient sends a message to a single session
func (sm *SessionManager) SendToClient(token string, message []byte) error {
    sm.Mu.RLock()
    session, ok := sm.Sessions[token]
    sm.Mu.RUnlock()

    if !ok || session.Conn == nil {
        return nil // session not found or disconnected
    }

    err := session.Conn.WriteMessage(websocket.TextMessage, message)
    if err != nil {
        // clean up disconnected client
        sm.Remove(token)
        session.Conn.Close()
    }
    return err
}

// GenerateToken creates a new random session token (hex string)
func (sm *SessionManager) GenerateToken(length int) (string, error) {
    b := make([]byte, length)
    _, err := rand.Read(b)
    if err != nil {
        return "", err
    }
    return hex.EncodeToString(b), nil
}
