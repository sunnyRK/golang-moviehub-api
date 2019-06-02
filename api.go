package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB 	   *sql.DB
}

//////////////////////////////////////////////////////////////
/////////////////// POSTGRES CREDENTIALS /////////////////////
//////////////////////////////////////////////////////////////

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "sunny"
	dbname   = "imdb"
)

////////////////////////////////////////////////
/////////////////// ROUTES /////////////////////
///////////////////////////////////////////////

func (a *App) Initialize() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
	host, port, user, password, dbname)
	var err error
	a.DB, err = sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(":8000", a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/movie", a.createMovie).Methods("POST")
	a.Router.HandleFunc("/movieById/{id}", a.getMovieById).Methods("GET")
	a.Router.HandleFunc("/movieByTitle/{title}", a.getMovieByTitle).Methods("GET")
	a.Router.HandleFunc("/movieByYear/{year}", a.getMovieByReleasedYear).Methods("GET")
	a.Router.HandleFunc("/movie/{id}", a.updateMovie).Methods("PUT")
}

///////////////////////////////////////////////////
/////////////////// REST APIS /////////////////////
///////////////////////////////////////////////////

func (a *App) createMovie(w http.ResponseWriter, r *http.Request) {
	var m movie
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := m.createMovie(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, m)
}

func (a *App) getMovieById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid movie ID")
		return
	}

	m := movie{Id: id}
	if err := m.getMovieById(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Movie not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, m)
}

func (a *App) getMovieByReleasedYear(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	year := vars["year"]

	m := movie{Year: year}
	if err := m.getMovieByReleasedYear(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Movie not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, m)
}

func (a *App) getMovieByTitle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title  := vars["title"]

	m := movie{Title: title}
	if err := m.getMovieByTitle(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Movie not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, m)
}

func (a *App) updateMovie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid movie ID")
		return
	}

	var m movie
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()
	m.Id = id

	if err := m.updateMovie(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, m)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

///////////////////////////////////////////////////////
/////////////////// MAIN FUNCTION /////////////////////
///////////////////////////////////////////////////////

func main() {
	a := App{}
	a.Initialize()
	a.Run(":8000")
}	