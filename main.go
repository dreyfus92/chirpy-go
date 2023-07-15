package main

import (
	"log"
	"net/http"

	"github.com/dreyfus92/chirpy-go/internal/database"
	"github.com/go-chi/chi/v5"
)

func main() {
	filepathRoot := "."
	port := "8080"

	apiCfg := &apiConfig{
		fileserverHits: 0,
	}

	_, err := database.NewDB(filepathRoot + "/database.json")

	if err != nil {
		log.Fatalf("Error creating database: %v", err)
	}

	r := chi.NewRouter()
	api := chi.NewRouter()
	admin := chi.NewRouter()

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))

	r.Handle("/app/*", fsHandler)
	r.Handle("/app", fsHandler)

	api.Get("/healthz", handlerReadiness)
	api.Post("/chirps", handlerValidateChirp)

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
