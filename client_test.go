package main

import (
	"fmt"
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient()

	services, err := client.GetServices()
	if err != nil {
		t.Error(err)
	}

	if len(services) == 0 {
		t.Error("No services found")
	}

	fmt.Println(services)

	service := services[0]
	locations, err := client.GetLocations(service.Id)
	if err != nil {
		t.Error(err)
	}

	if len(locations) == 0 {
		t.Error("No locations found")
	}

	fmt.Println(locations)

	appointment, err := client.GetAppointments(locations[0].Id, service.Id)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(appointment)
}
