# LogSentry

<div align="center">

```text
 _                 ____             _              
| |    ___   __ _ / ___|  ___ _ __ | |_ _ __ _   _ 
| |   / _ \ / _` |\___ \ / _ \ '_ \| __| '__| | | |
| |__| (_) | (_| | ___) |  __/ | | | |_| |  | |_| |
|_____\___/ \__, ||____/ \___|_| |_|\__|_|   \__, |
            |___/                             |___/ 
```

**High-performance distributed log processing, monitoring, analytics, and observability for Go systems.**

[![Go Version](https://img.shields.io/badge/Go-1.26.4-00ADD8?style=for-the-badge&logo=go)](https://go.dev/)
[![Docker](https://img.shields.io/badge/Docker-ready-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-17-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Prometheus](https://img.shields.io/badge/Prometheus-metrics-E6522C?style=for-the-badge&logo=prometheus&logoColor=white)](https://prometheus.io/)
[![Grafana](https://img.shields.io/badge/Grafana-dashboards-F46800?style=for-the-badge&logo=grafana&logoColor=white)](https://grafana.com/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow?style=for-the-badge)](#license)

[![GitHub Stars](https://img.shields.io/github/stars/amandx36/LogSentry?style=social)](https://github.com/amandx36/LogSentry/stargazers)
[![GitHub Forks](https://img.shields.io/github/forks/amandx36/LogSentry?style=social)](https://github.com/amandx36/LogSentry/network/members)
[![GitHub Issues](https://img.shields.io/github/issues/amandx36/LogSentry)](https://github.com/amandx36/LogSentry/issues)
[![GitHub Actions](https://img.shields.io/github/actions/workflow/status/amandx36/LogSentry/ci.yml?branch=main&label=CI)](https://github.com/amandx36/LogSentry/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/amandx36/LogSentry)](https://goreportcard.com/report/github.com/amandx36/LogSentry)
[![Release](https://img.shields.io/github/v/release/amandx36/LogSentry?include_prereleases)](https://github.com/amandx36/LogSentry/releases)
[![Maintained](https://img.shields.io/badge/Maintained-yes-brightgreen)](https://github.com/amandx36/LogSentry/commits/main)
[![Build Passing](https://img.shields.io/badge/build-passing-brightgreen)](#testing)
[![Code Coverage](https://img.shields.io/badge/coverage-project--level-blue)](#testing)

</div>

---

## The Problem

Most log processing setups in real backend systems do one of these badly:
- Reread entire log files on every check → wasted I/O, doesn't scale
- Lose all progress on restart → duplicate or missing data
- Parse sequentially → can't keep up with high log volume
- Give zero visibility into their own health → you don't know if the pipeline itself is broken until logs go missing

## The Solution

LogSentry is a distributed log processing pipeline in Go that fixes each of these directly:

| Problem | LogSentry's Fix |
|---|---|
| Rereading whole files | Persistent byte-offset tracking — only reads bytes appended since last read |
| Lost progress on crash/restart | Offsets persisted to disk (JSON), so processing resumes exactly where it left off |
| Slow sequential parsing | CPU-core-sized worker pool (goroutines + channels) parses files concurrently |
| No visibility into pipeline health | Self-instrumented with Prometheus metrics + Grafana dashboards |
| Manual setup pain | One-command Docker Compose stack: Postgres + API + worker + Prometheus + Grafana |

In short: **it watches logs in real time, processes only new data, stores it queryable in Postgres, and exposes both the logs and its own operational health — all containerized.**


## Table of Contents

- [Project Overview](#project-overview)
- [Features](#features)
- [Screenshots](#screenshots)
- [Demo GIF](#demo-gif)
- [System Architecture](#system-architecture)
- [Complete Project Workflow](#complete-project-workflow)
- [Folder Structure](#folder-structure)
- [Component Responsibilities](#component-responsibilities)
- [API Documentation](#api-documentation)
- [Prometheus Metrics](#prometheus-metrics)
- [Grafana Dashboards](#grafana-dashboards)
- [Database Schema](#database-schema)
- [Worker Pool Design](#worker-pool-design)
- [Offset Persistence](#offset-persistence)
- [Docker Deployment](#docker-deployment)
- [Installation Guide](#installation-guide)
- [Running Without Docker](#running-without-docker)
- [Running With Docker Compose](#running-with-docker-compose)
- [Configuration](#configuration)
- [Monitoring](#monitoring)
- [Performance](#performance)
- [Production Improvements](#production-improvements)
- [Challenges Solved](#challenges-solved)
- [Lessons Learned](#lessons-learned)
- [Why Recruiters Should Care](#why-recruiters-should-care)
- [Contribution Guide](#contribution-guide)
- [Development Workflow](#development-workflow)
- [Testing](#testing)
- [License](#license)
- [Acknowledgements](#acknowledgements)
- [Contact](#contact)

---

## Project Overview

LogSentry is a distributed log processing and monitoring system written in Go. It watches log files in real time, processes only newly appended bytes using persistent file offsets, parses logs concurrently with a worker pool, stores structured records in PostgreSQL, exposes operational metrics to Prometheus, and visualizes system health through Grafana dashboards.

The project models a production backend observability pipeline:

1. Applications write log lines to files.
2. LogSentry detects file writes through `fsnotify`.
3. The offset manager reads only bytes that were appended since the last successful read.
4. The parser converts raw lines into structured log entries.
5. PostgreSQL stores searchable records.
6. REST APIs expose search, dashboard, analytics, and health endpoints.
7. Prometheus scrapes runtime metrics.
8. Grafana turns those metrics into operational dashboards.

### Real-World Motivation

Application logs are usually high-volume, append-only, semi-structured data. A naive log processor often rereads entire files, loses progress after restart, blocks on sequential parsing, or provides no visibility into its own failure modes. LogSentry exists to solve those problems with a clear backend architecture:

- **Incremental processing** avoids duplicate ingestion and wasted I/O.
- **Persistent offsets** make restarts recoverable.
- **Worker pools** use CPU cores efficiently during batch parsing.
- **Structured storage** enables reliable search and analytics.
- **Metrics and dashboards** expose the health of the pipeline itself.
- **Docker Compose** makes the system reproducible for local and team environments.

### Production Use Cases

LogSentry can be used as a learning-grade or prototype foundation for:

- Backend service log ingestion.
- Local observability stacks.
- Incident debugging tools.
- Developer environment monitoring.
- Log analytics dashboards.
- API-driven operational reporting.
- Demonstrations of Go concurrency, observability, and containerized backend design.

### Architecture Philosophy

LogSentry follows a modular architecture where every package owns a narrow responsibility:

- `monitor` watches files and manages offsets.
- `parser` converts raw log lines into domain models.
- `repository/postgres` persists logs.
- `services` contain query, analytics, and dashboard logic.
- `api` exposes HTTP endpoints.
- `metrics` centralizes Prometheus collectors.
- `config` loads runtime settings.

This keeps I/O, parsing, persistence, API handling, and metrics decoupled enough to evolve independently.

---

## Features

### Live File Monitoring

LogSentry uses `fsnotify` to watch the configured input directory. Whenever a log file is written to, the watcher receives an event and triggers incremental processing.

**Why it exists:** production logs are append-heavy and continuous. A monitoring system must react to new data without requiring manual restarts or full directory rescans.

### Persistent Offset Management

The offset manager tracks the last byte position processed for every watched file and persists this state in JSON.

**Why it exists:** if the worker restarts, LogSentry can continue from the last known offset instead of rereading the whole file or missing appended data.

### Worker Pool Architecture

Batch parsing uses a fixed number of workers based on `runtime.NumCPU()`. Files are sent through a jobs channel, parsed concurrently, and merged into a final report.

**Why it exists:** log directories can contain many files. A bounded worker pool improves throughput while avoiding unbounded goroutine creation.

### Parallel Log Parsing

Each worker opens a `.log` file, parses it with a regular expression, categorizes entries, emits metrics, and sends results back through a results channel.

**Why it exists:** parsing is CPU and I/O sensitive. Parallelism improves initial ingestion speed and demonstrates safe Go concurrency.

### PostgreSQL Storage

Parsed entries are inserted into the `log_entries` table using batched SQL inserts.

**Why it exists:** structured storage enables API queries, analytics, dashboard totals, and durable retrieval.

### Prometheus Metrics

Both the API server and worker expose Prometheus-compatible metrics. Metrics include processed logs, parser failures, watcher events, DB inserts, active workers, and parse duration.

**Why it exists:** a monitoring service must also monitor itself. These metrics help identify throughput, parsing quality, ingest health, and runtime behavior.

### Grafana Dashboards

Grafana dashboard JSON files are included under `dashBoard/` for processing metrics, error monitoring, live activity, and log analysis.

**Why it exists:** dashboards make the system readable at a glance and support operational debugging.

### Search APIs

LogSentry exposes APIs for category search, source search, keyword search, date range search, and recent logs.

**Why it exists:** logs are valuable only when engineers can retrieve them quickly during debugging and incident response.

### Dashboard APIs

The `/dashboard` endpoint returns aggregate totals by category.

**Why it exists:** dashboards and UI clients need compact summary data without repeatedly running custom queries.

### Analytics APIs

The `/analytics` endpoint returns total volume, category counts, error rate, top source, and generation time.

**Why it exists:** operational teams need derived insights, not only raw log rows.

### Category Filtering

Logs can be filtered by category such as `ERROR`, `WARN`, `INFO`, or unknown categories.

**Why it exists:** severity filtering is the fastest way to narrow a noisy log stream.

### Source Filtering

Logs can be filtered by the service or component source captured inside square brackets in each log line.

**Why it exists:** production systems usually contain multiple services, and engineers need to isolate one component at a time.

### Keyword Search

The keyword endpoint performs case-insensitive matching against log details using PostgreSQL `ILIKE`.

**Why it exists:** incident debugging often starts with a partial error string, request ID, timeout marker, or domain keyword.

### Recent Logs

The recent logs endpoint returns newest records ordered by descending database ID.

**Why it exists:** live debugging frequently needs the latest events first.

### Error Tracking

Error counters, category counts, and analytics calculations all track error volume.

**Why it exists:** error rate is one of the clearest signals of application health.

### Metrics Monitoring

Prometheus scrapes both `server:8080/metrics` and `worker:9091/metrics`.

**Why it exists:** separating API metrics from worker metrics makes it easier to distinguish API health from ingestion health.

### Containerized Deployment

Docker Compose starts PostgreSQL, the API server, worker, Prometheus, and Grafana on a shared bridge network.

**Why it exists:** contributors and reviewers can run the full stack with one command and consistent service names.

### Modular Architecture

The codebase is organized around packages with clear responsibility boundaries.

**Why it exists:** modularity reduces coupling and makes future production upgrades, such as Kafka or OpenTelemetry, easier.

---

## Screenshots

> [!NOTE]
> Screenshot assets can be replaced as the UI and dashboards evolve. Existing dashboard image assets are available under `dashBoard/images/`.

| Area | Image |
|---|---|
| Application Dashboard | `images/dashboard.png` |
| Grafana Overview | `images/grafana.png` |
| API Response | `images/api.png` |
| Monitoring View | `images/monitoring.png` |
| Processing Metrics | `dashBoard/images/processingMatrix.jpeg` |
| Error Monitoring | `dashBoard/images/ErrorMoniotring.jpeg` |
| Log Analysis | `dashBoard/images/LogAnalysis.jpeg` |
| Live Activity | `dashBoard/images/LogActivity.jpeg` |

---

## Grafana Dashboard

> [!TIP]
> Add a short recording showing a log file being appended, LogSentry ingesting the new lines, Prometheus metrics changing, and Grafana updating.

Expected demo asset:

```text
images/demo.gif
```

---

## System Architecture

### Overall Architecture

```mermaid
flowchart LR
    A[Application Log Files] --> B[fsnotify Watcher]
    B --> C[Offset Manager]
    C --> D[Read Newly Appended Bytes]
    D --> E[Parser]
    E --> F[Structured Log Report]
    F --> G[Batch Insert]
    G --> H[(PostgreSQL)]
    H --> I[REST API Server]
    I --> J[Search / Dashboard / Analytics Clients]
    E --> K[Prometheus Metrics]
    B --> K
    G --> K
    K --> L[Prometheus]
    L --> M[Grafana Dashboards]
