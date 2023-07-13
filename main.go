package main

import (
	"fmt"
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func handlerMetrics(cfg *apiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits)))
	}
}

func main(){
	filepathRoot := "./"
	port := "8080"
	cfg := &apiConfig{}

	//initializing mux
	mux := http.NewServeMux()

	// adding app handler
	mux.Handle("/app/", cfg.middlewareMetricsInc((http.StripPrefix("/app", http.FileServer(http.Dir("."))))))
	//adding healthz handler
	mux.HandleFunc("/healthz", handlerReadiness)
	//adding metrics handler
	mux.HandleFunc("/metrics", handlerMetrics(cfg))
	//applying middleware
	corsMux := middlewareCors(mux)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Server started on port %s", port)
	log.Printf("Serving files from %s\n", filepathRoot)

	log.Fatal(srv.ListenAndServe())

}

func handlerReadiness(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

