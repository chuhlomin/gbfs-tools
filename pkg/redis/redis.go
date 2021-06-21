package redis

import (
	"context"
	"fmt"
	"strings"

	"github.com/mediocregopher/radix/v4"
	"github.com/pkg/errors"

	"github.com/chuhlomin/gbfs-go"
	"github.com/chuhlomin/gbfs-tools/pkg/structs"
)

// Client represents layer between server, writer and Redis
type Client struct {
	ctx    context.Context
	client radix.Client
}

// NewClient creates new Client
func NewClient(ctx context.Context, network, addr, auth string) (*Client, error) {
	dialer := radix.Dialer{}
	if auth != "" {
		dialer = radix.Dialer{AuthPass: auth}
	}

	client, err := (radix.PoolConfig{Dialer: dialer}).New(ctx, network, addr)
	if err != nil {
		return nil, err
	}

	var pong string
	err = client.Do(ctx, radix.Cmd(&pong, "PING"))
	if err != nil {
		return nil, err
	}
	if pong != "PONG" {
		return nil, fmt.Errorf("PING failed, got %q", pong)
	}

	return &Client{
		ctx:    ctx,
		client: client,
	}, nil
}

func (c *Client) WriteSystem(system gbfs.System) error {
	return c.client.Do(
		c.ctx,
		radix.Cmd(
			nil,
			"HMSET",
			fmt.Sprintf("system:%s", system.ID),
			"id", system.ID,
			"name", system.Name,
			"url", system.AutoDiscoveryURL,
			"www", system.URL,
			"cc", system.CountryCode,
			"location", system.Location,
		),
	)
}

func (c *Client) GetSystem(systemID string) (*structs.System, error) {
	var vv []string
	err := c.client.Do(c.ctx, radix.Cmd(&vv, "HGETALL", fmt.Sprintf("system:%s", systemID)))
	if err != nil {
		return nil, errors.Wrapf(err, "get system %q", systemID)
	}

	system := structs.System{}

	for i := 0; i < len(vv); i++ {
		key := vv[i]
		value := vv[i+1]

		switch key {
		case "id":
			system.ID = value
		case "name":
			system.Name = value
		case "url":
			system.AutoDiscoveryURL = value
		case "www":
			system.URL = value
		case "cc":
			system.CountryCode = value
		case "location":
			system.Location = value
		}
		i++
	}

	return &system, nil
}

func (c *Client) GetSystemURL(systemID string) (string, error) {
	var v string
	err := c.client.Do(c.ctx, radix.Cmd(&v, "HGET", fmt.Sprintf("system:%s", systemID), "url"))
	if err != nil {
		return "", errors.Wrapf(err, "get system URL %q", systemID)
	}

	return v, nil
}

func (c *Client) GetSystemsIDs() ([]string, error) {
	var result, keys []string

	err := c.client.Do(c.ctx, radix.Cmd(&keys, "KEYS", "system:*"))
	if err != nil {
		return nil, errors.Wrap(err, "get systems keys")
	}

	for _, key := range keys {
		result = append(result, key[len("system:"):])
	}

	return result, nil
}

func (c *Client) WriteFeeds(systemID, language string, feeds []gbfs.Feed) error {
	for _, feed := range feeds {
		err := c.client.Do(
			c.ctx,
			radix.Cmd(
				nil,
				"SET",
				fmt.Sprintf("feed:%s:%s:%s", systemID, feed.Name, language),
				feed.URL,
			),
		)
		if err != nil {
			return errors.Wrapf(err, "write feed %q: %q", feed.Name, feed.URL)
		}
	}
	return nil
}

func (c *Client) GetFeedURL(systemID, feedName, language string) (string, error) {
	var url string
	err := c.client.Do(c.ctx, radix.Cmd(&url, "GET", fmt.Sprintf("feed:%s:%s:%s", systemID, feedName, language)))
	if err != nil {
		return "", errors.Wrap(err, "GET")
	}

	if url == "" { // fallback to any other language
		var keys []string
		err = c.client.Do(c.ctx, radix.Cmd(&keys, "KEYS", fmt.Sprintf("feed:%s:%s:*", systemID, feedName)))
		if err != nil {
			return "", errors.Wrap(err, "KEYS")
		}

		if len(keys) > 0 {
			err = c.client.Do(c.ctx, radix.Cmd(&url, "GET", keys[0]))
		}
	}

	return url, err
}

func (c *Client) GetFeeds(systemID string) ([]structs.Feed, error) {
	var keys []string
	if err := c.client.Do(c.ctx, radix.Cmd(&keys, "KEYS", fmt.Sprintf("feed:%s:*", systemID))); err != nil {
		return nil, errors.Wrapf(err, "keys for %q feeds", systemID)
	}

	result := []structs.Feed{}

	for _, key := range keys {
		_, feedName, language := splitFeedKey(key)

		var url string
		err := c.client.Do(c.ctx, radix.Cmd(&url, "GET", key))
		if err != nil {
			return nil, err
		}

		result = append(
			result,
			structs.Feed{
				Name:     feedName,
				URL:      url,
				Language: language,
			},
		)
	}

	return result, nil
}

func (c *Client) GetFeedsLanguages(systemID string) ([]string, error) {
	var keys []string
	if err := c.client.Do(c.ctx, radix.Cmd(&keys, "KEYS", fmt.Sprintf("feed:%s:*", systemID))); err != nil {
		return nil, errors.Wrapf(err, "keys for %q feeds", systemID)
	}

	langs := map[string]struct{}{}
	for _, key := range keys {
		_, _, language := splitFeedKey(key)
		langs[language] = struct{}{}
	}

	var result []string
	for lang := range langs {
		result = append(result, lang)
	}

	return result, nil
}

func splitFeedKey(key string) (systemID, feedName, language string) {
	v := strings.Split(key, ":")
	systemID, feedName, language = v[0], v[1], v[2]
	return
}
