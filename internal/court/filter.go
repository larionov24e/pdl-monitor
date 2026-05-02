package court

import (
	"slices"
	"time"
	"vdng-bot/internal/scheduler"
)

type AvailableCourt struct {
	Name string   `json:"name,omitempty"`
	Time []string `json:"time,omitempty"`
}

func FilterByTimeRange(timeFrom time.Time, timeTo time.Time, courts []scheduler.Court) []AvailableCourt {
	var availableCourts []AvailableCourt
	var times []string
	for t := timeFrom; t.Before(timeTo) || t.Equal(timeTo); t = t.Add(time.Hour) {
		times = append(times, t.Format("15:04:05"))
	}

	for _, court := range courts {
		var avaTimes []string

		if len(court.Timings) == 0 {
			continue
		}

		for k := 0; k < len(times); k++ {
			if slices.Contains(court.Timings[0].AvailableTimeSlots, times[k]) {
				avaTimes = append(avaTimes, times[k])
			}
		}

		if len(avaTimes) > 0 {
			availableCourts = append(availableCourts, AvailableCourt{
				Name: court.CourtName,
				Time: avaTimes,
			})
		}
	}

	return availableCourts
}
