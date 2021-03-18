package main

import (
	"log"
	"net/http"

	"github.com/caarlos0/env/v6"

	"github.com/chuhlomin/gbfs-tools/pkg/gbfs"
)

type config struct {
	Hostname     string `env:"HOSTNAME" envDefault:"127.0.0.1"`
	Port         string `env:"PORT" envDefault:"8082"`
	AllowOrigin  string `env:"CORS_ALLOW_ORIGIN" envDefault:"*"`
	DatabasePath string `env:"DB_PATH,required"` // used in pkg/gbfs/data.go
}

func main() {
	log.Println("Parsing environment variables...")
	var c config
	if err := env.Parse(&c); err != nil {
		log.Fatalf("ERROR: Failed to parse environment variables: %v", err)
	}

	http.HandleFunc("/graphql", withLogging(withCORS(gbfs.Handler(), c.AllowOrigin)))
	http.HandleFunc("/geojson", withLogging(withCORS(gbfs.HandlerGeoJSON(), c.AllowOrigin)))

	bind := c.Hostname + ":" + c.Port
	log.Printf("Listening on %v", bind)
	log.Fatal(http.ListenAndServe(bind, nil))
}

func withCORS(next http.Handler, allowOrigin string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", allowOrigin)
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		next.ServeHTTP(w, r)
	}
}

func withLogging(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s %s", r.Method, r.URL, r.RemoteAddr, r.UserAgent())
		next.ServeHTTP(w, r)
	}
}
