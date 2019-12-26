package main

import (
	"log"
	"net/http"
	"os"

	"github.com/sp98/tickstore/pkg/apis/v1/basicauth"
	"github.com/sp98/tickstore/pkg/apis/v1/stocks"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/sp98/tickstore/pkg/apis/v1/ohlc"
)

var (
	userName string
	password string
)

func init() {
	userName = os.Getenv("API_USER_NAME")
	password = os.Getenv("API_PASSWORD")
	if userName == "" && password == "" {
		log.Fatal("failed to read the API username and password")
		panic(1)
	}
}

//Routes define all the global routes.
func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON), // Set content-Type headers as application/json
		middleware.Logger,          // Log API request calls
		middleware.DefaultCompress, // Compress results, mostly gzipping assets and json
		middleware.RedirectSlashes, // Redirect slashes to no slash URL versions
		middleware.Recoverer,       // Recover from panics without crashing server
		basicauth.New("TICKSTORE", map[string][]string{
			userName: {password},
		}),
	)

	router.Route("/v1", func(r chi.Router) {
		r.Mount("/api/ohlc", ohlc.Routes())
		r.Mount("/api/stocks", stocks.Routes())
	})

	return router
}

func main() {
	router := Routes()

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("%s %s\n", method, route) // Walk and print out all routes
		return nil
	}
	if err := chi.Walk(router, walkFunc); err != nil {
		log.Panicf("Logging err: %s\n", err.Error()) // panic if there is an error
	}

	log.Fatal(http.ListenAndServe(":3001", router)) // Note, the port is usually gotten from the environment.
}
