package main

import (
	"log"
	"net/http"

	"github.com/dreyfus92/chirpy-go/internal/database"
	"github.com/go-chi/chi/v5"
)

type apiConfig struct {
	fileserverHits int
	db             *database.DB
}

func main() {
	filepathRoot := "."
	port := "8080"
	databaseFile := "database.json"

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
