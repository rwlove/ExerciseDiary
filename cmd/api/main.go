package main

import (
	_ "time/tzdata"

	"github.com/rwlove/WorkoutDiary/internal/api"
)

func main() {
	api.Start()
}