```

### Docker Architecture

```mermaid
flowchart TB
    subgraph Host[Developer Machine]
        L[./logs volume]
        O[./internal/data volume]
        P[./configs/prometheus.yml]
    end

    subgraph Net[Docker Bridge Network: LogsentryNet]
        DB[(logsentry-DB<br/>PostgreSQL 17)]
        S[logsentry-server<br/>Gin API :8080]
        W[logsentry-worker<br/>Worker + Metrics :9091]
        PR[Prometheus :9090]
        G[Grafana :3000]
    end

    L --> W
    O --> W
    P --> PR
    S --> DB
    W --> DB
    PR --> S
    PR --> W
    G --> PR
```

### Data Flow

```mermaid
sequenceDiagram
    participant App as Application
    participant File as Log File
    participant Watcher as fsnotify Watcher
    participant Offset as Offset Manager
    participant Reader as Incremental Reader
    participant Parser as Parser
    participant DB as PostgreSQL
    participant API as REST API

    App->>File: Append log lines
    File-->>Watcher: Write event
    Watcher->>Offset: Get last processed byte
    Watcher->>Reader: Read from last offset
    Reader-->>Watcher: New bytes + new offset
    Watcher->>Offset: Save new offset
    Watcher->>Parser: Parse appended bytes
    Parser-->>Watcher: Structured report
    Watcher->>DB: Batch insert log entries
    API->>DB: Query logs and aggregates
    DB-->>API: JSON-ready results
