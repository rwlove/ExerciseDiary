package web

import (
	"log"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"

	"github.com/aceberg/ExerciseDiary/internal/models"
)

func statsHandler(c *gin.Context) {
	var guiData models.GuiData

	exs, err := dataStore.SelectEx()
	if err != nil {
		log.Println("ERROR statsHandler SelectEx:", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	sets, err := dataStore.SelectSet()
	if err != nil {
		log.Println("ERROR statsHandler SelectSet:", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	guiData.ExData.Exs = exs
	guiData.ExData.Sets = sets
	guiData.Config = appConfig

	guiData.GroupMap = make(map[string]string)
	for _, ex := range guiData.ExData.Sets {
		if _, ok := guiData.GroupMap[ex.Name]; !ok {
			guiData.GroupMap[ex.Name] = ex.Name
		}
	}

	sort.Slice(guiData.ExData.Sets, func(i, j int) bool {
		return guiData.ExData.Sets[i].Date < guiData.ExData.Sets[j].Date
	})

	c.HTML(http.StatusOK, "header.html", guiData)
	c.HTML(http.StatusOK, "stats.html", guiData)
}
