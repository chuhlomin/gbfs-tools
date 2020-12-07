package main

import (
	"log"
	"net/http"

	"github.com/caarlos0/env/v6"
	"github.com/graphql-go/handler"

	"github.com/chuhlomin/gbfs-tools/pkg/gbfs"
)

type config struct {
	Hostname string `env:"HOSTNAME" envDefault:"127.0.0.1"`
	Port     string `env:"PORT" envDefault:"8080"`
}

func main() {
	log.Println("Parsing environment variables...")
	var c config
	if err := env.Parse(&c); err != nil {
		log.Fatalf("ERROR: Failed to parse environment variables: %v", err)
	}

	h := handler.New(&handler.Config{
		Schema:   &gbfs.Schema,
		Pretty:   true,
		GraphiQL: true,
	})

	http.Handle("/graphql", h)

	bind := c.Hostname + ":" + c.Port
	log.Printf("Listening on %v", bind)
	log.Fatal(http.ListenAndServe(bind, nil))
}