```

### Worker Flow

```mermaid
flowchart TD
    A[Scan input_dir] --> B[Filter .log files]
    B --> C[Jobs Channel]
    C --> W1[Worker 1]
    C --> W2[Worker 2]
    C --> WN[Worker N]
    W1 --> R[Results Channel]
    W2 --> R
    WN --> R
    R --> M[MergeReports]
    M --> F[Final LogReport]
    F --> DB[Batch Insert]
```

### Monitoring Flow

```mermaid
flowchart LR
    Worker[Worker Runtime] --> WM[/worker:9091 metrics/]
    Server[API Server Runtime] --> SM[/server:8080 metrics/]
    WM --> Prom[Prometheus]
    SM --> Prom
    Prom --> Grafana[Grafana Dashboards]
```

### Prometheus Scraping Flow

```mermaid
flowchart TD
    C[configs/prometheus.yml] --> P[Prometheus]
    P -->|scrape_interval: 5s| S[server:8080]
    P -->|scrape_interval: 5s| W[worker:9091]
    S --> M1[/metrics]
    W --> M2[/metrics]
    P --> TSDB[(Prometheus TSDB)]
    TSDB --> G[Grafana Panels]
```

### API Request Flow

```mermaid
sequenceDiagram
    participant Client
    participant Gin as Gin Router
    participant Controller
    participant Service
    participant Postgres

    Client->>Gin: GET /search/category?category=ERROR
    Gin->>Controller: Route to handler
    Controller->>Controller: Validate query params
    Controller->>Service: SearchByCategory(db, category)
    Service->>Postgres: SELECT ...
    Postgres-->>Service: Rows
    Service-->>Controller: []LogEntry
    Controller-->>Client: 200 JSON
```

### Log Parsing Flow

```mermaid
flowchart TD
    A[Raw log line] --> B{Regex match?}
    B -->|No| PF[parser_failures_total++]
    B -->|Yes| C[Extract timestamp, category, source, details]
    C --> D{Category}
    D -->|ERROR| E[Errors slice + logs_error_total]
    D -->|WARN| W[Warns slice + logs_warn_total]
    D -->|INFO| I[Infos slice + logs_info_total]
    D -->|Other| U[Unknown slice + logs_unknown_total]
    E --> T[logs_processed_total++]
    W --> T
    I --> T
    U --> T
```

### Offset Persistence Flow

```mermaid
flowchart TD
    A[Write event] --> B[GetOffset(file)]
    B --> C[Seek to last byte]
    C --> D[Read new bytes]
    D --> E[Seek end to compute new offset]
    E --> F[UpdateOffset(file, newOffset)]
    F --> G[Marshal offset map to JSON]
    G --> H[Write internal/data/offset.json]
```

---

## Complete Project Workflow

```mermaid
flowchart TD
    A[Log File] --> B[Watcher]
    B --> C[Offset Manager]
    C --> D[Read New Bytes]
    D --> E[Worker Pool / Parser]
    E --> F[Merge]
    F --> G[Database]
    G --> H[Metrics]
    G --> I[REST API]
    H --> J[Prometheus]
    J --> K[Grafana]
    I --> L[User]
    K --> L
