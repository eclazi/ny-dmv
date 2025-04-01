package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/urfave/cli/v3"
)

func getAvailableLocationDatesFromNow(locationId int, typeId int) (*AvailableLocationDates, error) {
	result, err := getAvailableLocationDates(locationId, typeId, time.Now().Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	return result, nil
}

func getNextAvailableAppointment(locationId int, typeId int) (*time.Time, error) {
	locs, err := getAvailableLocationDatesFromNow(locationId, typeId)
	if err != nil {
		return nil, err
	}

	loc, _ := time.LoadLocation("America/New_York")
	result, err := time.ParseInLocation("2006-01-02T15:04:05", locs.FirstAvailableDate, loc)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func printServices() {
	siteData, err := getSiteData()
	if err != nil {
		panic(err)
	}

	for _, serviceType := range siteData.ServiceTypes {
		fmt.Println(serviceType.CategoryDescription)
		for _, service := range serviceType.ServiceTypes {
			fmt.Println("    ", service.TypeId, service.Name)
		}
	}
}

func collectLocations(serviceType int) map[int]string {
	locs, err := getLocationsByCounty(serviceType)
	if err != nil {
		panic(err)
	}

	result := make(map[int]string)
	for _, county := range *locs {
		for _, loc := range county.Locations {
			result[loc.Id] = loc.Name
		}
	}

	return result
}

func printLocations(serviceType int) {
	loc, err := getLocationsByCounty(serviceType)
	if err != nil {
		panic(err)
	}

	for _, county := range *loc {
		fmt.Println(county.County)
		for _, location := range county.Locations {
			fmt.Println("    ", location.Id, location.Name, location.City)
		}
	}
}

func earliestAppointment(service int, locations []int) (*time.Time, int) {
	var earliestTime time.Time
	var earliestLoc int

	for _, loc := range locations {
		apt, err := getNextAvailableAppointment(loc, service)
		if err != nil {
			panic(err)
		}

		if earliestTime.IsZero() || apt.Before(earliestTime) {
			earliestTime = *apt
			earliestLoc = loc
		}
	}

	return &earliestTime, earliestLoc
}

func getLocationName(locationId int, serviceId int) string {
	nameMap := collectLocations(serviceId)
	return nameMap[locationId]
}

func main() {

	cmd := &cli.Command{
		Commands: []*cli.Command{
			{
				Name: "services",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					printServices()
					return nil
				},
			},
			{
				Name: "locations",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					v, err := strconv.Atoi(cmd.Args().First())
					if err != nil {
						panic(err)
					}

					printLocations(v)
					return nil
				},
			},
			{
				Name: "appointment",
				Action: func(ctx context.Context, cmd *cli.Command) error {

					service, err := strconv.Atoi(cmd.Args().First())
					if err != nil {
						panic(err)
					}

					var locations []int
					for _, arg := range cmd.Args().Slice()[1:] {
						if v, err := strconv.Atoi(arg); err == nil {
							locations = append(locations, v)
						}
					}

					earliest, earliestLoc := earliestAppointment(service, locations)

					locName := getLocationName(earliestLoc, service)
					fmt.Println(locName, " - ", earliest.Format(time.RFC1123))

					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		panic(err)
	}
}
