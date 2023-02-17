package database

import (
	"context"
	"github.com/go-faster/errors"
	"time"
)

func (db Database) GetNotes(ctx context.Context, userID uint64) ([]Note, error) {
	// Create the query
	rows, err := db.db.Query(ctx, "SELECT id, title, note_text, created_at, deadline FROM normal_notes WHERE for_user=$1 ORDER BY note_order", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Iterate over rows
	notes := make([]Note, 0)
	for rows.Next() {
		var note Note
		err = rows.Scan(&note.ID, &note.Title, &note.Text, &note.CreatedAt, &note.Deadline)
		if err != nil {
			return nil, errors.Wrap(err, "cannot scan row")
		}
		// Get images
		note.Images, err = db.getNoteImages(ctx, note.ID)
		if err != nil {
			return nil, errors.Wrap(err, "cannot get note images")
		}
		// Add to list
		notes = append(notes, note)
	}
	return notes, nil
}

func (db Database) GetNote(ctx context.Context, userID uint64, noteID int32) (Note, error) {
	var note Note
	err := db.db.QueryRow(ctx, "SELECT id, title, note_text, created_at, deadline FROM normal_notes WHERE id=$1 AND for_user=$2", noteID, userID).
		Scan(&note.ID, &note.Title, &note.Text, &note.CreatedAt, &note.Deadline)
	if err != nil {
		return Note{}, err
	}
	note.Images, err = db.getNoteImages(ctx, note.ID)
	if err != nil {
		return Note{}, errors.Wrap(err, "cannot get note images")
	}
	return note, nil
}

func (db Database) getNoteImages(ctx context.Context, noteID int32) ([]SavedImage, error) {
	rows, err := db.db.Query(ctx, "SELECT id, filename FROM normal_note_photos WHERE for_note=$1", noteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Get rows
	images := make([]SavedImage, 0)
	for rows.Next() {
		var savedImage SavedImage
		err = rows.Scan(&savedImage.ID, &savedImage.Filename)
		if err != nil {
			return nil, errors.Wrap(err, "cannot scan row")
		}
	}
	return images, nil
}

func (db Database) InsertNote(ctx context.Context, forUser uint64, note Note) (Note, error) {
	// Create a transaction to make everything atomic
	tx, err := db.db.Begin(ctx)
	if err != nil {
		return Note{}, err
	}
	defer tx.Rollback(ctx)
	// At first insert the note itself
	err = tx.QueryRow(ctx, "INSERT INTO normal_notes (for_user, title, note_text, background_color, created_at, deadline, note_order) VALUES ($1, $2, $3, NULL, NOW(), $4, (SELECT MAX(note_order) FROM normal_notes WHERE for_user=$1)) RETURNING id",
		forUser, note.Title, note.Text, note.Deadline).Scan(&note.ID)
	if err != nil {
		return Note{}, errors.Wrap(err, "cannot insert note itself")
	}
	// Now insert the pics
	for i := range note.Images {
		err = tx.QueryRow(ctx, "INSERT INTO normal_note_photos (id, filename, for_note) VALUES (gen_random_uuid(), $1, $2)",
			note.Images[i].Filename, note.ID).Scan(&note.Images[i].ID)
		if err != nil {
			return Note{}, errors.Wrap(err, "cannot insert image of note")
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		return Note{}, errors.Wrap(err, "cannot commit")
	}
	return note, nil
}

func (db Database) DeleteNote(ctx context.Context, forUser uint64, noteID int32) error {
	// No transaction just fuck it
	rows, err := db.db.Exec(ctx, "DELETE FROM normal_notes WHERE id=$1 AND for_user=$2", noteID, forUser)
	if err != nil {
		return errors.Wrap(err, "cannot delete note itself")
	}
	if rows.RowsAffected() == 0 { // user trying to delete another note
		return nil // not errors. Make frontend happy
	}
	_, _ = db.db.Exec(ctx, "DELETE FROM normal_note_photos WHERE for_note=$1", noteID)
	// we don't care about photo errors
	return nil
}

func (db Database) ReorderNote(ctx context.Context, forUser uint64, note1, note2 int32) error {
	// Get old orders
	var note1Order, note2Order int32
	err := db.db.QueryRow(ctx, "SELECT note_order FROM normal_notes WHERE id=$1 AND for_user=$2",
		note1, forUser).Scan(&note1Order)
	if err != nil {
		return err
	}
	err = db.db.QueryRow(ctx, "SELECT note_order FROM normal_notes WHERE id=$1 AND for_user=$2",
		note2, forUser).Scan(&note2Order)
	if err != nil {
		return err
	}
	// Update them in database. Also, we don't care about errors
	_, _ = db.db.Exec(ctx, "UPDATE normal_notes SET note_order=$1 WHERE id=$2", note1Order, note2)
	_, _ = db.db.Exec(ctx, "UPDATE normal_notes SET note_order=$1 WHERE id=$2", note2Order, note1)
	return nil
}

func (db Database) NoteDeadline(ctx context.Context, forUser uint64, noteID int32, deadline *time.Time) error {
	_, err := db.db.Exec(ctx, "UPDATE normal_notes SET deadline=$1 WHERE id=$2 AND for_user=$3", deadline, noteID, forUser)
	return err
}
