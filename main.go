package main

import (
	"log"
	"net/http"
)

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


