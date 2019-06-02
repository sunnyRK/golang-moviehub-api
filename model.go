package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

//////////////////////////////////////////////////////////////
/////////////////// MOVIE STRUCTURE //////////////////////////
//////////////////////////////////////////////////////////////


type movie struct {
	Id int `json:"id"`
	Title string `json:"title"`
	Year  string `json:"year"`
	Rated string `json:"rated"`
	Genre string `json:"genre"`
}

/////////////////////////////////////////////////////////////
/////////////////// DATABASE OPERATIONS /////////////////////
/////////////////////////////////////////////////////////////

func (m *movie) createMovie(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO movie(title, year, rated, genre) VALUES($1, $2, $3, $4) RETURNING id",
		m.Title, m.Year, m.Rated, m.Genre).Scan(&m.Id)

	if err != nil {
		return err
	}
	return nil
}

func (m *movie) getMovieById(db *sql.DB) error {
	return db.QueryRow("SELECT title, year, rated, genre FROM movie WHERE id=$1",m.Id).Scan(&m.Title, &m.Year, &m.Rated, &m.Genre)
}

func (m *movie) getMovieByReleasedYear(db *sql.DB) error {
	return db.QueryRow("SELECT title, year, rated, genre FROM movie WHERE year=$1",m.Year).Scan(&m.Title, &m.Year, &m.Rated, &m.Genre)
}

func (m *movie) getMovieByTitle(db *sql.DB) error {
	var count int

	row := db.QueryRow("SELECT COUNT(*) FROM movie WHERE title=$1",m.Title)
	err := row.Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(count)

	if count == 0 {
		api := Init("a39d4a2a")
		query := &QueryData{Title: m.Title, Year: ""}
		res2, err := api.MovieByTitle(query)
		if err != nil {
			fmt.Println(err)
		}

		var m movie		
		m.Title = res2.Title
		m.Year = res2.Year
		m.Rated = res2.Rated
		m.Genre = res2.Genre

		if err := m.createMovie(db); err != nil {
			fmt.Println(err)
		}

		fmt.Println(res2.Title)
		fmt.Println(res2.Year)
		fmt.Println(res2.Rated)
		fmt.Println(res2.Genre)
	}

	return db.QueryRow("SELECT title, year, rated, genre FROM movie WHERE title=$1",m.Title).Scan(&m.Title, &m.Year, &m.Rated, &m.Genre)
}

func (m *movie) updateMovie(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE movie SET genre=$1, rated=$2 WHERE id=$3",
			m.Genre, m.Rated, m.Id)

	return err
}
