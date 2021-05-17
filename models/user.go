package models

import (
	"fmt"

	mysqlDB "github.com/nizigama/simple-note/services/database"
)

type User struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
	Picture   string
}

var (
	UserMigration mysqlDB.Migration = mysqlDB.Migration{
		TableName: "Users",
		Definition: []mysqlDB.Column{
			{
				Name:  "id",
				Type:  "int",
				Extra: "auto_increment",
				Key:   "primary key",
			},
			{
				Name: "firstName",
				Type: "varchar(255)",
			},
			{
				Name: "lastName",
				Type: "varchar(255)",
			},
			{
				Name: "email",
				Type: "varchar(255)",
			},
			{
				Name: "password",
				Type: "varchar(255)",
			},
			{
				Name: "picture",
				Type: "varchar(255)",
			},
		},
	}

	defaultPicture string = "avatar.png"
)

// Save persists the user in the struct in the database
func (u User) Save() (uint64, error) {

	query := fmt.Sprintf("INSERT INTO %s(firstName, lastName, email, password, picture) VALUES(?,?,?,?,\"%s\")", UserMigration.TableName, defaultPicture)

	r, err := mysqlDB.MysqlDB.Exec(query, u.FirstName, u.LastName, u.Email, u.Password)

	if err != nil {
		return 0, err
	}

	var userID int64

	if userID, err = r.LastInsertId(); err != nil {
		return 0, err
	}

	return uint64(userID), nil
}

func ReadUser(userID uint64) (User, error) {

	query := fmt.Sprintf("SELECT firstName, lastName, email, password, picture FROM %s WHERE id = ?", UserMigration.TableName)
	row := mysqlDB.MysqlDB.QueryRow(query, int(userID))

	var user User

	err := row.Scan(&user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Picture)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func ReadAllUsers() ([]User, error) {

	var users []User
	query := fmt.Sprintf("SELECT firstName, lastName, email, password, picture FROM %s", UserMigration.TableName)
	rows, err := mysqlDB.MysqlDB.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		user := User{}

		err := rows.Scan(&user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Picture)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func ReadSingleUserByEmail(userEmail string) (User, uint64, error) {

	query := fmt.Sprintf("SELECT id, firstName, lastName, email, password, picture FROM %s WHERE email = ?", UserMigration.TableName)
	row := mysqlDB.MysqlDB.QueryRow(query, userEmail)

	var user User
	var userID int

	err := row.Scan(&userID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Picture)

	if err != nil {
		return User{}, 0, err
	}

	return user, uint64(userID), nil
}

func UpdateUser(u User, itemID int) error {

	query := fmt.Sprintf("UPDATE %s SET firstName = ?, lastName = ?, email = ?, password = ?, picture = ? WHERE id = ?", UserMigration.TableName)

	_, err := mysqlDB.MysqlDB.Exec(query, u.FirstName, u.LastName, u.Email, u.Password, u.Picture, itemID)

	if err != nil {
		return err
	}

	return nil
}

func DeleteUser(itemID int) error {

	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", UserMigration.TableName)

	_, err := mysqlDB.MysqlDB.Exec(query, itemID)

	if err != nil {
		return err
	}

	return nil
}
