package model

import "database/sql"

type Movie struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Overview string `json:"overview"`
	Counter  int    `json:"counter"`
}
type MovieCounter struct {
	ID      int `json:"id"`
	Counter int `json:"counter"`
}

type Movies struct {
	Results []Movie
}

func GetMovieByID(db *sql.DB, id int) (*MovieCounter, error) {
	var movieC MovieCounter
	err := db.QueryRow("SELECT id, counter FROM movies WHERE id = ?", id).Scan(&movieC.ID, &movieC.Counter)
	if err != nil {
		return nil, err
	}
	return &movieC, nil
}

func CreateMovie(db *sql.DB, id int) error {
	query := "INSERT INTO movies (id, counter) VALUES (?,1)"
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func UpdateMovie(db *sql.DB, id int) error {
	query := "UPDATE movies SET counter = counter + 1 WHERE id = ?"
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