```

1. **Log File**: application services append lines to files inside `logs/inputs`.
2. **Watcher**: `fsnotify` detects write events in the configured input directory.
3. **Offset Manager**: LogSentry retrieves the last processed byte for the changed file.
4. **Read New Bytes**: the reader seeks to the offset and reads only appended data.
5. **Worker Pool**: initial directory parsing uses a CPU-sized worker pool for `.log` files.
6. **Parser**: raw lines are transformed into `models.LogEntry` values.
7. **Merge**: partial worker reports are merged into a single report.
8. **Database**: entries are batch inserted into PostgreSQL.
9. **Metrics**: parser, watcher, worker, and database activity are recorded as Prometheus metrics.
10. **REST API**: Gin routes expose search, dashboard, analytics, health, version, and metrics endpoints.
11. **Grafana**: dashboards visualize Prometheus time series.
12. **User**: engineers inspect logs, metrics, and operational summaries.

---

## Folder Structure

```text
LogSentry/
├── cmd/
│   ├── server/                  # Gin API process
│   └── worker/                  # Parser, file watcher, worker metrics process
├── configs/
│   └── prometheus.yml           # Prometheus scrape configuration
├── dashBoard/
│   ├── *.json                   # Grafana dashboard exports
│   └── images/                  # Dashboard screenshots
├── docker/
│   ├── server.Dockerfile        # API container build
│   └── worker.Dockerfile        # Worker container build
├── docs/
│   ├── WorkerPool.md            # Worker pool theory notes
│   ├── WorkerPoolImpl.md        # Worker pool implementation notes
│   └── images/                  # Documentation images
├── internal/
│   ├── api/
│   │   ├── controller/          # HTTP handlers
│   │   └── routes/              # Route registration
│   ├── config/                  # JSON configuration loading
│   ├── data/                    # Offset persistence data
│   ├── metrics/                 # Prometheus collectors
│   ├── models/                  # Domain models and schema setup
│   ├── monitor/                 # fsnotify watcher, reader, offsets
│   ├── parser/                  # File and byte parsing, worker pool
│   ├── repository/postgres/     # PostgreSQL connection and inserts
│   ├── services/
│   │   ├── analytics/           # Analytics queries
│   │   ├── dashboard/           # Dashboard aggregate queries
│   │   └── search/              # Search queries
│   └── writer/                  # File output writer
├── logs/                        # Mounted log directory
├── docker-compose.yml           # Full local observability stack
├── go.mod
├── go.sum
└── README.md
```

### Why This Structure Exists

| Path | Responsibility | Why it exists |
|---|---|---|
| `cmd/server` | API entrypoint | Separates the HTTP service process from worker logic. |
| `cmd/worker` | Worker entrypoint | Runs initial parsing, live monitoring, and worker metrics. |
| `internal/api/controller` | HTTP handlers | Keeps request validation and HTTP responses near Gin. |
| `internal/api/routes` | Route registration | Centralizes URL mapping and makes endpoint discovery simple. |
| `internal/config` | Runtime config loading | Provides a single source for directories, buffers, and DB connection. |
| `internal/metrics` | Prometheus collectors | Avoids duplicate metric definitions across packages. |
| `internal/models` | Domain structs and schema | Defines shared data contracts. |
| `internal/monitor` | Live file monitoring | Owns offsets, watcher events, and incremental reads. |
| `internal/parser` | Parsing and worker pool | Owns conversion from raw logs to structured reports. |
| `internal/repository/postgres` | Persistence | Keeps SQL connection and insert logic out of services. |
| `internal/services` | Business queries | Encapsulates search, analytics, and dashboard calculations. |
| `dashBoard` | Grafana assets | Stores importable dashboard definitions and screenshots. |
| `docker` | Container builds | Keeps Dockerfiles separate from application source. |
| `configs` | Infrastructure config | Stores Prometheus scrape settings. |

---

## Component Responsibilities

### `controller`

The controller package contains Gin handlers. It validates query parameters, calls service functions, handles service errors, and writes JSON responses.

Key handlers:

- `Ping`: health response at `/ping`.
- `Apiversion`: version response at `/api/version`.
- `GetDashboard`: aggregate dashboard stats.
- `GetAnalyticsEP`: analytics summary.
- `GetRecentLog`: recent log retrieval with optional limit.
- `SearchByCategory`: category filter.
- `SearchBySource`: source filter.
- `SearchByKeyWord`: keyword filter.
- `SearchByDateEP`: date range filter.

### `routes`

The routes package maps HTTP paths to controllers. `RegestAllRoutes` registers every public endpoint on the Gin engine.

### `parser`

The parser package owns both batch and incremental parsing.

- `ParseSingleFile` parses an opened file.
- `ParseByte` parses appended byte slices.
- `ScanAndSendJobs` scans `input_dir` and sends `.log` files to workers.
- `Worker` consumes jobs and sends parsed reports.
- `MergeReports` combines category slices and counts.
- `LoadingBuffer` coordinates workers, results, and dashboard totals.

The log format expected by the parser is:

```text
<timestamp date> <timestamp time> <CATEGORY> [<source>] <details>
```

Example:

```text
2026-07-07 10:30:00 ERROR [payments] timeout while processing invoice
```

### `repository`

The PostgreSQL repository package handles:

- `Connect`: open and ping a PostgreSQL connection using `pgx`.
- `CreateTables`: create `log_entries` if it does not exist.
- `BatchInsert`: insert multiple log entries in one transaction.
- `InsertLogs`: insert all category slices from a parsed report.

### `monitor`

The monitor package handles live ingestion:

- `OffsetManager`: thread-safe map of file offsets.
- `LoadOffsets`: restore offsets from JSON.
- `Save`: persist offsets to disk.
- `ReadNewLogs`: seek and read only appended bytes.
- `DirWatching`: listen for write events and process new data.
- `ProcessLogs`: parse appended bytes and insert them into PostgreSQL.

### `metrics`

The metrics package defines Prometheus collectors used throughout the application. Keeping these in one package prevents mismatched names and duplicate registration.

### `writer`

The writer package writes categorized parser output to files. This is useful for development inspection and local debugging.

### `dashboard`

The dashboard service queries PostgreSQL category counts and returns `DashBoardDetails`.

### `analytics`

The analytics service computes higher-level insights such as total logs, error rate, warning counts, top source, and generation time.

### `search`

The search service contains reusable SQL query helpers for category, source, keyword, recent, and date range searches.

### `models`

The models package defines shared data structures:

- `LogEntry`
- `LogReport`
- `DashBoardDetails`
- `Analytics`
- `MonitringConfig`

It also initializes the database schema.

### `config`

The config package loads `internal/config/config.json` and maps JSON fields into a typed Go configuration struct.

---

## API Documentation

Base URL when running locally:

```text
http://localhost:8080
```

### Endpoint Summary

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/ping` | Health check. |
| `GET` | `/api/version` | API version and endpoint status. |
| `GET` | `/metrics` | Prometheus metrics from the API server. |
| `GET` | `/dashboard` | Aggregate category counts. |
| `GET` | `/analytics` | Derived analytics summary. |
| `GET` | `/logs/recent?limit=10` | Latest logs ordered by descending ID. |
| `GET` | `/search/category?category=ERROR` | Logs matching a category. |
| `GET` | `/search/source?source=payments` | Logs matching a source. |
| `GET` | `/search/keywords?keyword=timeout` | Logs whose details contain a keyword. |
| `GET` | `/search/date?start_date=2026-06-24&end_date=2026-06-26` | Logs within a timestamp range. |

### `GET /ping`

Health check endpoint.

```bash
curl http://localhost:8080/ping
```

Response:

```json
{
  "message": "pong"
}
```

| Status | Meaning |
|---|---|
| `200` | Server is reachable. |

### `GET /api/version`

Returns the API version.

```bash
curl http://localhost:8080/api/version
```

