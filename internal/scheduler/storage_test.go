package scheduler

import (
	"testing"
)

func TestStorage_AddCache(t *testing.T) {
	storage := NewStorage()
	var sessionId int64

	sessionId = 1
	storage.AddCache(sessionId, DateKey, "value")

	checkVal := storage.GetSessionCacheByKey(sessionId, DateKey)
	if checkVal != "value" {
		t.Error("expected value but got ", checkVal)
	}
}

func TestStorage_GetCourtsByDate(t *testing.T) {
	storage := NewStorage()

	storage.scheduleCourts = map[string][]Court{
		"2026-05-03": {
			{
				Id:        1,
				CourtName: "Court A",
				Timings: []Timings{
					{
						Date:               "2026-05-03",
						StartTime:          "09:00:00",
						EndTime:            "12:00:00",
						AvailableTimeSlots: []string{"09:00:00", "10:00:00", "11:00:00"},
					},
				},
			},
			{
				Id:        2,
				CourtName: "Court B",
				Timings: []Timings{
					{
						Date:               "2026-05-03",
						StartTime:          "16:00:00",
						EndTime:            "20:00:00",
						AvailableTimeSlots: []string{"16:00:00", "17:00:00", "18:00:00"},
					},
				},
			},
		},
		"2026-05-04": {
			{
				Id:        3,
				CourtName: "Court C",
				Timings: []Timings{
					{
						Date:               "2026-05-04",
						StartTime:          "10:00:00",
						EndTime:            "14:00:00",
						AvailableTimeSlots: []string{"10:00:00", "11:00:00"},
					},
				},
			},
		},
	}

	courtsForDate1 := storage.GetCourtsByDate("2026-05-03")

	if len(courtsForDate1) != 2 {
		t.Error("expected 2 courts but got ", len(courtsForDate1))
	}

	courtsForDate2 := storage.GetCourtsByDate("2026-05-04")

	if len(courtsForDate2) != 1 {
		t.Error("expected 1 courts but got ", len(courtsForDate2))
	}
}
