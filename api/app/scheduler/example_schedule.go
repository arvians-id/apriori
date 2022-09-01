package scheduler

import "log"

type ExampleSchedule struct {
}

func (scheduler *ExampleSchedule) Run() {
	log.Println("this is example")
}
