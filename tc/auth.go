package tc

import (
	"fmt"
	"time"
	"errors"
)

type Textcat struct {
	Function LogicBridge
	//SessionManager *SessionManager
}

var ErrNotFound = errors.New("data does not found")

func (tc *Textcat) CreateUser(username string, password string) error {
	fmt.Printf("user register: %s: %s", username, password)

	tblErr := tc.Function.CreateTable("users")
	if tblErr != nil {
		return MakeError("server_error:", tblErr)
	}

	// Validate username
	if !IsValidUsername(username) {
		return errors.New("error: invalid username: must be alphanumeric, can contain '-' and '_'")
	}

	// Build user object
	user := User{
		Username:  username,
		Password:  password,
		Created:   time.Now(),
		LastLogin: time.Now(),
	}

	// Delegate existence check to the Handler
	var existing User
	err := tc.Function.GetData("users", func(v any) bool {
		u := v.(*User)
		return u.Username == username // return the sub-func() if user exists
	}, &existing)
	

	if err == nil {
		return errors.New("error: username is taken")
	}

	// Delegate storage to the Handler
	err = tc.Function.StoreData("users", user)
	if err != nil {
		return MakeError("server_error", err)
	}
	return errors.New("ok")
}

func (tc *Textcat) LoginUser(username string, password string) error {
	fmt.Printf("user register: %s: %s", username, password)

	// Validate username
	if !IsValidUsername(username) {
		return errors.New("invalid username: must be alphanumeric, can contain '-' and '_'")
	}

	// Build user object
	/**user := User{
		Username:  username,
		Password:  password,
		Created:   time.Now(),
		LastLogin: time.Now(),
	}**/

	// Delegate existence check to the Handler
	var existing User
	err := tc.Function.GetData("users", func(v any) bool {
		u := v.(*User)
		return u.Username == username
	}, &existing)
	if err == nil {
		return errors.New("username is taken")
	} else if err != nil && !errors.Is(err, ErrNotFound) {
		return errors.New("user does not exist")
	}

	// Delegate storage to the Handler
	return nil
}