package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main(){
	filepathRoot := "."
	port := "8080"

	apiCfg := &apiConfig{
		fileserverHits: 0,
	}

	// mux := http.NewServeMux()
	r := chi.NewRouter()

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))

	r.Handle("/app/*", fsHandler)
	r.Handle("/app", fsHandler)

	r.Get("/healthz", handlerReadiness)
	r.Get("/metrics", apiCfg.handlerMetrics)

	corsMux := middlewareCors(r)


	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Server started on port %s", port)
	log.Printf("Serving files from %s\n", filepathRoot)

	log.Fatal(srv.ListenAndServe())

}


