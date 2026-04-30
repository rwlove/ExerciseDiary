package web

import (
	"embed"

	"github.com/aceberg/ExerciseDiary/internal/auth"
	"github.com/aceberg/ExerciseDiary/internal/models"
	"github.com/aceberg/ExerciseDiary/internal/store"
)

var (
	// appConfig - config for Web Gui
	appConfig models.Conf

	// authConf - config for auth
	authConf auth.Conf

	// Exercise data
	exData models.AllExData

	// dataStore is the active data source.
	// SQLiteStore in monolith mode, APIClient in split-frontend mode.
	dataStore store.Store

	// apiClient is non-nil only in split-frontend mode; used for config
	// operations that fall outside the Store interface.
	apiClient *store.APIClient
)

// templFS - html templates
//
//go:embed templates/*
var templFS embed.FS

// pubFS - public folder
//
//go:embed public/*
var pubFS embed.FS
