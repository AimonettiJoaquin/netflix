package handlers

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"netflix/internal/config"
	"netflix/pkg/model"
	models "netflix/pkg/model"
	"strconv"

	"github.com/gorilla/mux"
)

func MovieRouterHandlers(router *mux.Router, db *sql.DB) {
	router.HandleFunc("/movies/popular/{cant}", getListMoviePopular()).Methods("GET")
	router.HandleFunc("/movies/{id}", getMovieById(db)).Methods("GET")

}

func getMovieById(db *sql.DB) http.HandlerFunc {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error al cargar la configuración")
	}
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])

		if err != nil {
			http.Error(w, "Invalid movie ID", http.StatusBadRequest)
			return
		}
		idUser, err := strconv.Atoi(r.Header.Get("id_user"))
		if err != nil {
			http.Error(w, "Please login", http.StatusInternalServerError)
			return
		} else if idUser == 0 {
			http.Error(w, "Please login", http.StatusInternalServerError)
			return
		}

		url := "https://api.themoviedb.org/3/movie/" + strconv.Itoa(id) + "?language=en-US"
		req, _ := http.NewRequest("GET", url, nil)

		req.Header.Add("accept", "application/json")
		req.Header.Add("Authorization", cfg.AUTH)
		res, _ := http.DefaultClient.Do(req)
		if res.StatusCode != 200 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("ID invalido, prueba con otro"))
			return
		}

		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var movie models.Movie

		err = json.Unmarshal(body, &movie)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		comment, err := model.GetCommentByMovieUser(db, id, idUser)
		if err != nil {
			http.Error(w, "This user have no comment in this movie", http.StatusBadRequest)
			return
		}
		var movieComment models.MovieComment
		movieComment.Movie = movie
		movieComment.Comment = comment.Comment

		_, err = model.GetMovieByID(db, id)
		if err != nil {
			model.CreateMovie(db, id)
		} else {
			model.UpdateMovie(db, id)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(movieComment)
	}
}

func getListMoviePopular() http.HandlerFunc {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error al cargar la configuración")
	}
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		cant, err := strconv.Atoi(vars["cant"])

		if err != nil {
			http.Error(w, "Invalid quantity", http.StatusBadRequest)
			return
		} else if cant > 35 {
			http.Error(w, "35 results max", http.StatusBadRequest)
			return
		} else if cant < 1 {
			http.Error(w, "Invalid quantity", http.StatusBadRequest)
			return
		}

		url := "https://api.themoviedb.org/3/movie/popular?language=en-US"
		req, _ := http.NewRequest("GET", url, nil)

		req.Header.Add("accept", "application/json")
		req.Header.Add("Authorization", cfg.AUTH)
		res, _ := http.DefaultClient.Do(req)
		if res.StatusCode != 200 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error en la solicitud"))
			return
		}

		defer res.Body.Close()

		// Decodificar la respuesta Json
		var movies models.Movies
		if err := json.NewDecoder(res.Body).Decode(&movies); err != nil {
			log.Fatalf("Error al decodificar la respuesta json: %v", err)
		}
		//seleccionamos los primeros n resultados
		topMovies := movies.Results[:cant]

		//Configurar el encabezado de respuesta como JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		//Codificar la respuesta en JSON formateado y escribirla
		jsonResponse, err := json.MarshalIndent(topMovies, "", "  ")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(jsonResponse)

	}
}
