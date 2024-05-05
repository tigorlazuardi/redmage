package schedulehistoriesview

import (
	"time"

	"github.com/tigorlazuardi/redmage/models"
)

type Data struct {
	Error     string
	Schedules models.ScheduleHistorySlice
	Total     int64
	Timezone  *time.Location
}

func (data *Data) splitByDays() (out []*splitByDaySchedules) {
	for _, schedule := range data.Schedules {
		day := time.Unix(schedule.CreatedAt, 0).In(data.Timezone)
		date := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())

		var found bool
	inner:
		for _, split := range out {
			if split.Date.Equal(date) {
				found = true
				split.Schedules = append(split.Schedules, schedule)
				break inner
			}
		}
		if !found {
			out = append(out, &splitByDaySchedules{
				Day:       day,
				Date:      date,
				Schedules: models.ScheduleHistorySlice{schedule},
			})
		}
	}

	return out
}

type splitByDaySchedules struct {
	Day       time.Time
	Date      time.Time
	Schedules models.ScheduleHistorySlice
}
