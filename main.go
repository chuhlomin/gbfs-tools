package main

import (
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
			log.Printf("ERROR For system %q: %v", s.AutoDiscoveryURL, err)
		}

		data = append(data, row)
	}

	return data
}

func getSystemsRow(client *gbfs.Client, s gbfs.System) ([]string, error) {
	log.Printf("GET %s", s.AutoDiscoveryURL)
	gbfs, err := client.LoadGBFS(s.AutoDiscoveryURL)

	row := []string{
		s.Name,
		getSystemStatusEmoji(gbfs, err),
		getSystemLangs(gbfs, err),
		hasFeed(client, gbfs, err, "gbfs_versions"),
		hasFeed(client, gbfs, err, "system_information"),
		hasFeed(client, gbfs, err, "vehicle_types"),
		hasFeed(client, gbfs, err, "station_information"),
		hasFeed(client, gbfs, err, "station_status"),
		hasFeed(client, gbfs, err, "free_bike_status"),
		hasFeed(client, gbfs, err, "system_hours"),
		hasFeed(client, gbfs, err, "system_calendar"),
		hasFeed(client, gbfs, err, "system_regions"),
		hasFeed(client, gbfs, err, "system_pricing_plans"),
		hasFeed(client, gbfs, err, "system_alerts"),
		hasFeed(client, gbfs, err, "geofencing_zones"),
	}

	return row, err
}

func renderTable(data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "S", "Lang", "Vers", "SyI", "VT", "StI", "StSt", "FBS", "Hr", "Cal", "Reg", "Prc", "Alr", "Geo"})
	table.SetBorder(false)

	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}

func getSystemStatusEmoji(gbfs *gbfs.GBFS, err error) string {
	if err != nil {
		return "×"
	}

	if gbfs.Data == nil {
		return "!"
	}

	return "✔️"
}

func hasFeed(client *gbfs.Client, gbfs *gbfs.GBFS, err error, feedName string) string {
	if err != nil {
		return " "
	}

	for _, data := range gbfs.Data {
		for _, feed := range data.Feeds {
			if feed.Name == feedName {
				switch feedName {
				case "system_information":
					return "*"

				case "station_information":
					return "*"

				case "station_status":
					return "*"

				case "free_bike_status":
					return "*"
				// 	response, err := client.LoadFreeBikeStatus(feed.URL) // todo: log
				// 	if err != nil {
				// 		return "❗"
				// 	}
				// 	if len(response.Data.Bikes) == 0 {
				// 		return "0️⃣ "
				// 	}
				case "system_hours":
					// log.Printf("system_hours %s", feed.URL)
					return "*"
				// 	response, err := client.LoadSystemHours(feed.URL) // todo: log
				// 	if err != nil {
				// 		return "❗"
				// 	}
				// 	if len(response.Data.RentalHours) == 0 {
				// 		return "✅0️"
				// 	}
				case "system_calendar":
					// log.Printf("system_calendar %s", feed.URL)
					return "*"

				case "system_regions":
					// log.Printf("system_regions %s", feed.URL)
					return "*"

				case "system_pricing_plans":
					// log.Printf("system_pricing_plans %s", feed.URL)
					return "*"

				case "system_alerts":
					// log.Printf("system_alerts %s", feed.URL)
					return "*"

				case "geofencing_zones":
					// log.Printf("geofencing_zones %s", feed.URL)
					return "*"

				case "gbfs_versions":
					// log.Printf("gbfs_versions %s", feed.URL)
					return "*"

				case "vehicle_types":
					// log.Printf("vehicle_types %s", feed.URL)
					return "*"

				default:
					return "?"
				}
			}
		}
	}

	return " "
}

func getSystemLangs(gbfs *gbfs.GBFS, err error) string {
	if err != nil {
		return " "
	}

	langs := []string{}
	for lang := range gbfs.Data {
		langs = append(langs, lang)
	}

	return strings.Join(langs, ",")
}
