package scheduler

import (
	"github.com/roylee0704/gron"
	"github.com/roylee0704/gron/xtime"
)

var (
	everySecond       = gron.Every(1 * xtime.Second)
	everyThirtySecond = gron.Every(30 * xtime.Second)
	everyMinute       = gron.Every(1 * xtime.Minute)
	everyFiveMinute   = gron.Every(5 * xtime.Minute)
	everyThirtyMinute = gron.Every(30 * xtime.Minute)
	everyOneHour      = gron.Every(1 * xtime.Hour)
	everySixHour      = gron.Every(6 * xtime.Hour)
	daily             = gron.Every(1 * xtime.Day)
	weekly            = gron.Every(1 * xtime.Week)
	monthly           = gron.Every(30 * xtime.Day)
	yearly            = gron.Every(365 * xtime.Day)
	everyMidnight     = daily.At("24:00")
	everyMorning      = daily.At("05:00")
	dailyCustom       = daily.At("11:20")
)

func Execute() {
	schedule := gron.New()

	var clearTablePasswordResets ClearTablePasswordResetsSchedule
	var example ExampleSchedule

	schedule.Add(daily, &clearTablePasswordResets)
	schedule.Add(daily, &example)

	schedule.Start()
}
