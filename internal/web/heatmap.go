package web

import (
	"strconv"
	"time"

	"github.com/aceberg/ExerciseDiary/internal/models"
)

// HeatMapResult contains both intensity and color heatmaps
type HeatMapResult struct {
	IntensityMap []models.HeatMapData
	ColorMap     []models.HeatMapData
}

// WorkoutData stores colors and names for workouts on a specific date
type WorkoutData struct {
	Colors      []string
	Names       []string
	Intensities []int
	Weights     []string
	Reps        []int
}

// generateHeatMap builds both heatmaps from the provided sets slice.
// Accepting sets as a parameter (instead of reading exData directly) keeps this
// function pure and compatible with both monolith and split-frontend modes.
func generateHeatMap(sets []models.Set) HeatMapResult {
	var intensityMap []models.HeatMapData
	var colorMap []models.HeatMapData
	var heat models.HeatMapData

	w := 52

	max := time.Now()
	min := max.AddDate(0, 0, -7*w)

	startDate := weekStartDate(min)
	countMap := countHeat(sets)
	workoutData := getWorkoutData(sets)

	dow := []string{"Mo", "Tu", "We", "Th", "Fr", "Sa", "Su"}

	for _, day := range dow {
		heat.Y = day

		for i := 0; i < w+1; i++ {
			heat.X = strconv.Itoa(i)
			heat.D = startDate.AddDate(0, 0, 7*i).Format("2006-01-02")

			heat.V = countMap[heat.D]
			intensityMap = append(intensityMap, heat)

			if data, exists := workoutData[heat.D]; exists {
				heat.V = len(data.Colors)
				heat.Colors = data.Colors
				heat.WorkoutNames = data.Names
				heat.WorkoutIntensities = data.Intensities
				heat.WorkoutWeights = data.Weights
				heat.WorkoutReps = data.Reps
				colorMap = append(colorMap, heat)
			} else {
				heat.V = 0
				heat.Colors = []string{}
				heat.WorkoutNames = []string{}
				heat.WorkoutIntensities = []int{}
				heat.WorkoutWeights = []string{}
				heat.WorkoutReps = []int{}
				colorMap = append(colorMap, heat)
			}
		}

		startDate = startDate.AddDate(0, 0, 1)
	}

	return HeatMapResult{
		IntensityMap: intensityMap,
		ColorMap:     colorMap,
	}
}

func weekStartDate(date time.Time) time.Time {
	offset := (int(time.Monday) - int(date.Weekday()) - 7) % 7
	result := date.Add(time.Duration(offset*24) * time.Hour)
	return result
}

func countHeat(sets []models.Set) map[string]int {
	countMap := make(map[string]int)
	for _, ex := range sets {
		countMap[ex.Date] += ex.Intensity + 1
	}
	return countMap
}

func getWorkoutData(sets []models.Set) map[string]WorkoutData {
	workoutMap := make(map[string]WorkoutData)
	for _, set := range sets {
		date := set.Date
		if set.WorkoutColor == "" {
			continue
		}
		if _, exists := workoutMap[date]; !exists {
			workoutMap[date] = WorkoutData{
				Colors:      []string{},
				Names:       []string{},
				Intensities: []int{},
				Weights:     []string{},
				Reps:        []int{},
			}
		}
		data := workoutMap[date]
		data.Colors = append(data.Colors, set.WorkoutColor)
		data.Names = append(data.Names, set.Name)
		data.Intensities = append(data.Intensities, set.Intensity)
		data.Weights = append(data.Weights, set.Weight.String())
		data.Reps = append(data.Reps, set.Reps)
		workoutMap[date] = data
	}
	return workoutMap
}
