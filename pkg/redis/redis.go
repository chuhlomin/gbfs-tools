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
	return c.client.Do(c.ctx, radix.Cmd(nil, "SET", "system:"+system.ID, packSystem(system)))
}

var systems []*structs.System

func (c *Client) CacheAllSystems() error {
	var keys []string
	if err := c.client.Do(c.ctx, radix.Cmd(&keys, "KEYS", "system:*")); err != nil {
		return errors.Wrap(err, "get systems keys")
	}

	var vals []string
	if err := c.client.Do(c.ctx, radix.Cmd(&vals, "MGET", keys...)); err != nil {
		return errors.Wrap(err, "get systems keys")
	}

	systems = []*structs.System{}
	for _, val := range vals {
		systems = append(systems, unpackSystem(val))
	}

	return nil
}

func (c *Client) GetSystems() ([]*structs.System, error) {
	if err := c.CacheAllSystems(); err != nil {
		return nil, err
	}

	return systems, nil
}

func (c *Client) GetSystem(systemID string) (*structs.System, error) {
	var v string
	if err := c.client.Do(c.ctx, radix.Cmd(&v, "GET", "system:"+systemID)); err != nil {
		return nil, errors.Wrapf(err, "get system %q", systemID)
	}

	return unpackSystem(v), nil
}

func packSystem(system gbfs.System) string {
	return strings.Join(
		[]string{
			system.ID,
			system.Name,
			system.AutoDiscoveryURL,
			system.URL,
			system.CountryCode,
			system.Location,
		},
		"\n",
	)
}

func unpackSystem(val string) *structs.System {
	vv := strings.Split(val, "\n")
	return &structs.System{
		ID:               vv[0],
		Name:             vv[1],
		AutoDiscoveryURL: vv[2],
		URL:              vv[3],
		CountryCode:      vv[4],
		Location:         vv[5],
	}
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

var allFeeds map[string][]structs.Feed

func (c *Client) CacheAllFeeds() error {
	var keys []string
	if err := c.client.Do(c.ctx, radix.Cmd(&keys, "KEYS", "feed:*")); err != nil {
		return errors.Wrapf(err, "keys for all feeds")
	}

	if len(keys) == 0 {
		return nil
	}

	allFeeds = map[string][]structs.Feed{}

	var urls []string
	if err := c.client.Do(c.ctx, radix.Cmd(&urls, "MGET", keys...)); err != nil {
		return errors.Wrap(err, "mget for all feeds")
	}

	for i, key := range keys {
		systemID, feedName, language := splitFeedKey(key)

		feed := structs.Feed{
			Name:     feedName,
			URL:      urls[i],
			Language: language,
		}

		if _, ok := allFeeds[systemID]; !ok {
			allFeeds[systemID] = []structs.Feed{feed}
		} else {
			allFeeds[systemID] = append(allFeeds[systemID], feed)
		}
	}

	return nil
}

func (c *Client) GetFeeds(systemID string) ([]structs.Feed, error) {
	if feeds, ok := allFeeds[systemID]; ok {
		return feeds, nil
	}

	var keys []string
	if err := c.client.Do(c.ctx, radix.Cmd(&keys, "KEYS", fmt.Sprintf("feed:%s:*", systemID))); err != nil {
		return nil, errors.Wrapf(err, "keys for %q feeds", systemID)
	}

	if len(keys) == 0 {
		return nil, nil
	}

	result := []structs.Feed{}

	var urls []string
	if err := c.client.Do(c.ctx, radix.Cmd(&urls, "MGET", keys...)); err != nil {
		return nil, errors.Wrapf(err, "mget for %q feeds", systemID)
	}

	for i, key := range keys {
		_, feedName, language := splitFeedKey(key)

		result = append(
			result,
			structs.Feed{
				Name:     feedName,
				URL:      urls[i],
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
	systemID, feedName, language = v[1], v[2], v[3]
	return
}
