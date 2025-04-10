package main

import (
	"context"
	"dmv-ny/pkg/dmvapi"
	"fmt"
	"os"
	"sort"

	"github.com/urfave/cli/v3"
)

// handleError provides consistent error handling throughout the application
func handleError(err error, message string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", message, err)
		os.Exit(1)
	}
}

func printServices() {
	client := dmvapi.NewClient()
	services, err := client.GetServices()
	handleError(err, "Failed to get services")

	for _, service := range services {
		fmt.Printf("%d: %s\n", service.Id, service.Name)
	}
}

func printLocations(serviceId int64) {
	client := dmvapi.NewClient()
	locations, err := client.GetLocations(int(serviceId))
	handleError(err, "Failed to get locations")

	for _, location := range locations {
		fmt.Printf("%d: %s, %s\n", location.Id, location.Name, location.City)
	}
}

// getLocationNames retrieves and maps location IDs to their names
func getLocationNames(client *dmvapi.Client, serviceId int) (map[int]string, error) {
	locations, err := client.GetLocations(serviceId)
	if err != nil {
		return nil, err
	}

	locationNames := make(map[int]string, len(locations))
	for _, location := range locations {
		locationNames[location.Id] = location.Name
	}
	return locationNames, nil
}

func printAppointments(locationIds []int64, serviceId int64) {
	client := dmvapi.NewClient()

	locationNames, err := getLocationNames(client, int(serviceId))
	handleError(err, "Failed to get location names")

	var appointments []dmvapi.Appointment
	for _, locationId := range locationIds {
		appts, err := client.GetAppointments(int(locationId), int(serviceId))
		handleError(err, fmt.Sprintf("Failed to get appointments for location %d", locationId))
		appointments = append(appointments, appts...)
	}

	sort.Slice(appointments, func(i, j int) bool {
		return appointments[i].DateTime.Before(appointments[j].DateTime)
	})

	for _, appt := range appointments {
		fmt.Printf("%v - %s - Location ID: %d, Slot ID: %d, Duration: %d\n",
			appt.DateTime.Format("2006-01-02 15:04:05"),
			locationNames[appt.LocationId],
			appt.LocationId,
			appt.SlotId,
			appt.Duration)
	}
}

func bookAppointment(locationId, serviceId, slotId int64, firstName, lastName, email, phone string) {
	client := dmvapi.NewClient()
	appts, err := client.GetAppointments(int(locationId), int(serviceId))
	handleError(err, "Failed to get appointments")

	// Find the appointment with the given slot ID
	var appointment dmvapi.Appointment
	found := false

	for _, appt := range appts {
		if appt.SlotId == int(slotId) {
			appointment = appt
			found = true
			break
		}
	}

	if !found {
		fmt.Println("No appointment found with the given slot ID.")
		return
	}

	err = client.BookAppointment(appointment, firstName, lastName, email, phone)
	handleError(err, "Failed to book appointment")

	fmt.Println("Appointment booked successfully!")
}

func main() {
	var serviceId int64
	var locationIds []int64
	var locationId int64
	var slotId int64
	var firstName, lastName, email, phone string

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
				Arguments: []cli.Argument{
					&cli.IntArg{
						Name:        "serviceId",
						Destination: &serviceId,
						Min:         1,
						Max:         1,
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					printLocations(serviceId)
					return nil
				},
			},
			{
				Name: "appointments",
				Arguments: []cli.Argument{
					&cli.IntArg{
						Name:        "serviceId",
						Destination: &serviceId,
						Min:         1,
						Max:         1,
					},
					&cli.IntArg{
						Name:   "locationIds",
						Values: &locationIds,
						Min:    1,
						Max:    10,
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					printAppointments(locationIds, serviceId)
					return nil
				},
			},
			{
				Name: "book",
				Arguments: []cli.Argument{
					&cli.IntArg{
						Name:        "serviceId",
						Destination: &serviceId,
						Min:         1,
						Max:         1,
					},
					&cli.IntArg{
						Name:        "locationId",
						Destination: &locationId,
						Min:         1,
						Max:         1,
					},
					&cli.IntArg{
						Name:        "slotId",
						Destination: &slotId,
						Min:         1,
						Max:         1,
					},
					&cli.StringArg{
						Name:        "firstName",
						Destination: &firstName,
						Min:         1,
						Max:         1,
					},
					&cli.StringArg{
						Name:        "lastName",
						Destination: &lastName,
						Min:         1,
						Max:         1,
					},
					&cli.StringArg{
						Name:        "email",
						Destination: &email,
						Min:         1,
						Max:         1,
					},
					&cli.StringArg{
						Name:        "phone",
						Destination: &phone,
						Min:         1,
						Max:         1,
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					bookAppointment(locationId, serviceId, slotId, firstName, lastName, email, phone)
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		handleError(err, "Command failed")
	}
}
