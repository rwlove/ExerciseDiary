package web

import (
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"

	"github.com/rwlove/WorkoutDiary/internal/models"
)

func indexHandler(c *gin.Context) {
	var guiData models.GuiData

	exs, err := dataStore.SelectEx()
	if err != nil {
		log.Println("ERROR indexHandler SelectEx:", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	sets, err := dataStore.SelectSet()
	if err != nil {
		log.Println("ERROR indexHandler SelectSet:", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	weights, err := dataStore.SelectW()
	if err != nil {
		log.Println("ERROR indexHandler SelectW:", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	exData.Exs = exs
	exData.Sets = sets
	exData.Weight = weights

	guiData.Config = appConfig
	guiData.ExData = exData
	guiData.GroupMap = createGroupMap()

	heatmaps := generateHeatMap(exData.Sets)
	guiData.IntensityHeatMap = heatmaps.IntensityMap
	guiData.ColorHeatMap = heatmaps.ColorMap

	sort.Slice(guiData.ExData.Exs, func(i, j int) bool {
		return guiData.ExData.Exs[i].Place < guiData.ExData.Exs[j].Place
	})
	sort.Slice(guiData.ExData.Weight, func(i, j int) bool {
		return guiData.ExData.Weight[i].Date < guiData.ExData.Weight[j].Date
	})

	c.HTML(http.StatusOK, "header.html", guiData)
	c.HTML(http.StatusOK, "index.html", guiData)
}

func createGroupMap() map[string]string {
	i := 0
	grMap := make(map[string]string)
	for _, ex := range exData.Exs {
		if _, ok := grMap[ex.Group]; !ok {
			grMap[ex.Group] = "grID" + fmt.Sprintf("%d", i)
			i++
		}
	}
	return grMap
}
