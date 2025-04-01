package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"
)

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

type Location struct {
	Id   int
	Name string
	City string
}

type LocationsByCounty = []struct {
	County    string
	Locations []Location
}

func getLocationsByCounty(serviceTypeId int) (*LocationsByCounty, error) {
	startDate := time.Now().Format(time.RFC3339)
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