Response:

```json
{
  "Version": "1.0.0",
  "EndPoints": "Working"
}
```

| Status | Meaning |
|---|---|
| `200` | Version returned successfully. |

### `GET /metrics`

Exposes API server metrics in Prometheus text format.

```bash
curl http://localhost:8080/metrics
```

| Status | Meaning |
|---|---|
| `200` | Prometheus metrics returned. |

### `GET /dashboard`

Returns aggregate category totals.

```bash
curl http://localhost:8080/dashboard
```

Example response:

```json
{
  "TotalLogs": 1250,
  "Errors": 42,
  "Warns": 188,
  "Infos": 990,
  "Unknown": 30
}
```

| Status | Meaning |
|---|---|
| `200` | Dashboard totals returned. |
| `500` | Failed to query aggregate statistics. |

### `GET /analytics`

Returns derived analytics.

```bash
curl http://localhost:8080/analytics
```

Example response:

```json
{
  "TotalLogs": 1250,
  "TotalErrors": 42,
  "TotalWarns": 188,
  "TotalInfos": 990,
  "TotalUnknown": 30,
  "ErrorRate": 3.36,
  "WarningRate": 0,
  "TopSource": "payments",
  "TopSourceCount": 320,
  "GeneratedAt": "2026-07-07T10:30:00.000000000+05:30"
}
```

| Status | Meaning |
|---|---|
| `200` | Analytics returned. |
| `500` | Failed to compute analytics. |

### `GET /logs/recent`

Returns newest logs. Default limit is `10`.

```bash
curl "http://localhost:8080/logs/recent?limit=5"
```

Example response:

```json
[
  {
    "TimeStamp": "2026-07-07 10:30:00",
    "Category": "ERROR",
    "Source": "payments",
    "Details": "timeout while processing invoice"
  }
]
```

| Query Parameter | Required | Description |
|---|---:|---|
| `limit` | No | Number of recent records to return. Defaults to `10`. |

| Status | Meaning |
|---|---|
| `200` | Logs returned. |
| `400` | `limit` is not an integer. |
| `500` | Query failed. |

### `GET /search/category`

Searches by exact category.

```bash
curl "http://localhost:8080/search/category?category=ERROR"
```

Example response:

```json
[
  {
    "TimeStamp": "2026-07-07 10:30:00",
    "Category": "ERROR",
    "Source": "payments",
    "Details": "timeout while processing invoice"
  }
]
```

| Query Parameter | Required | Description |
|---|---:|---|
| `category` | Yes | Category such as `ERROR`, `WARN`, `INFO`, or another parsed value. |

| Status | Meaning |
|---|---|
| `200` | Matching logs returned. |
| `400` | Missing `category`. |
| `500` | Query failed. |

### `GET /search/source`

Searches by exact source.

```bash
curl "http://localhost:8080/search/source?source=payments"
```

Example response:

```json
{
  "count": 1,
  "data": [
    {
      "TimeStamp": "2026-07-07 10:30:00",
      "Category": "ERROR",
      "Source": "payments",
      "Details": "timeout while processing invoice"
    }
  ]
}
```

| Query Parameter | Required | Description |
|---|---:|---|
| `source` | Yes | Source captured inside square brackets in the log line. |

| Status | Meaning |
|---|---|
| `200` | Matching logs returned. |
| `400` | Missing `source`. |
| `500` | Query failed. |

### `GET /search/keywords`

Searches log details with case-insensitive matching.

```bash
curl "http://localhost:8080/search/keywords?keyword=timeout"
```

Example response:

```json
{
  "count": 1,
  "data": [
    {
      "TimeStamp": "2026-07-07 10:30:00",
      "Category": "ERROR",
      "Source": "payments",
      "Details": "timeout while processing invoice"
    }
  ]
}
```

| Query Parameter | Required | Description |
|---|---:|---|
| `keyword` | Yes | Case-insensitive keyword matched against `details`. |

| Status | Meaning |
|---|---|
| `200` | Matching logs returned. |
| `400` | Missing `keyword`. |
| `500` | Query failed. |

### `GET /search/date`

Searches by timestamp range.

```bash
curl "http://localhost:8080/search/date?start_date=2026-06-24&end_date=2026-06-26"
```

Example response:

```json
[
  {
    "TimeStamp": "2026-06-25 09:00:00",
    "Category": "WARN",
    "Source": "gateway",
    "Details": "retrying upstream request"
  }
]
```

| Query Parameter | Required | Description |
|---|---:|---|
| `start_date` | Yes | Lower timestamp bound. |
| `end_date` | Yes | Upper timestamp bound. |

| Status | Meaning |
|---|---|
| `200` | Matching logs returned. |
| `400` | Missing `start_date` or `end_date`. |
| `500` | Query failed. |

> [!WARNING]
> The current schema stores `timestamp` as `TEXT`, so date range behavior depends on lexicographic ordering. Use ISO-like sortable timestamp strings for predictable results.

---

## Prometheus Metrics

| Metric | Type | Source | What it measures | Dashboard use |
|---|---|---|---|---|
| `logs_processed_total` | Counter | Parser | Total parsed log entries. | Throughput and ingestion volume. |
| `logs_error_total` | Counter | Parser | Parsed `ERROR` entries. | Error trend panels. |
| `logs_warn_total` | Counter | Parser | Parsed `WARN` entries. | Warning trend panels. |
| `logs_info_total` | Counter | Parser | Parsed `INFO` entries. | Normal activity panels. |
| `logs_unknown_total` | Counter | Parser | Parsed entries with unrecognized categories. | Parser quality and format drift. |
| `db_insert_total` | Counter | Repository | Number of rows inserted into PostgreSQL. | Persistence throughput. |
| `watcher_events_total` | Counter | Monitor | File watcher events received. | Live file activity. |
| `parser_failures_total` | Counter | Parser | Lines that failed parsing or scanner errors. | Bad log format detection. |
| `log_write_events_total` | Counter | Monitor | Write events detected by watcher. | Live append activity. |
| `read_failures_total` | Counter | Monitor | Failures reading appended file data. | File I/O health. |
| `files_processed_total` | Counter | Parser | `.log` files submitted for parsing. | Batch ingestion progress. |
| `active_workers` | Gauge | Parser | Current active worker goroutines. | Worker pool visibility. |
| `parse_duration_seconds` | Histogram | Parser | Time spent parsing file or byte input. | Latency distribution. |

