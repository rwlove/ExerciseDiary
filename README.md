[![Publish Container Images](https://github.com/rwlove/ExerciseDiary/actions/workflows/container-publish.yml/badge.svg)](https://github.com/rwlove/ExerciseDiary/actions/workflows/container-publish.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/rwlove/WorkoutDiary)](https://goreportcard.com/report/github.com/rwlove/WorkoutDiary)

<h1>Workout Diary</h1>

Workout diary with GitHub-style year visualization. Log daily sets, track body weight, and visualize training history with intensity heatmaps.

- [Architecture](#architecture)
- [Quick start](#quick-start)
- [Configuration](#configuration)
- [API server options](#api-server-options)
- [Frontend options](#frontend-options)
- [Local network only](#local-network-only)
- [Thanks](#thanks)

## Architecture

Workout Diary runs as two independent services:

```
┌──────────────────────────┐        ┌──────────────────────────┐
│  workoutdiary-frontend  │─HTTP──▶│  workoutdiary-api       │
│  Web UI  (default :8080) │        │  JSON API  (default :8851│
└──────────────────────────┘        └───────────┬──────────────┘
                                                 │
                                            SQLite DB
```

| Service | Image | Description |
|---|---|---|
| API backend | `ghcr.io/rwlove/workoutdiary-api` | Owns the SQLite database, exposes a JSON REST API |
| Web frontend | `ghcr.io/rwlove/workoutdiary-frontend` | Serves the browser UI, talks to the API over HTTP |

## Quick start

```sh
docker compose up
```

Or run each service manually:

```sh
# Start the API backend (stores data in /data/WorkoutDiary)
docker run --name exdiary-api \
  -v ~/.dockerdata/WorkoutDiary:/data/WorkoutDiary \
  -p 8851:8851 \
  ghcr.io/rwlove/workoutdiary-api

# Start the web frontend
docker run --name exdiary-frontend \
  -e API_URL=http://<YOUR_HOST_IP>:8851 \
  -p 8080:8080 \
  ghcr.io/rwlove/workoutdiary-frontend
```

Then open **http://localhost:8080** in your browser.

## Configuration

Both services are configured exclusively via environment variables. No config file is required.

### API server (`workoutdiary-api`)

| Variable | Description | Default |
|---|---|---|
| `PORT` | Listen port | `8851` |
| `HOST` | Listen address | `0.0.0.0` |
| `DATA_DIR` | SQLite data directory (also settable via `-d` flag) | `/data/WorkoutDiary` |
| `API_KEY` | Require this value on every `X-Api-Key` request header; empty = no auth | `""` |
| `THEME` | Any [Bootswatch](https://bootswatch.com) theme (lowercase) or extras: `emerald`, `grass`, `grayscale`, `ocean`, `sand`, `wood` | `grass` |
| `COLOR` | Background: `light` or `dark` | `dark` |
| `HEATCOLOR` | Heatmap cell color | `#03a70c` |
| `PAGESTEP` | Rows per page | `10` |
| `TZ` | Timezone | `""` |

### Frontend server (`workoutdiary-frontend`)

| Variable | Description | Default |
|---|---|---|
| `PORT` | Listen port | `8080` |
| `API_URL` | Base URL of the API server | `http://localhost:8851` |
| `API_KEY` | `X-Api-Key` value sent to the API (must match API server `API_KEY`) | `""` |
| `NODE_PATH` | URL of a [node-bootstrap](https://github.com/aceberg/my-dockerfiles/tree/main/node-bootstrap) instance for offline use | `""` |
| `TZ` | Timezone | `""` |

## Local network only

By default the app loads themes, icons, and fonts from the internet. For an air-gapped setup, run the [node-bootstrap](https://github.com/aceberg/my-dockerfiles/tree/main/node-bootstrap) sidecar and set `NODE_PATH` on the frontend:

```sh
docker run --name node-bootstrap \
  -v ~/.dockerdata/icons:/app/icons \
  -p 8850:8850 \
  aceberg/node-bootstrap

docker run --name exdiary-frontend \
  -e API_URL=http://<YOUR_HOST_IP>:8851 \
  -e NODE_PATH=http://<YOUR_HOST_IP>:8850 \
  -p 8080:8080 \
  ghcr.io/rwlove/workoutdiary-frontend
```

## Features

- **Exercise library** — organize exercises into groups, store default weight/reps/intensity
- **Daily workout log** — autosaves on every change (no Save button)
- **Body weight tracking** — log weight and view a rolling chart
- **Heatmap history** — GitHub-style workout intensity and per-exercise color heatmaps
- **Stats page** — per-exercise intensity charts with period filtering
- **Dark mode by default** — full Bootstrap dark theme
- **PWA support** — installable as a home screen app

## Thanks

- All Go packages listed in [go.mod](go.mod)
- [Bootstrap](https://getbootstrap.com/) and [Bootswatch](https://bootswatch.com) themes
- [Chart.js](https://github.com/chartjs/Chart.js) and [chartjs-chart-matrix](https://github.com/kurkle/chartjs-chart-matrix)
- [Gin](https://github.com/gin-gonic/gin)
- Favicon and logo: [Flaticon](https://www.flaticon.com/icons/)
