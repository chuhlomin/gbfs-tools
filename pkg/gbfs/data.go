package gbfs

import (
	"log"
	"time"

	"github.com/chuhlomin/gbfs-go"
)

var client *gbfs.Client

func init() {
	client = gbfs.NewClient("github.com/chuhlomin/gbfs-tools", 30*time.Second)
}

func GetSystems() ([]gbfs.System, error) {
	log.Println("GET Systems")
	return client.LoadSystems(gbfs.SystemsNABSA)
}
