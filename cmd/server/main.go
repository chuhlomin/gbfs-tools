package main

import (
	"log"
	"net/http"

	"github.com/caarlos0/env/v6"
	"github.com/graphql-go/handler"

	"github.com/chuhlomin/gbfs-tools/pkg/gbfs"
)

type config struct {
	Hostname     string `env:"HOSTNAME" envDefault:"127.0.0.1"`
	Port         string `env:"PORT" envDefault:"8082"`
	DatabasePath string `env:"DB_PATH,required"`
}

func main() {
	log.Println("Parsing environment variables...")
	var c config
	if err := env.Parse(&c); err != nil {
		log.Fatalf("ERROR: Failed to parse environment variables: %v", err)
	}

	h := handler.New(&handler.Config{
		Schema:     &gbfs.Schema,
		Pretty:     true,
		GraphiQL:   true,
		Playground: true,
	})

	// all := nestedMiddleware(withCORS)
	http.HandleFunc("/graphql", withCORS(h))

	bind := c.Hostname + ":" + c.Port
	log.Printf("Listening on %v", bind)
	log.Fatal(http.ListenAndServe(bind, nil))
}

func withCORS(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		next.ServeHTTP(w, r)
	}
}
