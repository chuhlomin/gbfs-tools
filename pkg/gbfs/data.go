package gbfs

import (
	"log"
	"time"

	"github.com/chuhlomin/gbfs-go"
)

var client *gbfs.Client

var systems []gbfs.System
var feeds map[string]gbfs.LanguageFeeds // url → gbfs

func init() {
	client = gbfs.NewClient("github.com/chuhlomin/gbfs-tools", 30*time.Second)
	feeds = map[string]gbfs.LanguageFeeds{}
}

func GetSystems(url string) ([]gbfs.System, error) {
	var err error
	if systems == nil {
		log.Printf("GET Systems %q", url)
		systems, err = client.LoadSystems(url)
	}

	return systems, err
}

func GetGBFS(url string) (gbfs.LanguageFeeds, error) {
	var err error
	if f, ok := feeds[url]; ok {
		return f, nil
	}

	log.Printf("GET GBFS %q", url)
	resp, err := client.LoadGBFS(url)
	if err != nil {
		return nil, err
	}

	feeds[url] = resp.Data
	return resp.Data, nil
}
