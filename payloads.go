package main

import (
	"encoding/json"
	"io"
)

type booking struct {
	ServiceTypeID   int       `json:"serviceTypeId"`
	ServiceTypeID2  *int      `json:"serviceTypeId2,omitempty"`
	BookingDateTime string    `json:"bookingDateTime"`
	BookingDuration int       `json:"bookingDuration"`
	FirstName       string    `json:"firstName"`
	LastName        string    `json:"lastName"`
	Email           string    `json:"email"`
	CellPhone       string    `json:"cellPhone"`
	SendSms         bool      `json:"sendSms"`
	SiteID          int       `json:"siteId"`
	SlotID          int       `json:"slotId"`
	DateOfBirth     *string   `json:"dateOfBirth,omitempty"`
}

func writePayload[T any](payload T, writer io.Writer) (int, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return 0, err
	}

	n, err := writer.Write(jsonData)
	if err != nil {
		return 0, err
	}

	return n, nil
}
