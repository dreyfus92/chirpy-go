package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dreyfus92/chirpy-go/internal/database"
	"github.com/go-chi/chi/v5"
)

type apiConfig struct {
	fileserverHits int
	db             *database.DB
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	const databaseFile = "database.json"

	//flag to delete the database file when debugging
	dbg := flag.Bool("debug", false, "enable debug mode")
	flag.Parse()

	if *dbg {
		fmt.Printf("Debug mode enabled, deleting database file \n")
		os.Remove(filepathRoot + "/" + databaseFile)
	}

	//create the DB
	db, err := database.NewDB(filepathRoot + "/" + databaseFile)
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := &apiConfig{
		fileserverHits: 0,
		db:             db,
	}

	r := chi.NewRouter()
	api := chi.NewRouter()
	admin := chi.NewRouter()

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))

	r.Handle("/app/*", fsHandler)
	r.Handle("/app", fsHandler)

	api.Get("/healthz", readinessHandler)

	// create new chirps
	api.Post("/chirps", apiCfg.createChirpHandler)

	// read all chirps
	api.Get("/chirps", apiCfg.getChirpsHandler)
	// read a single chirp
	api.Get("/chirps/{id}", apiCfg.getChirpHandler)

	// create new users
	api.Post("/users", apiCfg.createUserHandler)

	admin.Get("/metrics", apiCfg.handlerMetrics)

	r.Mount("/api", api)
	r.Mount("/admin", admin)

	corsMux := middlewareCors(r)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Server started on port %s", port)
	log.Printf("Serving files from %s\n", filepathRoot)

	log.Fatal(srv.ListenAndServe())

}
