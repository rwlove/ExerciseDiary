[![Container Publish](https://github.com/rwlove/ExerciseDiary/actions/workflows/container-publish.yml/badge.svg)](https://github.com/rwlove/ExerciseDiary/actions/workflows/container-publish.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/rwlove/WorkoutDiary)](https://goreportcard.com/report/github.com/rwlove/WorkoutDiary)

<h1>Workout Diary</h1>

A modern workout tracking app with GitHub-style heatmap visualization, body weight logging, and per-exercise stats.

- [Quick start](#quick-start)
- [Architecture](#architecture)
- [Configuration](#configuration)
- [Container images](#container-images)

## Quick start

Run both services with Docker Compose:

```yaml
services:
  api:
    image: ghcr.io/rwlove/workoutdiary-api:latest
    environment:
      - DATA_DIR=/data
      - PORT=8851
      - API_KEY=changeme
    volumes:
      - ./data:/data
    ports:
      - "8851:8851"

  frontend:
    image: ghcr.io/rwlove/workoutdiary-frontend:latest
    environment:
      - PORT=8080
      - API_URL=http://api:8851
      - API_KEY=changeme
    ports:
      - "8080:8080"
    depends_on:
      - api
```

Then open `http://localhost:8080`.

## Architecture

Workout Diary is split into two independent services:

| Service | Binary | Default port | Description |
| ------- | ------ | ------------ | ----------- |
| API | `workoutdiary-api` | 8851 | SQLite database + JSON REST API |
| Frontend | `workoutdiary-frontend` | 8080 | Web UI, proxies data from the API |

## Configuration

All configuration is done via environment variables — no config files needed.

### API service

| Variable | Description | Default |
| -------- | ----------- | ------- |
| `DATA_DIR` | Directory for the SQLite database | `/data/WorkoutDiary` |
| `PORT` | Listen port | `8851` |
| `HOST` | Listen address | `0.0.0.0` |
| `API_KEY` | Secret key required by the frontend | `""` |
| `THEME` | Bootswatch theme name | `darkly` |
| `COLOR` | Background color: `light` or `dark` | `dark` |
| `HEATCOLOR` | Heatmap accent color (hex) | `#03a70c` |
| `TZ` | Timezone for correct date handling | `""` |

### Frontend service

| Variable | Description | Default |
| -------- | ----------- | ------- |
| `PORT` | Listen port | `8080` |
| `API_URL` | Base URL of the API service | `http://localhost:8851` |
| `API_KEY` | Must match the API service key | `""` |
| `NODE_PATH` | Optional URL for local Bootstrap/icons bundle | `""` |

## Container images

Pre-built multi-arch images (amd64, arm64, arm/v7) are published to GHCR on every push to `main` and on version tags:

```
ghcr.io/rwlove/workoutdiary-api:latest
ghcr.io/rwlove/workoutdiary-frontend:latest
```

## Features

- **Exercise library** — organize exercises into groups, store default weight/reps/intensity per exercise
- **Daily workout log** — autosaves on every change (no Save button needed)
- **Body weight tracking** — log weight and view a rolling chart
- **Heatmap history** — GitHub-style workout intensity and per-exercise color heatmaps
- **Stats page** — per-exercise intensity charts with period filtering
- **Dark mode by default** — full Bootstrap dark theme out of the box
- **PWA support** — installable as a home screen app

## Thanks

- All Go packages listed in [go.mod](go.mod)
- [Bootstrap](https://getbootstrap.com/) / [Bootswatch](https://bootswatch.com)
- [Chart.js](https://github.com/chartjs/Chart.js) and [chartjs-chart-matrix](https://github.com/kurkle/chartjs-chart-matrix)
- [Gin](https://github.com/gin-gonic/gin)
- Favicon and logo: [Flaticon](https://www.flaticon.com/icons/)
