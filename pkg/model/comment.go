package model

import (
	"database/sql"
	"errors"
)

type Comment struct {
	Id       int    `json:"id"`
	Comment  string `json:"comment"`
	Id_User  int    `json:"id_user"`
	Id_Movie int    `json:"id_movie"`
}

func GetComments(db *sql.DB) ([]Comment, error) {
	rows, err := db.Query(`SELECT id, comment, id_user, id_movie FROM comments`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.Id, &comment.Comment, &comment.Id_User, &comment.Id_Movie); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func CreateComment(db *sql.DB, comment *Comment) error {
	query := "INSERT INTO comments (comment, id_user, id_movie) VALUES (?, ?, ?)"

	if comment.Comment == "" {
		err := errors.New("Comment is required")
		return err
	}
	result, err := db.Exec(query, comment.Comment, comment.Id_User, comment.Id_Movie)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	comment.Id = int(id)

	return nil
}

func GetCommentByID(db *sql.DB, id int) (*Comment, error) {
	var comment Comment
	query := "SELECT id, comment, id_user, id_movie FROM comments WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&comment.Id, &comment.Comment, &comment.Id_User, &comment.Id_Movie)
	if err != nil {
		return nil, err
	}

	return &comment, nil

}

func GetCommentByMovieUser(db *sql.DB, id_movie int, id_user int) (*Comment, error) {
	var comment Comment
	query := "SELECT id, comment, id_user, id_movie FROM comments WHERE id_movie = ? AND id_user = ?"
	err := db.QueryRow(query, id_movie, id_user).Scan(&comment.Id, &comment.Comment, &comment.Id_User, &comment.Id_Movie)
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func UpdateComment(db *sql.DB, comment *Comment) error {
	query := "UPDATE comments SET comment = ? WHERE id = ?"
	_, err := db.Exec(query, comment.Comment, comment.Id)
	if err != nil {
		return err
	}
	return nil
}

func DeleteComment(db *sql.DB, id int) error {
	query := "DELETE FROM comments WHERE id = ?"
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
