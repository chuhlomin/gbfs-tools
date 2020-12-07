package main

import (
	"log"
	"net/http"

	"github.com/graphql-go/handler"

	"github.com/chuhlomin/gbfs-tools/pkg/gbfs"
)

func main() {
	h := handler.New(&handler.Config{
		Schema:   &gbfs.Schema,
		Pretty:   true,
		GraphiQL: true,
	})

	http.Handle("/graphql", h)

	bind := "127.0.0.1:8080"
	log.Printf("Listening on %v", bind)
	log.Fatal(http.ListenAndServe(bind, nil))
}
