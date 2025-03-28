package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

//type SiteData = map[string]interface{}

type SiteData struct {
	ServiceTypes []struct {
		CategoryDescription string
		ServiceTypes        []struct {
			TypeId    int
			SubTypeId int
			Name      string
		}
	}
}

func getSiteData() (*SiteData, error) {
	resp, err := http.Get("https://publicwebsiteapi.nydmvreservation.com/api/SiteData")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil // TODO, useful error
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result SiteData
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

//https://publicwebsiteapi.nydmvreservation.com/api/LocationsByCounty?serviceTypeId=203&startDate=2025-03-27T00:41:58.040Z

type LocationsByCounty = []struct {
	County    string
	Locations []struct {
		Id   int
		Name string
		City string
	}
}

func getLocationsByCounty(serviceTypeId int, startDate string) (*LocationsByCounty, error) {
	resp, err := http.Get("https://publicwebsiteapi.nydmvreservation.com/api/LocationsByCounty?serviceTypeId=" + strconv.Itoa(serviceTypeId) + "&startDate=" + startDate)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil // TODO, useful error
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result LocationsByCounty
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

//https://publicwebsiteapi.nydmvreservation.com/api/AvailableLocationDates?locationId=34&typeId=203&startDate=2025-03-27T00:44:00.269Z

type AvailableLocationDates = struct {
	FirstAvailableDate string
}

func getAvailableLocationDates(locationId int, typeId int, startDate string) (*AvailableLocationDates, error) {
	resp, err := http.Get("https://publicwebsiteapi.nydmvreservation.com/api/AvailableLocationDates?locationId=" + strconv.Itoa(locationId) + "&typeId=" + strconv.Itoa(typeId) + "&startDate=" + startDate)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil // TODO, useful error
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result AvailableLocationDates
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

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

func main() {
	serviceTypeId := flag.Int("serviceTypeId", 0, "Service Type ID")
	locationId := flag.Int("locationId", 0, "Location ID")
	flag.Parse()

	if *serviceTypeId == 0 {
		// List available services
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
	} else if *locationId == 0 {
		// List available locations
		startDate := time.Now().Format(time.RFC3339)
		locations, err := getLocationsByCounty(*serviceTypeId, startDate)
		if err != nil {
			panic(err)
		}

		for _, county := range *locations {
			fmt.Println(county.County)
			for _, location := range county.Locations {
				fmt.Println("    ", location.Id, location.Name, location.City)
			}
		}
	} else {
		// List available dates
		t, err := getNextAvailableAppointment(*locationId, *serviceTypeId)
		if err != nil {
			panic(err)
		}

		fmt.Println("First available date:", t.Format(time.RFC1123))
	}

}
