package scheduler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Timings struct {
	Date               string   `json:"date"`
	StartTime          string   `json:"start_time"`
	EndTime            string   `json:"end_time"`
	AvailableTimeSlots []string `json:"available_time_slots"`
}

type Court struct {
	Id              int         `json:"id"`
	CourtName       string      `json:"court_name"`
	MaximumCapacity interface{} `json:"maximum_capacity"`
	Timings         []Timings   `json:"timings"`
}

type T struct {
	Status string `json:"status"`
	Data   struct {
		Bookings []struct {
			Id       int `json:"id"`
			Location struct {
				Id           int     `json:"id"`
				LocationName string  `json:"location_name"`
				Courts       []Court `json:"courts"`
			} `json:"location"`
		} `json:"bookings"`
	} `json:"data"`
}

type Scheduler struct {
	Retries          time.Duration
	BaseUrl          string
	CourAvailableMap map[string]Court
	ScheduleStorage  *Storage
}

func NewScheduler(retries time.Duration, baseUrl string) *Scheduler {
	storage := NewStorage()

	return &Scheduler{
		Retries:          retries,
		BaseUrl:          baseUrl,
		CourAvailableMap: make(map[string]Court),
		ScheduleStorage:  storage,
	}
}

func (s *Scheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(s.Retries)
	defer ticker.Stop()
	s.SyncCourts()

mainLoop:
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			break mainLoop
		case <-ticker.C:
			s.ScheduleStorage.Clear()
			s.SyncCourts()
		}
	}
}

func (s *Scheduler) SyncCourts() {
	ws := sync.WaitGroup{}
	generateDates := GenerateDates()

	for _, date := range generateDates {
		ws.Add(1)

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)

			defer cancel()
			defer ws.Done()

			url := strings.Replace(s.BaseUrl, "{date}", date, -1)
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

			if err != nil {
				slog.Error("error in request ", "err", err)
				return
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				slog.Error("error in response ", "err", err)
				return
			}

			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					slog.Error("error closing body", "error", err)
				}
			}(resp.Body)

			body, err := io.ReadAll(resp.Body)

			if err != nil {
				slog.Error("error in response ", "error", err)
			}

			response := T{}
			err = json.NewDecoder(bytes.NewReader(body)).Decode(&response)

			for _, book := range response.Data.Bookings {
				if book.Id != 494 {
					continue
				}

				for _, court := range book.Location.Courts {
					s.ScheduleStorage.AddCourt(date, court)
				}
			}
		}()
	}

	ws.Wait()
}

func GenerateDates() []string {
	var dates []string

	for i := 0; i < 8; i++ {
		todayDate := time.Now()

		if i == 1 {
			dates = append(dates, todayDate.AddDate(0, 0, 1).Format("2006-01-02"))
		} else {
			dates = append(dates, todayDate.AddDate(0, 0, i).Format("2006-01-02"))
		}
	}

	return dates
}
