package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"

	g "github.com/chuhlomin/gbfs-go"

	"github.com/chuhlomin/gbfs-tools/pkg/gbfs"
	"github.com/chuhlomin/gbfs-tools/pkg/redis"
)

type config struct {
	Hostname     string `env:"HOSTNAME" envDefault:"127.0.0.1"`
	Port         string `env:"PORT" envDefault:"8082"`
	AllowOrigin  string `env:"CORS_ALLOW_ORIGIN" envDefault:"*"`
	RedisNetwork string `env:"REDIS_NETWORK" envDefault:"tcp"`
	RedisAddr    string `env:"REDIS_ADDR" envDefault:"redis:6379"`
	RedisAuth    string `env:"REDIS_AUTH"`
}

func main() {
	log.Print("Stating...")
	if err := run(); err != nil {
		log.Fatalf("ERROR %v", err)
	}
}

func run() error {
	log.Println("Parsing environment variables...")
	var c config
	if err := env.Parse(&c); err != nil {
		return errors.Wrap(err, "parse environment variables")
	}

	redisClient, err := redis.NewClient(
		context.Background(),
		c.RedisNetwork,
		c.RedisAddr,
		c.RedisAuth,
	)
	if err != nil {
		return errors.Wrap(err, "create Redis client")
	}

	gbfs.Client = g.NewClient("github.com/chuhlomin/gbfs-tools", 30*time.Second)
	gbfs.RedisClient = redisClient

	http.HandleFunc("/", ok)
	http.HandleFunc("/graphql", withLogging(withCORS(gbfs.Handler(), c.AllowOrigin)))
	http.HandleFunc("/geojson", withLogging(withCORS(gbfs.HandlerGeoJSON(), c.AllowOrigin)))

	bind := c.Hostname + ":" + c.Port
	log.Printf("Listening on %v", bind)
	return http.ListenAndServe(bind, nil)
}

func ok(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("OK"))
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
