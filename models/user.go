package users

import (
	boltDB "github.com/nizigama/simple-note/services/database"
)

type User struct {
	FirstName string
	LastName  string
	Email     string
	Password  []byte
}

const (
	TableName string = "Users"
)

// Save persists the user in the struct in the database
func (u User) Save() error {

	userMap := map[string]interface{}{
		"firstName": u.FirstName,
		"lastName":  u.LastName,
		"email":     u.LastName,
		"password":  u.Password,
	}

	return boltDB.Store(userMap, TableName)
}
