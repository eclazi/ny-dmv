package main

// import (
// 	"encoding/json"
// 	"io"
// 	"net/http"
// 	"strconv"
// 	"time"
// )

// type LocationId int
// type ServiceType = int

// type SiteDataResponse struct {
// 	ServiceTypes []struct {
// 		CategoryDescription string
// 		ServiceTypes        []struct {
// 			TypeId    int
// 			SubTypeId int
// 			Name      string
// 		}
// 	}
// }

// func getSiteData() (*SiteDataResponse, error) {
// 	resp, err := http.Get("https://publicwebsiteapi.nydmvreservation.com/api/SiteData")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return nil, nil // TODO, useful error
// 	}

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var result SiteDataResponse
// 	err = json.Unmarshal(body, &result)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &result, nil
// }

// //https://publicwebsiteapi.nydmvreservation.com/api/LocationsByCounty?serviceTypeId=203&startDate=2025-03-27T00:41:58.040Z

// type Location struct {
// 	Id   LocationId
// 	Name string
// 	City string
// }

// type LocationsByCountyResponse = []struct {
// 	County    string
// 	Locations []Location
// }

// func getLocationsByCounty(serviceTypeId ServiceType) (*LocationsByCountyResponse, error) {
// 	startDate := time.Now().Format(time.RFC3339)
// 	resp, err := http.Get("https://publicwebsiteapi.nydmvreservation.com/api/LocationsByCounty?serviceTypeId=" + strconv.Itoa(serviceTypeId) + "&startDate=" + startDate)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return nil, nil // TODO, useful error
// 	}

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var result LocationsByCountyResponse
// 	err = json.Unmarshal(body, &result)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &result, nil
// }

// //https://publicwebsiteapi.nydmvreservation.com/api/AvailableLocationDates?locationId=34&typeId=203&startDate=2025-03-27T00:44:00.269Z

// type AvailableLocationDatesResponse = struct {
// 	FirstAvailableDate string
// }

// func getAvailableLocationDates(locationId LocationId, typeId ServiceType, startDate string) (*AvailableLocationDatesResponse, error) {
// 	resp, err := http.Get("https://publicwebsiteapi.nydmvreservation.com/api/AvailableLocationDates?locationId=" + strconv.Itoa(int(locationId)) + "&typeId=" + strconv.Itoa(typeId) + "&startDate=" + startDate)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return nil, nil // TODO, useful error
// 	}

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var result AvailableLocationDatesResponse
// 	err = json.Unmarshal(body, &result)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &result, nil
// }

// func getAvailableLocationDatesFromNow(locationId LocationId, typeId int) (*AvailableLocationDatesResponse, error) {
// 	result, err := getAvailableLocationDates(locationId, typeId, time.Now().Format(time.RFC3339))
// 	if err != nil {
// 		return nil, err
// 	}
// 	return result, nil
// }

// func getNextAvailableAppointment(locationId LocationId, typeId int) (*time.Time, error) {
// 	locs, err := getAvailableLocationDatesFromNow(locationId, typeId)
// 	if err != nil {
// 		return nil, err
// 	}

// 	loc, _ := time.LoadLocation("America/New_York")
// 	result, err := time.ParseInLocation("2006-01-02T15:04:05", locs.FirstAvailableDate, loc)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &result, nil
// }

// type AppointmentsList = []time.Time

// type AppointmentsMap = map[Location]AppointmentsList

// func getAppointements(service ServiceType, locations []Location) AppointmentsMap {
// 	results := make(AppointmentsMap)

// 	for _, loc := range locations {
// 		date, err := getNextAvailableAppointment(loc.Id, service)
// 		if err != nil {
// 			panic(err)
// 		}

// 		results[loc] = append(results[loc], *date)
// 	}

// 	return results
// }