### Example PromQL

```promql
rate(logs_processed_total[1m])
```

```promql
rate(logs_error_total[5m])
```

```promql
histogram_quantile(0.95, rate(parse_duration_seconds_bucket[5m]))
```

```promql
active_workers
```

---

## Grafana Dashboards

Dashboard definitions live in `dashBoard/`.

| Dashboard | File | Purpose |
|---|---|---|
| Processing Metrics | `dashBoard/ProcessingMatrix.json` | Tracks parsing throughput, active workers, DB inserts, and parse duration. |
| Error Monitoring | `dashBoard/ErrorMOnitoring.json` | Focuses on error logs, parser failures, and read failures. |
| Live Activity | `dashBoard/LiveActivity.json` | Shows watcher events and live write activity. |
| Log Analytics | `dashBoard/LogAnalysis.json` | Visualizes category distribution and ingestion trends. |

### Processing Metrics

Use this dashboard to answer:

- Is the worker parsing logs?
- Are inserts reaching PostgreSQL?
- Are active workers visible?
- Is parsing latency increasing?

Screenshot:


![alt text](ProcessingMatrix.png)


### Error Monitoring

Use this dashboard to answer:

- Are application error logs rising?
- Are parser failures increasing?
- Are file reads failing?

Screenshot:

![alt text](ErrorMonitoring.png)

### Live Activity

Use this dashboard to answer:

- Is `fsnotify` receiving write events?
- Are log files still active?
- Is the live monitoring loop processing appends?

Screenshot:

![alt text](LiveActivity.png)

### Log Analytics

Use this dashboard to answer:

- Which categories dominate the log stream?
- Is normal activity stable?
- Are unknown logs appearing due to format drift?

Screenshot:

![alt text](LogAnalytics.png)

---

## Database Schema

LogSentry creates a single table on startup.

```sql
CREATE TABLE IF NOT EXISTS log_entries (
    id SERIAL PRIMARY KEY,
    timestamp TEXT NOT NULL,
    category VARCHAR(10) NOT NULL,
    source VARCHAR(100),
    details TEXT NOT NULL
);
```

### ER Diagram

```mermaid
erDiagram
    LOG_ENTRIES {
        int id PK
        text timestamp
        varchar category
        varchar source
        text details
    }
```

### Field Details

| Column | Type | Description |
|---|---|---|
| `id` | `SERIAL PRIMARY KEY` | Monotonic database identifier. |
| `timestamp` | `TEXT NOT NULL` | Timestamp extracted from the log line. |
| `category` | `VARCHAR(10) NOT NULL` | Severity or category such as `ERROR`, `WARN`, `INFO`. |
| `source` | `VARCHAR(100)` | Source extracted from square brackets. |
| `details` | `TEXT NOT NULL` | Remaining log message. |

### Recommended Indexes

The current schema creates the primary key automatically. For larger datasets, add indexes for common search patterns:

```sql
CREATE INDEX IF NOT EXISTS idx_log_entries_category
ON log_entries (category);

CREATE INDEX IF NOT EXISTS idx_log_entries_source
ON log_entries (source);

CREATE INDEX IF NOT EXISTS idx_log_entries_timestamp
ON log_entries (timestamp);

CREATE INDEX IF NOT EXISTS idx_log_entries_details_trgm
ON log_entries USING gin (details gin_trgm_ops);
```

> [!NOTE]
> The trigram index requires the PostgreSQL `pg_trgm` extension.

```sql
CREATE EXTENSION IF NOT EXISTS pg_trgm;
```

---

## Worker Pool Design

The worker pool is used during initial batch parsing. It is built from standard Go concurrency primitives:

- A buffered `jobs` channel.
- A buffered `results` channel.
- `runtime.NumCPU()` workers.
- A `sync.WaitGroup`.
- A final merge step.

### Lifecycle

```mermaid
sequenceDiagram
    participant Main
    participant Jobs
    participant Worker
    participant Results
    participant Merge

    Main->>Jobs: Send .log file jobs
    Main->>Jobs: Close jobs channel
    loop for each worker
        Worker->>Jobs: Receive file path
        Worker->>Worker: ParseSingleFile
        Worker->>Results: Send Result
    end
    Main->>Main: WaitGroup waits
    Main->>Results: Close results channel
    Main->>Merge: Merge reports
```

### Advantages

- Bounded concurrency.
- Better CPU utilization.
- Simple synchronization through channels.
- Clean separation of scanning, parsing, and merging.
- No shared mutable parser state inside workers.

### Synchronization

Workers do not mutate the final report directly. Each worker returns a `Result`, and the main goroutine merges those results. This avoids race conditions around report slices and counters.

### Failure Handling

If a worker cannot open or parse a file, it sends a `Result` with an error. The coordinator skips failed results and continues processing the remaining files.

---

## Offset Persistence

LogSentry persists offsets in JSON using a map:

```json
{
  "logs/inputs/app.log": 2048
}
```

### Why JSON

JSON is simple, human-readable, portable, and sufficient for a local offset registry. It also makes debugging easy because contributors can inspect the exact byte offset for each file.

### Recovery Process

1. Worker starts.
2. `LoadOffsets` reads the JSON file.
3. If the file does not exist, an empty offset map is created.
4. When a write event occurs, the watcher fetches the stored offset.
5. The reader seeks to that byte position.
6. New bytes are parsed and inserted.
7. The new end-of-file offset is saved back to JSON.

### Crash Recovery

If the process crashes after saving an offset but before database insertion, a small amount of data may be skipped. If it crashes after insertion but before offset save, data may be reread and duplicated. A production version should make offset commits transactional with persistence or use an external durable queue.

> [!WARNING]
> The worker currently loads `internal/data/data.json` during startup but saves live offsets to `internal/data/offset.json` inside the watcher. Align these paths before relying on restart recovery in production.

---

## Docker Deployment

The Compose stack defines five services.

