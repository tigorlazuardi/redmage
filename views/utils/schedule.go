package utils

import (
	"time"

	"github.com/robfig/cron/v3"
)

var cronParser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)

func NextScheduleTime(in string) (out time.Time) {
	s, err := cronParser.Parse(in)
	if err != nil {
		return time.Time{}
	}
	return s.Next(time.Now())
}
