package models

import (
	boltDB "github.com/nizigama/simple-note/services/database"
)

type User struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
	Picture   string
}

const (
	UserTableName string = "Users"
)

// Save persists the user in the struct in the database
func (u User) Save() (uint64, error) {

	userMap := map[string]interface{}{
		"firstName": u.FirstName,
		"lastName":  u.LastName,
		"email":     u.Email,
		"password":  u.Password,
		"picture":   "avatar.png",
	}

	return boltDB.Store(userMap, UserTableName)
}

func ReadUser(userID uint64) (User, error) {

	user, err := boltDB.Show(userID, UserTableName)

	if err != nil {
		return User{}, err
	}

	return User{
		FirstName: user["firstName"].(string),
		LastName:  user["lastName"].(string),
		Email:     user["email"].(string),
		Password:  user["password"].(string),
		Picture:   user["picture"].(string),
	}, nil
}

func ReadAllUsers() ([]User, error) {

	var users []User
	res, err := boltDB.All(UserTableName)

	if err != nil {
		return nil, err
	}

	for _, user := range res {
		users = append(users, User{
			FirstName: user["firstName"].(string),
			LastName:  user["lastName"].(string),
			Email:     user["email"].(string),
			Password:  user["password"].(string),
			Picture:   user["picture"].(string),
		})
	}

	return users, nil
}

func ReadSingleUserByEmail(userEmail string) (User, uint64, error) {

	user, index, err := boltDB.SingleByStringField(UserTableName, "email", userEmail)

	if err != nil {
		return User{}, 0, err
	}

	return User{
		FirstName: user["firstName"].(string),
		LastName:  user["lastName"].(string),
		Email:     user["email"].(string),
		Password:  user["password"].(string),
		Picture:   user["picture"].(string),
	}, index, nil
}

func UpdateUser(u User, itemID int) error {
	userMap := map[string]interface{}{
		"firstName": u.FirstName,
		"lastName":  u.LastName,
		"email":     u.Email,
		"password":  u.Password,
		"picture":   u.Picture,
	}

	return boltDB.Update(userMap, UserTableName, uint64(itemID))
}