| Service | Container | Port | Purpose |
|---|---|---:|---|
| `db` | `logsentry-DB` | `5432` | PostgreSQL database. |
| `server` | `logsentry-server` | `8080` | Gin REST API and API metrics. |
| `worker` | `logsentry-worker` | `9091` | Initial parsing, live monitoring, worker metrics. |
| `prometheus` | `Prometheus` | `9090` | Metrics scraping and storage. |
| `grafana` | `Grafana` | `3000` | Dashboards. |

### Network

All services run on the `LogsentryNet` bridge network. Containers communicate using service names:

- API server uses `db:5432`.
- Worker uses `db:5432`.
- Prometheus scrapes `server:8080` and `worker:9091`.
- Grafana connects to Prometheus.

### Volumes

| Volume | Mounted at | Purpose |
|---|---|---|
| `postgres-data` | `/var/lib/postgresql/data` | Durable PostgreSQL data. |
| `prometheus-data` | `/prometheus` | Prometheus TSDB data. |
| `grafana-data` | `/var/lib/grafana` | Grafana settings and dashboards. |
| `./logs` | `/app/logs` in worker | Local log files. |
| `./internal/data` | `/app/internal/data` in worker | Offset persistence. |

---

## Installation Guide

### Prerequisites

| Tool | Recommended version | Required for |
|---|---|---|
| Go | `1.26.4` as declared in `go.mod` | Running binaries and tests locally. |
| Docker | Current stable | Containerized stack. |
| Docker Compose | Compose v2 | Multi-container deployment. |
| PostgreSQL | 17 or compatible | Local non-Docker database. |
| Prometheus | Current stable | Manual metrics setup. |
| Grafana | OSS current stable | Manual dashboard setup. |

### Linux

```bash
git clone git@github.com:amandx36/LogSentry.git
cd LogSentry
go mod download
docker compose up --build
```

### macOS

```bash
git clone git@github.com:amandx36/LogSentry.git
cd LogSentry
go mod download
docker compose up --build
```

### Windows PowerShell

```powershell
git clone git@github.com:amandx36/LogSentry.git
cd LogSentry
go mod download
docker compose up --build
```

---

## Running Without Docker

1. Start PostgreSQL.
2. Create a database named `logsentry`.
3. Update `internal/config/config.json` with a local connection string.

Example:

```json
{
  "output_dir": "logs/output/",
  "buffer_size": 4096,
  "input_dir": "logs/inputs",
  "connectionString": "postgres://admin:admin@localhost:5432/logsentry?sslmode=disable"
}
```

4. Install dependencies.

```bash
go mod download
```

5. Run the API server.

```bash
go run ./cmd/server
```

6. Run the worker in another terminal.

```bash
go run ./cmd/worker
```

7. Append a log line.

```bash
printf '2026-07-07 10:30:00 ERROR [payments] timeout while processing invoice\n' >> logs/inputs/app.log
```

8. Query recent logs.

```bash
curl "http://localhost:8080/logs/recent?limit=5"
```

---

## Running With Docker Compose

Start the full stack:

```bash
docker compose up --build
```

Expected containers:

```text
logsentry-DB
logsentry-server
logsentry-worker
Prometheus
Grafana
```

Expected URLs:

| Service | URL |
|---|---|
| API | `http://localhost:8080` |
| API metrics | `http://localhost:8080/metrics` |
| Worker metrics | `http://localhost:9091/metrics` |
| Prometheus | `http://localhost:9090` |
| Grafana | `http://localhost:3000` |
| PostgreSQL | `localhost:5432` |

Verify health:

```bash
curl http://localhost:8080/ping
```

Append a log line from the host:

```bash
mkdir -p logs/inputs
printf '2026-07-07 10:30:00 INFO [api] request completed successfully\n' >> logs/inputs/app.log
```

Check ingestion:

```bash
curl "http://localhost:8080/search/source?source=api"
```

---

## Configuration

### Application Config

File: `internal/config/config.json`

```json
{
  "output_dir": "logs/output/",
  "buffer_size": 4096,
  "input_dir": "logs/inputs",
  "connectionString": "postgres://admin:admin@db:5432/logsentry?sslmode=disable"
}
```

| Field | Description |
|---|---|
| `output_dir` | Directory used by the writer package for generated output. |
| `buffer_size` | Buffer configuration value reserved for log processing behavior. |
| `input_dir` | Directory watched and scanned for `.log` files. |
| `connectionString` | PostgreSQL connection URL. In Docker, host is `db`. |

### Prometheus Config

File: `configs/prometheus.yml`

```yaml
global:
  scrape_interval: 5s

scrape_configs:
  - job_name: "server"
    static_configs:
      - targets:
          - server:8080

  - job_name: "worker"
    static_configs:
      - targets:
          - worker:9091
```

### Ports

| Port | Owner | Purpose |
|---:|---|---|
| `3000` | Grafana | Dashboard UI. |
| `5432` | PostgreSQL | Database access. |
| `8080` | Server | REST API and API metrics. |
| `9090` | Prometheus | Prometheus UI and query API. |
| `9091` | Worker | Worker metrics. |

---

## Monitoring

### Prometheus Setup

Prometheus is configured to scrape:

- `server:8080`
- `worker:9091`

Both expose `/metrics`.

### Grafana Setup

1. Open `http://localhost:3000`.
2. Add Prometheus as a data source.
3. Use `http://prometheus:9090` from inside Docker.
4. Import dashboard JSON files from `dashBoard/`.

### Dashboard Import

In Grafana:

1. Go to **Dashboards**.
2. Select **New**.
3. Select **Import**.
4. Upload one of the JSON files from `dashBoard/`.
5. Choose the Prometheus data source.
6. Save the dashboard.

---

## Performance

### Complexity

| Operation | Complexity | Notes |
|---|---|---|
| Scan input directory | `O(n)` | `n` is number of directory entries. |
| Parse log lines | `O(m)` | `m` is number of lines processed. |
| Merge reports | `O(k)` | `k` is number of parsed entries. |
| Category search | `O(n)` without index | Add index on `category` for scale. |
| Source search | `O(n)` without index | Add index on `source` for scale. |
| Keyword search | `O(n)` without trigram index | Use `pg_trgm` for production search. |

### Memory

The parser currently accumulates categorized slices in memory before batch insertion. This is straightforward and fast for moderate datasets. For very large logs, a streaming insertion model or bounded batch flushing would reduce memory pressure.

### CPU

