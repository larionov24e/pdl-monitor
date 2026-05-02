package court

import (
	"reflect"
	"testing"
	"time"
	"vdng-bot/internal/scheduler"
)

func TestFilterByTimeRange(t *testing.T) {
	courts := []scheduler.Court{
		{
			Id:        1,
			CourtName: "test",
			Timings: []scheduler.Timings{
				{
					Date:               "2026-04-27",
					StartTime:          "09:00:00",
					EndTime:            "13:00:00",
					AvailableTimeSlots: []string{"09:00:00", "10:00:00", "11:00:00", "12:00:00"},
				},
			},
		},
		{
			Id:        2,
			CourtName: "test1",
			Timings: []scheduler.Timings{
				{
					Date:               "2026-04-27",
					StartTime:          "16:00:00",
					EndTime:            "20:00:00",
					AvailableTimeSlots: []string{"16:00:00", "17:00:00", "18:00:00", "19:00:00"},
				},
			},
		},
		{
			Id:        3,
			CourtName: "test2",
			Timings: []scheduler.Timings{
				{
					Date:               "2026-04-27",
					StartTime:          "16:00:00",
					EndTime:            "22:00:00",
					AvailableTimeSlots: []string{"16:00:00", "17:00:00", "18:00:00", "20:00:00", "21:00:00"},
				},
			},
		},
	}

	t.Run("Compare 1", func(t *testing.T) {
		dateFrom, err := time.Parse("15:04:05", "16:00:00")

		if err != nil {
			t.Fatal(err)
		}

		dateTo, err := time.Parse("15:04:05", "17:00:00")

		if err != nil {
			t.Fatal(err)
		}

		expected := []AvailableCourt{
			{
				Name: "test1",
				Time: []string{"16:00:00", "17:00:00"},
			},
			{
				Name: "test2",
				Time: []string{"16:00:00", "17:00:00"},
			},
		}

		compareResult := FilterByTimeRange(dateFrom, dateTo, courts)

		if !reflect.DeepEqual(expected, compareResult) {
			t.Errorf("Step 1: FilterByTimeRange: expected %v, got %v", expected, compareResult)
		}
	})

	t.Run("Compare 2", func(t *testing.T) {
		dateFrom, err := time.Parse("15:04:05", "09:00:00")

		if err != nil {
			t.Fatal(err)
		}

		dateTo, err := time.Parse("15:04:05", "12:00:00")

		if err != nil {
			t.Fatal(err)
		}

		expected := []AvailableCourt{
			{
				Name: "test",
				Time: []string{"09:00:00", "10:00:00", "11:00:00", "12:00:00"},
			},
		}

		compareResult := FilterByTimeRange(dateFrom, dateTo, courts)

		if !reflect.DeepEqual(expected, compareResult) {
			t.Errorf("Step 2: FilterByTimeRange: expected %v, got %v", expected, compareResult)
		}
	})

}
