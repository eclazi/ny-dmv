package main

import (
	"context"
	"dmv-ny/pkg/dmvapi"
	"fmt"
	"os"
	"sort"

	"github.com/urfave/cli/v3"
)

func printServices() {
	client := dmvapi.NewClient()
	services, err := client.GetServices()
	if err != nil {
		panic(err)
	}

	for _, service := range services {
		fmt.Println(service.Id, ": ", service.Name)
	}
}
func printLocations(serviceId int64) {
	client := dmvapi.NewClient()
	locations, err := client.GetLocations(int(serviceId))
	if err != nil {
		panic(err)
	}

	for _, location := range locations {
		fmt.Println(location.Id, ": ", location.Name, ", ", location.City)
	}
}

func printAppointments(locationIds []int64, serviceId int64) {
	client := dmvapi.NewClient()

	locations, err := client.GetLocations(int(serviceId))
	if err != nil {
		panic(err)
	}

	locationNames := make(map[int]string, 0)
	for _, location := range locations {
		locationNames[location.Id] = location.Name
	}

	var appointments []dmvapi.Appointment
	for _, locationId := range locationIds {
		appts, err := client.GetAppointments(int(locationId), int(serviceId))
		if err != nil {
			panic(err)
		}
		appointments = append(appointments, appts...)
	}

	sort.Slice(appointments, func(i, j int) bool {
		return appointments[i].DateTime.Before(appointments[j].DateTime)
	})

	for _, appointment := range appointments {
		fmt.Printf("%v - %s - Location ID : %d,  Slot ID : %d, Duration : %d\n", appointment.DateTime.Format("2006-01-02 15:04:05"),
			locationNames[appointment.LocationId], appointment.LocationId, appointment.SlotId, appointment.Duration)
	}
}

func bookAppointment(locationId int64, serviceId int64, slotId int64, firstName string, lastName string, email string, phone string) {
	client := dmvapi.NewClient()
	apps, err := client.GetAppointments(int(locationId), int(serviceId))
	if err != nil {
		panic(err)
	}

	// Find the appointment with the given slot ID
	var appointment dmvapi.Appointment
	for _, app := range apps {
		if app.SlotId == int(slotId) {
			appointment = app
			break
		}
	}

	if appointment.SlotId == 0 {
		fmt.Println("No appointment found with the given slot ID.")
		return
	}

	if err := client.BookAppointment(appointment, firstName, lastName, email, phone); err != nil {
		panic(err)
	}

	fmt.Println("Appointment booked successfully!")
}
func main() {
	var serviceId int64
	var locationIds []int64

	var locationId int64
	var slotId int64
	var firstName string
	var lastName string
	var email string
	var phone string

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
		panic(err)
	}
}