Initial parsing uses `runtime.NumCPU()` workers, which maps worker count to available logical CPU cores. This is a practical default for CPU-bound parsing.

### Scalability

Current scale path:

1. Add database indexes.
2. Stream large files in chunks.
3. Make offset commits transactional.
4. Add queue-based buffering.
5. Split workers horizontally.
6. Move deployment to Kubernetes.

---

## Production Improvements

The current project is a strong backend foundation. Production deployments should consider:

- Kafka or Redpanda for durable log event buffering.
- Redis for short-lived cache and rate-limiting state.
- Elasticsearch or OpenSearch for full-text log search.
- OpenTelemetry traces and metrics.
- Kubernetes manifests and Helm charts.
- CI/CD with GitHub Actions.
- JWT authentication for APIs.
- RBAC for dashboard and search access.
- API rate limiting.
- Request IDs and structured API logs.
- Database migrations instead of startup schema creation.
- Timestamp storage as `TIMESTAMPTZ`.
- Partitioned log tables for large retention windows.
- Deduplication keys for crash-safe ingestion.
- Horizontal worker scaling.
- Cloud deployment on AWS, GCP, or Azure.
- Secrets management for database credentials.
- TLS for public endpoints.
- Readiness and liveness probes.

---

## Challenges Solved

### Concurrency

The worker pool coordinates scanning, parsing, result collection, and merging without shared mutable parser state.

### Offsets

Persistent byte offsets prevent full-file rereads during live monitoring and allow restart-aware processing.

### Live Monitoring

`fsnotify` removes the need for polling and responds quickly to log file writes.

### Parsing

The parser extracts timestamp, category, source, and details with one regular expression and routes entries into category-specific slices.

### Race Conditions

The offset manager uses `sync.RWMutex`, and report merging happens in a single coordinator path.

### Docker

Docker Compose wires the database, API, worker, Prometheus, and Grafana into a reproducible local stack.

### Metrics

Prometheus counters, gauges, and histograms expose both business signals and internal runtime behavior.

---

## Lessons Learned

LogSentry demonstrates practical backend engineering concepts:

- Designing append-only data pipelines.
- Using Go channels and goroutines safely.
- Persisting state for restart recovery.
- Building REST APIs with Gin.
- Writing batched SQL inserts.
- Exposing Prometheus metrics.
- Using Grafana dashboards for operational visibility.
- Running multi-service systems with Docker Compose.
- Separating packages by responsibility.
- Thinking about failure modes before production deployment.

---

## Why Someone  Should Care

LogSentry demonstrates skills that map directly to backend, platform, infrastructure, and observability roles:

| Skill | Evidence in project |
|---|---|
| Go backend development | Gin APIs, packages, domain models, worker process. |
| Concurrency | Worker pool, channels, goroutines, WaitGroup, mutex-protected offsets. |
| System design | Ingestion pipeline, API layer, storage layer, monitoring layer. |
| Databases | PostgreSQL schema, batch inserts, search queries, aggregate queries. |
| Observability | Prometheus metrics and Grafana dashboards. |
| DevOps | Dockerfiles, Docker Compose, service networking, persistent volumes. |
| Production mindset | Offsets, incremental reads, metrics, modular architecture, deployment docs. |
| API design | Search, dashboard, analytics, health, and version endpoints. |
| Documentation | Architecture diagrams, operational setup, and contribution workflows. |

---

## Contribution Guide

1. Fork the repository.
2. Create a feature branch.
3. Make focused changes.
4. Run tests and formatting.
5. Open a pull request with a clear description.

```bash
git checkout -b feat/add-new-metric
go fmt ./...
go test ./...
git commit -m "feat: add parser throughput metric"
git push origin feat/add-new-metric
```

### Contribution Standards

- Keep package responsibilities clear.
- Add tests for parser, service, and repository behavior when practical.
- Prefer small pull requests.
- Document new endpoints and metrics in this README.
- Do not commit local database data or generated logs.
- Keep Docker and non-Docker workflows working.

---

## Development Workflow

### Branch Naming

| Branch prefix | Use |
|---|---|
| `feat/` | New feature. |
| `fix/` | Bug fix. |
| `docs/` | Documentation-only change. |
| `test/` | Test-only change. |
| `refactor/` | Internal code improvement. |
| `chore/` | Tooling or maintenance. |

### Commit Naming

Use concise conventional commits:

```text
feat: add source histogram metric
fix: align offset load and save paths
docs: document docker deployment
test: add parser category coverage
refactor: split monitor processing flow
```

### Versioning

Use semantic versioning for releases:

```text
MAJOR.MINOR.PATCH
```

- `MAJOR`: breaking API or storage changes.
- `MINOR`: backward-compatible features.
- `PATCH`: bug fixes and documentation corrections.

---

## Testing

Run all tests:

```bash
go test ./...
```

Run parser tests:

```bash
go test ./internal/parser -v
```

Run formatting:

```bash
go fmt ./...
```

Run static package listing:

```bash
go list ./...
```

> [!NOTE]
> Some tests may require `logs/inputs` and a valid configuration file. Database-backed tests require PostgreSQL to be reachable using the configured connection string.

---

## License

LogSentry is released under the MIT License.

```text
MIT License

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files, to deal in the Software
without restriction, including without limitation the rights to use, copy,
modify, merge, publish, distribute, sublicense, and/or sell copies of the
Software, subject to inclusion of the copyright notice and permission notice
in all copies or substantial portions of the Software.
```

---

## Acknowledgements

LogSentry builds on excellent open-source projects:

- [Go](https://go.dev/) for systems programming and concurrency.
- [Gin](https://gin-gonic.com/) for HTTP routing.
- [PostgreSQL](https://www.postgresql.org/) for durable structured storage.
- [pgx](https://github.com/jackc/pgx) for PostgreSQL connectivity.
- [fsnotify](https://github.com/fsnotify/fsnotify) for file system notifications.
- [Prometheus](https://prometheus.io/) for metrics.
- [Grafana](https://grafana.com/) for dashboards.
- [Docker](https://www.docker.com/) for reproducible deployment.

---

## Contact

Project repository:

```text
https://github.com/amandx36/LogSentry
```

Maintainer profile:

```text
https://github.com/amandx36
```

-------amandx36
