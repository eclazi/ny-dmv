package main

import (
	"dmv-ny/pkg/dmvapi"
	"testing"
	"time"
)

func TestFilterAppointments(t *testing.T) {
	now := time.Now()
	withinDuration := 2 * time.Hour

	appointments := []dmvapi.Appointment{
		{DateTime: now.Add(1 * time.Hour)},
		{DateTime: now.Add(3 * time.Hour)},
		{DateTime: now.Add(-1 * time.Hour)},
	}

	filtered := filterAppointments(appointments, now, withinDuration)

	if len(filtered) != 1 {
		t.Errorf("Expected 1 appointment, got %d", len(filtered))
	}

	if !filtered[0].DateTime.Equal(appointments[0].DateTime) {
		t.Errorf("Expected appointment at %v, got %v", appointments[0].DateTime, filtered[0].DateTime)
	}
}