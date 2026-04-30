package api

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rwlove/WorkoutDiary/internal/auth"
	"github.com/rwlove/WorkoutDiary/internal/check"
	"github.com/rwlove/WorkoutDiary/internal/conf"
	"github.com/rwlove/WorkoutDiary/internal/db"
	"github.com/rwlove/WorkoutDiary/internal/migrate"
	"github.com/rwlove/WorkoutDiary/internal/models"
	"github.com/rwlove/WorkoutDiary/internal/store"
)

var (
	appConfig models.Conf
	authConf  auth.Conf
	dataStore store.Store
)

// Start initialises the database, reads config from environment variables,
// and begins serving the JSON API.
//
// dataDir is the only required parameter — it sets where the SQLite database
// lives. All other settings (port, API key, theme, auth, …) are read from
// environment variables:
//
//	PORT        listen port              (default: 8851)
//	API_KEY     required X-Api-Key value (default: "", no auth)
//	HOST        listen host              (default: 0.0.0.0)
//	THEME       UI theme                 (default: grass)
//	COLOR       light or dark            (default: light)
//	HEATCOLOR   heatmap colour           (default: #03a70c)
//	PAGESTEP    rows per page            (default: 10)
//	AUTH        enable session auth      (default: false)
//	AUTH_USER   username                 (default: "")
//	AUTH_PASSWORD bcrypt password        (default: "")
//	AUTH_EXPIRE session expiry           (default: 7d)
func Start(dataDir string) {
	appConfig, authConf = conf.GetFromEnv()
	appConfig.DirPath = dataDir
	appConfig.DBPath = dataDir + "/sqlite.db"
	// ConfPath left empty — Settings changes are in-memory only in env-var mode.
	check.Path(appConfig.DBPath)

	apiKey := os.Getenv("API_KEY")
	postgresDSN := os.Getenv("POSTGRES_DSN")

	if postgresDSN != "" {
		log.Println("INFO: starting API server with PostgreSQL backend")
		pgStore, err := store.NewPostgres(postgresDSN)
		if err != nil {
			log.Fatalf("ERROR: connect to postgres: %v", err)
		}
		if err := db.MigratePostgres(pgStore.Pool()); err != nil {
			log.Fatalf("ERROR: postgres schema migration: %v", err)
		}
		dataStore = pgStore
		appConfig.PostgresEnabled = true
	} else {
		log.Println("INFO: starting API server with SQLite backend, db =", appConfig.DBPath)
		db.Create(appConfig.DBPath)
		dataStore = store.NewSQLite(appConfig.DBPath)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// API key middleware (optional – skip when API_KEY is unset)
	if apiKey != "" {
		r.Use(apiKeyMiddleware(apiKey))
	}

	// Exercises
	r.GET("/api/exercises", getExercises)
	r.POST("/api/exercises", postExercise)
	r.PUT("/api/exercises/:id", putExercise)
	r.DELETE("/api/exercises/:id", deleteExercise)

	// Sets
	r.GET("/api/sets", getSets)
	r.PUT("/api/sets/date/:date", putSetsByDate)

	// Body weight
	r.GET("/api/weight", getWeight)
	r.POST("/api/weight", postWeight)
	r.DELETE("/api/weight/:id", deleteWeight)

	// Config
	r.GET("/api/config", getConfig)
	r.PUT("/api/config", putConfig)
	r.PUT("/api/config/auth", putConfigAuth)

	// Migration (only registered when PostgreSQL is active)
	if appConfig.PostgresEnabled {
		r.POST("/api/migrate/sqlite-to-postgres", migrateFromSQLite)
	}

	address := appConfig.Host + ":" + appConfig.Port
	log.Println("=================================== ")
	log.Printf("API server at http://%s", address)
	log.Println("=================================== ")

	check.IfError(r.Run(address))
}

// ─── middleware ───────────────────────────────────────────────────────────────

func apiKeyMiddleware(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-Api-Key") != key {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.Next()
	}
}

// ─── exercises ────────────────────────────────────────────────────────────────

func getExercises(c *gin.Context) {
	exs, err := dataStore.SelectEx()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, exs)
}

func postExercise(c *gin.Context) {
	var ex models.Exercise
	if err := c.ShouldBindJSON(&ex); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := dataStore.InsertEx(ex); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}

func putExercise(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var ex models.Exercise
	if err := c.ShouldBindJSON(&ex); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ex.ID = id
	// Replace: delete old, insert new
	if err := dataStore.DeleteEx(id); err != nil {
		log.Println("WARN putExercise DeleteEx:", err)
	}
	if err := dataStore.InsertEx(ex); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func deleteExercise(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := dataStore.DeleteEx(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// ─── sets ─────────────────────────────────────────────────────────────────────

func getSets(c *gin.Context) {
	sets, err := dataStore.SelectSet()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sets)
}

func putSetsByDate(c *gin.Context) {
	date := c.Param("date")
	var sets []models.Set
	if err := c.ShouldBindJSON(&sets); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := dataStore.BulkReplaceSetsByDate(date, sets); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// ─── weight ───────────────────────────────────────────────────────────────────

func getWeight(c *gin.Context) {
	ws, err := dataStore.SelectW()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ws)
}

func postWeight(c *gin.Context) {
	var w models.BodyWeight
	if err := c.ShouldBindJSON(&w); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := dataStore.InsertW(w); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}

func deleteWeight(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := dataStore.DeleteW(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// ─── config ───────────────────────────────────────────────────────────────────

func getConfig(c *gin.Context) {
	c.JSON(http.StatusOK, appConfig)
}

func putConfig(c *gin.Context) {
	var cfg models.Conf
	if err := c.ShouldBindJSON(&cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	appConfig.Host = cfg.Host
	appConfig.Port = cfg.Port
	appConfig.Theme = cfg.Theme
	appConfig.Color = cfg.Color
	appConfig.HeatColor = cfg.HeatColor
	appConfig.PageStep = cfg.PageStep
	conf.Write(appConfig, authConf)
	c.Status(http.StatusOK)
}

// ─── migration ────────────────────────────────────────────────────────────────

func migrateFromSQLite(c *gin.Context) {
	pgStore, ok := dataStore.(*store.PostgresStore)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "not using postgres store"})
		return
	}
	result, err := migrate.SQLiteToPostgres(appConfig.DBPath, pgStore)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func putConfigAuth(c *gin.Context) {
	var body struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Expire   string `json:"expire"`
		Auth     bool   `json:"auth"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	authConf.User = body.User
	authConf.ExpStr = body.Expire
	authConf.Auth = body.Auth
	appConfig.Auth = body.Auth
	if body.Password != "" {
		authConf.Password = auth.HashPassword(body.Password)
	}
	authConf.Expire = auth.ToTime(authConf.ExpStr)
	if authConf.Auth && (authConf.User == "" || authConf.Password == "") {
		log.Println("WARNING: Auth won't work with empty login or password.")
		authConf.Auth = false
		appConfig.Auth = false
	}
	conf.Write(appConfig, authConf)
	c.Status(http.StatusOK)
}

