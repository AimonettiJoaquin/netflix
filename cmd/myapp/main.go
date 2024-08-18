package main

import (
	"log"
	"net/http"
	"netflix/internal/config"
	"netflix/internal/database"
	"netflix/internal/handlers"

	"github.com/gorilla/mux"
)

func main() {
	// Lectura del archivo confg.yaml -configuración

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error al cargar la configuración")
	}

	// Conectamos a la DB

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}
	defer db.Close()

	// Creamos la tabla usuario

	if err := database.CreateUsersTable(db); err != nil {
		log.Fatalf("Error al crear la tabla usuario: %v", err)
	}
	if err := database.CreateCommentsTable(db); err != nil {
		log.Fatalf("Error al crear la tabla comments: %v", err)
	}
	if err := database.CreateMovieTable(db); err != nil {
		log.Fatalf("Error al crear la tabla movie: %v", err)
	}

	//Instanciamos el router de gorilla mux
	router := mux.NewRouter()
	handlers.UserRouterHandlers(router, db)
	handlers.MovieRouterHandlers(router, db)
	handlers.CommentRouterHandler(router, db)

	// Iniciar el servidor HTTP

	// Iniciar el servidor
	log.Printf("Servidor corriendo en %s", cfg.ServerAddress)
	if err := http.ListenAndServe(cfg.ServerAddress, router); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}

}
