# gbfs-go

## What is GBFS

[The General Bikeshare Feed Specification][], known as GBFS, is the open data standard for bikeshare. GBFS makes real-time data feeds in a uniform format publicly available online, with an emphasis on findability. GBFS is intended to make information publicly available online; therefore information that is personally identifiable is not currently and will not become part of the core specification.

## What is gbfs-go

gbfs-go is a [Go][] library that aimed to simplify reading and writing GBFS data, it provides Go structures and HTTP-client wrapper.

## Examples

* [systems](https://github.com/chuhlomin/gbfs-go/blob/main/examples/systems/main.go)
* [gbfs](https://github.com/chuhlomin/gbfs-go/blob/main/examples/gbfs/main.go)
* [gbfs_version](https://github.com/chuhlomin/gbfs-go/blob/main/examples/gbfs_version/main.go)
* [system_information](https://github.com/chuhlomin/gbfs-go/blob/main/examples/system_information/main.go)
* [vehicle_types](https://github.com/chuhlomin/gbfs-go/blob/main/examples/vehicle_types/main.go)
* [station_information](https://github.com/chuhlomin/gbfs-go/blob/main/examples/station_information/main.go)
* [station_status](https://github.com/chuhlomin/gbfs-go/blob/main/examples/station_status/main.go)
* [free_bike_status](https://github.com/chuhlomin/gbfs-go/blob/main/examples/free_bike_status/main.go)
* [system_hours](https://github.com/chuhlomin/gbfs-go/blob/main/examples/system_hours/main.go)
* [system_calendar](https://github.com/chuhlomin/gbfs-go/blob/main/examples/system_calendar/main.go)
* [system_regions](https://github.com/chuhlomin/gbfs-go/blob/main/examples/system_regions/main.go)
* [system_pricing_plans](https://github.com/chuhlomin/gbfs-go/blob/main/examples/system_pricing_plans/main.go)
* [system_alerts](https://github.com/chuhlomin/gbfs-go/blob/main/examples/system_alerts/main.go)

Try examples:

```bash
go run examples/system_information/main.go
```

## License

[MIT License][]

Keep in mind that data feeds may have a [different license] (like CC0-1.0, CC-BY-4.0, etc).

[The General Bikeshare Feed Specification]: https://github.com/NABSA/gbfs
[Go]: https://golang.org
[different license]: https://github.com/NABSA/gbfs/blob/master/data-licenses.md
[MIT License]: https://github.com/chuhlomin/gbfs-go/blob/main/LICENSE
