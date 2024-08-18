package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"netflix/pkg/model"
	"strconv"

	"github.com/gorilla/mux"
)

func CommentRouterHandler(router *mux.Router, db *sql.DB) {
	router.HandleFunc("/comments", createComment(db)).Methods("POST")
	router.HandleFunc("/comments/{id}", getCommentByID(db)).Methods("GET")
	router.HandleFunc("/comments", getComments(db)).Methods("GET")
	router.HandleFunc("/comments/{id}", deleteComment(db)).Methods("DELETE")
	router.HandleFunc("/comments/{id}", updateComment(db)).Methods("PUT")
}

func createComment(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var comment model.Comment
		if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := model.CreateComment(db, &comment); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(comment)
	}
}

func getCommentByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		comment, err := model.GetCommentByID(db, id)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Comment not found", http.StatusNotFound)
				return
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(comment)
	}
}

func getComments(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		comments, err := model.GetComments(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(comments)
	}
}

func deleteComment(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = model.GetCommentByID(db, id)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Comment not found", http.StatusNotFound)
				return
			}
		}

		err = model.DeleteComment(db, id)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Comment not found", http.StatusNotFound)
				return
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		type response struct {
			Message string `json:"message"`
		}
		res := response{
			Message: "Comment deleted successfully",
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res.Message)
	}
}

func updateComment(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid comment ID", http.StatusBadRequest)
			return
		}

		_, err = model.GetCommentByID(db, id)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Comment not found", http.StatusNotFound)
				return
			}
		}

		var comment model.Comment
		if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		comment.Id = id
		if err := model.UpdateComment(db, &comment); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(comment)
	}
}
