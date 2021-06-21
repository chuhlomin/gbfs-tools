package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/caarlos0/env/v6"

	"github.com/chuhlomin/gbfs-go"
	"github.com/chuhlomin/gbfs-tools/pkg/redis"
	"github.com/pkg/errors"
)

type config struct {
	SystemsURL   string        `env:"SYSTEMS_CSV_URL" envDefault:"https://raw.githubusercontent.com/NABSA/gbfs/master/systems.csv"`
	RedisNetwork string        `env:"REDIS_NETWORK" envDefault:"tcp"`
	RedisAddr    string        `env:"REDIS_ADDR" envDefault:"redis:6379"`
	RedisAuth    string        `env:"REDIS_AUTH"`
	FeedsDelay   time.Duration `env:"FEEDS_DELAY" envDefault:"2s"`
}

func main() {
	log.Print("Starting...")
	if err := run(); err != nil {
		log.Fatalf("ERROR %v", err)
	}
	log.Print("Finished")
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

	gbfsClient := gbfs.NewClient("github.com/chuhlomin/gbfs-tools/writer", 30*time.Second)

	resp, err := http.Get(c.SystemsURL)
	if err != nil {
		return errors.Wrapf(err, "get %q", c.SystemsURL)
	}
	defer resp.Body.Close()

	systems, err := gbfs.ParseSystemsCSV(resp.Body)
	if err != nil {
		return errors.Wrap(err, "parse systems CSV")
	}

	log.Print("Writing systems...")
	if err := writeSystems(systems, redisClient); err != nil {
		return errors.Wrap(err, "write systems")
	}

	log.Print("Writing feeds...")
	if err := writeFeeds(systems, redisClient, gbfsClient, c.FeedsDelay); err != nil {
		return errors.Wrap(err, "write feeds")
	}

	return nil
}

func writeSystems(systems []gbfs.System, redisClient *redis.Client) error {
	for _, system := range systems {
		if err := redisClient.WriteSystem(system); err != nil {
			return errors.Wrap(err, system.ID)
		}
	}

	return nil
}

func writeFeeds(
	systems []gbfs.System,
	redisClient *redis.Client,
	client *gbfs.Client,
	delay time.Duration,
) error {
	for _, system := range systems {
		resp, err := client.LoadGBFS(system.AutoDiscoveryURL)
		if err != nil {
			log.Printf("Failed to load GBFS for %q by URL %q: %v", system.ID, system.URL, err)
			continue
		}

		for lang, feeds := range resp.Data {
			if err := redisClient.WriteFeeds(system.ID, lang, feeds.Feeds); err != nil {
				return errors.Wrapf(err, "%s %s", system.ID, lang)
			}
		}

		time.Sleep(delay)
	}
	return nil
}
