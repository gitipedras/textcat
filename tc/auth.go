package tc

type Auth struct {
	AuthManager *AuthManager
	//SessionManager *SessionManager
}

type AuthManager struct {
	
}

type (am *AuthManager) AddUser(username string, password string) error {
	
}