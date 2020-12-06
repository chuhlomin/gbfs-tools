package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/chuhlomin/gbfs-go"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		log.Printf("ERROR: Failed to %v", err)
	}
}

func run() error {
	client := gbfs.NewClient("github.com/chuhlomin/gbfs-tools", 30*time.Second)

	systems, err := client.LoadSystems(gbfs.SystemsNABSA)
	if err != nil {
		return errors.Wrap(err, "load systems")
	}

	_ = []string{
		"system_information",
		"station_information",
		"station_status",
		"free_bike_status",
		"system_hours",
		"system_calendar",
		"system_regions",
		"system_alerts",
	}

	// Picking random 10 systems
	rand.Seed(time.Now().Unix())
	rand.Shuffle(len(systems), func(i, j int) {
		systems[i], systems[j] = systems[j], systems[i]
	})
	systems = systems[0:10]

	renderTable(getSystemsData(client, systems))

	return nil
}

func getSystemsData(client *gbfs.Client, systems []gbfs.System) [][]string {
	data := [][]string{}

	for _, s := range systems {
		row, err := getSystemsRow(client, s)
		if err != nil {
			row = append(row, fmt.Sprintf("%v", err))
		}

		data = append(data, row)
	}

	return data
}

func getSystemsRow(client *gbfs.Client, s gbfs.System) ([]string, error) {
	log.Printf("Get GBFS: %s", s.AutoDiscoveryURL)
	gbfs, err := client.LoadGBFS(s.AutoDiscoveryURL)

	row := []string{
		s.Name,
		getSystemStatusEmoji(gbfs, err),
		getSystemLangs(gbfs, err),
	}

	return row, err
}

func getSystemStatusEmoji(gbfs *gbfs.GBFS, err error) string {
	if err != nil {
		return "🛑"
	}

	if gbfs.Data == nil {
		return "0️⃣"
	}

	return "✅"
}

func getSystemLangs(gbfs *gbfs.GBFS, err error) string {
	if err != nil {
		return "–"
	}

	langs := []string{}
	for lang := range gbfs.Data {
		langs = append(langs, lang)
	}

	return strings.Join(langs, ",")
}

func renderTable(data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "GBFS", "Lang"}) //, "Sy", "St", "SS", "FBS", "SH"})

	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}
