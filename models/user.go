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
func (u User) Save() (uint64, error) {

	userMap := map[string]interface{}{
		"firstName": u.FirstName,
		"lastName":  u.LastName,
		"email":     u.Email,
		"password":  u.Password,
	}

	return boltDB.Store(userMap, TableName)
}

func Read(userID uint64) (User, error) {

	user, err := boltDB.Show(userID, TableName)

	if err != nil {
		return User{}, err
	}

	return User{
		FirstName: user["firstName"].(string),
		LastName:  user["lastName"].(string),
		Email:     user["email"].(string),
	}, nil
}
