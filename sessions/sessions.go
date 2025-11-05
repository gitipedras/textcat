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

    /* internal */
    "textcat/models"
    "log/slog"
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

// Returns true if the session exists and matches, otherwise false.
func (sm *SessionManager) CheckByUsername(username string) bool {
    // Acquire a Read Lock to safely iterate over the map
    sm.Mu.RLock()
    defer sm.Mu.RUnlock()

    // The range loop iterates over the keys (token) and values (session) of the map
    for _, session := range sm.Sessions {
        // You should also check if the session pointer is nil, just in case
        if session != nil && session.Username == username {
            // Found a matching session for the given username
            return true
        }
    }

    // If the loop completes without finding a match
    return false
}

func (sm *SessionManager) Exists(token string) bool {
    sm.Mu.RLock()
    defer sm.Mu.RUnlock()

    session, ok := sm.Sessions[token]
    if !ok || session == nil || session.Conn == nil {
        return false
    }

    return true
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
    models.App.Log.Info("sending message to client", slog.String("token", token))

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
