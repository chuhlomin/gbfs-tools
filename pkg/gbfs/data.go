package gbfs

import (
	"log"
	"time"

	"github.com/chuhlomin/gbfs-go"
)

var client *gbfs.Client

var systems []gbfs.System

func init() {
	client = gbfs.NewClient("github.com/chuhlomin/gbfs-tools", 30*time.Second)
}

func GetSystems() ([]gbfs.System, error) {
	var err error
	if systems == nil {
		log.Println("GET Systems")
		systems, err = client.LoadSystems(gbfs.SystemsNABSA)
	}

	return systems, err
}
