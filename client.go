package main

import (
	"fmt"
	"net/http"
	"time"
)

const (
	baseURL                = "https://publicwebsiteapi.nydmvreservation.com/api/"
)

type Service struct {
	Id   int
	Name string
}

type Location struct {
	Id   int
	Name string
	City string
}

type Appointment struct {
	LocationId int
	DateTime   time.Time
	SlotId     int
	Duration   int
	ServiceId  int
}

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{},
	}
}

func (c *Client) GetServices() ([]Service, error) {
	resp, err := http.Get(baseURL + "SiteData")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	sdr, err := loadResponse[siteDataResponse](resp.Body)
	if err != nil {
		return nil, err
	}

	var services []Service
	for _, serviceType := range sdr.ServiceTypes {
		for _, service := range serviceType.ServiceTypes {
			services = append(services, Service{
				Id:   service.TypeId,
				Name: service.Name,
			})
		}
	}

	return services, nil
}

func (c *Client) GetLocations(service *Service) ([]Location, error) {
	resp, err := http.Get(baseURL + "LocationsByCounty" + "?serviceTypeId=" + fmt.Sprint(service.Id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	locs, err := loadResponse[LocationsByCountyResponse](resp.Body)
	if err != nil {
		return nil, err
	}
	var locations []Location
	for _, county := range *locs {
		for _, loc := range county.Locations {
			locations = append(locations, Location{
				Id:   loc.Id,
				Name: loc.Name,
				City: loc.City,
			})
		}
	}
	return locations, nil
}

func (c *Client) GetAppointments(location *Location, service *Service) ([]Appointment, error) {
	resp, err := http.Get(baseURL + "AvailableLocationDates" + "?locationId=" + fmt.Sprint(location.Id) + "&typeId=" + fmt.Sprint(service.Id) + "&startDate=" + time.Now().Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	aldr, err := loadResponse[availableLocationDatesResponse](resp.Body)
	if err != nil {
		return nil, err
	}

	loc, _ := time.LoadLocation("America/New_York")

	var appointments []Appointment
	for _, locationAvailabilityDates := range aldr.LocationAvailabilityDates {
		for _, availableTimeSlot := range locationAvailabilityDates.AvailableTimeSlots {
			tim, _ := time.ParseInLocation("2006-01-02T15:04:05", availableTimeSlot.StartDateTime, loc)
			appointments = append(appointments, Appointment{
				LocationId: locationAvailabilityDates.LocationId,
				DateTime:   tim,
				SlotId:     availableTimeSlot.SlotId,
				Duration:   availableTimeSlot.Duration,
				ServiceId:  service.Id,
			})
		}
	}

	return appointments, nil
}
