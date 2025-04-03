package main

import (
	"encoding/json"
	"io"
)

type siteDataResponse struct {
	ServiceTypes []struct {
		CategoryDescription string
		ServiceTypes        []struct {
			TypeId    int
			SubTypeId int
			Name      string
		}
	}
}

type LocationsByCountyResponse = []struct {
	County    string
	Locations []struct {
		Id   int
		Name string
		City string
	}
}

type availableLocationDatesResponse = struct {
	LocationAvailabilityDates []struct {
		LocationId int
		AvailableTimeSlots []struct {
			StartDateTime string
			Duration      int
			SlotId        int
		}
	}
	FirstAvailableDate string
}

func loadResponse[T any](reader io.Reader) (*T, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var result T
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
