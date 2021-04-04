package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-openapi/runtime/middleware"
	goHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/grkmk/glm-api/handlers"
)

func main() {
	// create logger
	logger := log.New(os.Stdout, "[product-api]: ", log.LstdFlags)

	// create handlers
	productsHandler := handlers.NewProducts(logger)

	// create serve mux & register handlers

	// serveMux := http.NewServeMux()
	// serveMux.Handle("/", productsHandler)
	// serveMux.Handle("/products", productsHandler)
	serveMux := mux.NewRouter()

	getRouter := serveMux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/products", productsHandler.GetProducts)
	getRouter.HandleFunc("/products/{id:[0-9]+}", productsHandler.GetProduct)

	postRouter := serveMux.Methods(http.MethodPost).Subrouter()
	postRouter.Use(productsHandler.ValidateProduct)
	postRouter.HandleFunc("/products", productsHandler.AddProduct)

	patchRouter := serveMux.Methods(http.MethodPatch).Subrouter()
	patchRouter.Use(productsHandler.ValidateProduct)
	patchRouter.HandleFunc("/products/{id:[0-9]+}", productsHandler.UpdateProducts)

	deleteRouter := serveMux.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/products/{id:[0-9]+}", productsHandler.DeleteProduct)

	specialHandlerOptions := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	specialHandler := middleware.Redoc(specialHandlerOptions, nil)
	getRouter.Handle("/docs", specialHandler)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	// CORS
	corsHandler := goHandlers.CORS(goHandlers.AllowedOrigins([]string{"http://localhost:3000"}))

	// create server
	httpServer := &http.Server{
		Addr:         ":9090",               // configure bind addrees
		Handler:      corsHandler(serveMux), // set default handler
		ReadTimeout:  1 * time.Second,       // max time to read request from client
		WriteTimeout: 1 * time.Second,       // max time to write response to client
		IdleTimeout:  120 * time.Second,     // max time for connections using TCP Keep-Alive
	}

	// start server
	go func() {
		logger.Println("Starting server on port 9090")

		err := httpServer.ListenAndServe()
		if err != nil {
			logger.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interrupt & gracefully shutdown server
	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt)
	signal.Notify(signalChannel, os.Kill)

	sig := <-signalChannel
	logger.Println("Received terminate, gracefully shutting down", sig)

	timeoutContext, _ := context.WithTimeout(context.Background(), 30*time.Second)
	httpServer.Shutdown(timeoutContext)

	http.ListenAndServe(":9090", serveMux) // creates a web server
}
