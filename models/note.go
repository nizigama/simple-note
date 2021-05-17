package models

import (
	"fmt"

	mysqlDB "github.com/nizigama/simple-note/services/database"
)

type Note struct {
	ID      int
	Title   string
	Body    string
	OwnerID int
}

var (
	NoteMigration mysqlDB.Migration = mysqlDB.Migration{
		TableName: "Notes",
		Definition: []mysqlDB.Column{
			{
				Name:  "id",
				Type:  "int",
				Extra: "auto_increment",
				Key:   "primary key",
			},
			{
				Name: "title",
				Type: "varchar(255)",
			},
			{
				Name: "note",
				Type: "longtext",
			},
			{
				Name: "ownerID",
				Type: "int",
			},
		},
	}
)

// Save persists the user in the struct in the database
func (n Note) Save() (uint64, error) {

	query := fmt.Sprintf("INSERT INTO %s(title, note, ownerID) VALUES(?,?,?)", NoteMigration.TableName)

	r, err := mysqlDB.MysqlDB.Exec(query, n.Title, n.Body, n.OwnerID)

	if err != nil {
		return 0, err
	}

	var noteID int64

	if noteID, err = r.LastInsertId(); err != nil {
		return 0, err
	}

	return uint64(noteID), nil
}

func ReadNote(noteID uint64) (Note, error) {

	query := fmt.Sprintf("SELECT id, title, note, ownerID FROM %s WHERE id = ?", NoteMigration.TableName)
	row := mysqlDB.MysqlDB.QueryRow(query, int(noteID))

	var note Note

	err := row.Scan(&note.ID, &note.Title, &note.Body, &note.OwnerID)

	if err != nil {
		return Note{}, err
	}

	return note, nil
}

func ReadAllUserNotes(userID int) ([]Note, error) {

	var notes []Note

	query := fmt.Sprintf("SELECT id, title, note, ownerID FROM %s WHERE ownerID = ?", NoteMigration.TableName)
	rows, err := mysqlDB.MysqlDB.Query(query, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		note := Note{}

		err := rows.Scan(&note.ID, &note.Title, &note.Body, &note.OwnerID)

		if err != nil {
			return nil, err
		}

		notes = append(notes, note)
	}

	return notes, nil
}

func UpdateNote(n Note, itemID int) error {

	query := fmt.Sprintf("UPDATE %s SET title = ?, note = ? WHERE id = ? AND ownerID = ?", NoteMigration.TableName)

	_, err := mysqlDB.MysqlDB.Exec(query, n.Title, n.Body, itemID, n.OwnerID)

	if err != nil {
		return err
	}

	return nil
}

func DeleteNote(itemID int) error {

	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", NoteMigration.TableName)

	_, err := mysqlDB.MysqlDB.Exec(query, itemID)

	if err != nil {
		return err
	}

	return nil
}
