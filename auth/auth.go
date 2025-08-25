package auth

import (
	/* websocket */
	"github.com/gorilla/websocket"

	/* other stuff */
	"time"
	"sync"

	/* internal  */
	"textcat/database"
	"textcat/models"
)

func UserLogin(msg models.WsIncome) {
	ok := database.CheckUser(msg.Username)
	if ok {
		good := database.CheckPass(msg.Username, msg.SessionToken)
		if good {
			// LOGIN THE USER
		}

	} else {
		// user not found
	}
}

func UserRegister() {
	
}


// ===================================== //
//            SESSION MANAGER            //
// ------------------------------------- //

// Session represents a single connected user
type Session struct {
    Username     string
    SessionToken string
    Conn         *websocket.Conn
    ConnectedAt  time.Time
}

// SessionManager manages all active sessions
type SessionManager struct {
    sessions map[string]*Session // key: session token
    mu       sync.RWMutex
}

// NewSessionManager creates a new session manager
func NewSessionManager() *SessionManager {
    return &SessionManager{
        sessions: make(map[string]*Session),
    }
}

// Add adds a new session
func (sm *SessionManager) Add(session *Session) {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    sm.sessions[session.SessionToken] = session
}

// Get retrieves a session by token
func (sm *SessionManager) Get(token string) (*Session, bool) {
    sm.mu.RLock()
    defer sm.mu.RUnlock()
    s, ok := sm.sessions[token]
    return s, ok
}

// Remove deletes a session
func (sm *SessionManager) Remove(token string) {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    delete(sm.sessions, token)
}

// SendToClient sends a message to a single session
func (sm *SessionManager) SendToClient(token string, message []byte) error {
    sm.mu.RLock()
    session, ok := sm.sessions[token]
    sm.mu.RUnlock()

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