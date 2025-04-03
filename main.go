package main

// import (
// 	"context"
// 	"fmt"
// 	"os"
// 	"strconv"
// 	"time"

// 	"github.com/urfave/cli/v3"
// )

// func printServices() {
// 	siteData, err := getSiteData()
// 	if err != nil {
// 		panic(err)
// 	}

// 	for _, serviceType := range siteData.ServiceTypes {
// 		fmt.Println(serviceType.CategoryDescription)
// 		for _, service := range serviceType.ServiceTypes {
// 			fmt.Println("    ", service.TypeId, service.Name)
// 		}
// 	}
// }

// func collectLocations(serviceType int) map[LocationId]string {
// 	locs, err := getLocationsByCounty(serviceType)
// 	if err != nil {
// 		panic(err)
// 	}

// 	result := make(map[LocationId]string)
// 	for _, county := range *locs {
// 		for _, loc := range county.Locations {
// 			result[loc.Id] = loc.Name
// 		}
// 	}

// 	return result
// }

// func printLocations(serviceType int) {
// 	loc, err := getLocationsByCounty(serviceType)
// 	if err != nil {
// 		panic(err)
// 	}

// 	for _, county := range *loc {
// 		fmt.Println(county.County)
// 		for _, location := range county.Locations {
// 			fmt.Println("    ", location.Id, location.Name, location.City)
// 		}
// 	}
// }

// func earliestAppointment(service int, locations []LocationId) (*time.Time, LocationId) {
// 	var earliestTime time.Time
// 	var earliestLoc LocationId

// 	for _, loc := range locations {
// 		apt, err := getNextAvailableAppointment(loc, service)
// 		if err != nil {
// 			panic(err)
// 		}

// 		if earliestTime.IsZero() || apt.Before(earliestTime) {
// 			earliestTime = *apt
// 			earliestLoc = loc
// 		}
// 	}

// 	return &earliestTime, earliestLoc
// }

// func getLocationName(locationId LocationId, serviceId int) string {
// 	nameMap := collectLocations(serviceId)
// 	return nameMap[locationId]
// }

// func main() {

// 	cmd := &cli.Command{
// 		Commands: []*cli.Command{
// 			{
// 				Name: "services",
// 				Action: func(ctx context.Context, cmd *cli.Command) error {
// 					printServices()
// 					return nil
// 				},
// 			},
// 			{
// 				Name: "locations",
// 				Action: func(ctx context.Context, cmd *cli.Command) error {
// 					v, err := strconv.Atoi(cmd.Args().First())
// 					if err != nil {
// 						panic(err)
// 					}

// 					printLocations(v)
// 					return nil
// 				},
// 			},
// 			{
// 				Name: "appointmenat",
// 				Action: func(ctx context.Context, cmd *cli.Command) error {

// 					service, err := strconv.Atoi(cmd.Args().First())
// 					if err != nil {
// 						panic(err)
// 					}

// 					var locations []LocationId
// 					for _, arg := range cmd.Args().Slice()[1:] {
// 						if v, err := strconv.Atoi(arg); err == nil {
// 							locations = append(locations, LocationId(v))
// 						}
// 					}

// 					earliest, earliestLoc := earliestAppointment(service, locations)

// 					locName := getLocationName(earliestLoc, service)
// 					fmt.Println(locName, " - ", earliest.Format(time.RFC1123))

// 					return nil
// 				},
// 			},
// 		},
// 	}

// 	if err := cmd.Run(context.Background(), os.Args); err != nil {
// 		panic(err)
// 	}
// }
