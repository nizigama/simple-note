package models

import (
	boltDB "github.com/nizigama/simple-note/services/database"
)

type Note struct {
	ID    int
	Title string
	Body  string
}

const (
	NoteTableName string = "Notes"
)

// Save persists the user in the struct in the database
func (n Note) Save() (uint64, error) {

	simpleNote := map[string]interface{}{
		"title": n.Title,
		"note":  n.Body,
	}

	return boltDB.Store(simpleNote, NoteTableName)
}

func ReadNote(noteID uint64) (Note, error) {

	simpleNote, err := boltDB.Show(noteID, NoteTableName)

	if err != nil {
		return Note{}, err
	}

	return Note{
		Title: simpleNote["title"].(string),
		Body:  simpleNote["note"].(string),
	}, nil
}

func ReadAllNotes() ([]Note, error) {

	var notes []Note
	res, err := boltDB.All(NoteTableName)

	if err != nil {
		return nil, err
	}

	for _, simpleNote := range res {
		notes = append(notes, Note{
			ID:    int(simpleNote["itemID"].(uint64)),
			Title: simpleNote["title"].(string),
			Body:  simpleNote["note"].(string),
		})
	}

	return notes, nil
}

func UpdateNote(n Note, itemID int) error {
	simpleNote := map[string]interface{}{
		"title": n.Title,
		"note":  n.Body,
	}

	return boltDB.Update(simpleNote, NoteTableName, uint64(itemID))
}
