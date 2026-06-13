# API Gateway with Event Ingestion Pipeline Gateway

A production-style API gateway and event ingestion pipeline for security telemetry, built in Go.

## Overview

 API Gateway with Event Ingestion Pipeline is the entry point for a security event ingestion pipeline. It accepts events from distributed security agents (e.g. host-based monitors reporting authentication failures, network anomalies, file integrity changes), authenticates and rate-limits clients, validates event structure against a defined schema, and forwards accepted events to a NATS JetStream queue for asynchronous processing.

The gateway is intentionally minimal by design, it does not parse logs deeply, run detections, store events, or generate alerts. Those responsibilities belong to downstream services (Ingestion Worker, Event Storage), keeping the gateway fast, stateless, and easy to scale horizontally.

## Architecture

```
Security Agent
     |
     v
+-------------+
| API Gateway |
+-------------+
     |
     v
+------------------+
| Event Queue      |
| (NATS JetStream) |
+------------------+
     |
     v
+------------------+
| Ingestion Worker |
+------------------+
     |
     v
+------------------+
| Event Storage    |
+------------------+
```

## Status

| Component         | Status |
|-------------------|--------|
| API Gateway       | Done |
| Event Queue (NATS)| Done |
| Ingestion Worker  | Done |
| Event Storage     | Planned |

## What the Gateway Does

- Accepts events via `POST /api/v1/events`
- Authenticates clients via API key (`Authorization: Bearer <key>`)
- Validates event structure and required fields
- Rate limits per client IP
- Forwards accepted events to a NATS JetStream stream

## What the Gateway Deliberately Does Not Do

- Deep log parsing
- Threat detection or correlation
- Permanent event storage
- Alert generation

These responsibilities are handled by downstream services, keeping the gateway's surface area small and its failure modes predictable.

## Event Format

All events must conform to the following schema:

```json
{
  "source": "linux-agent",
  "host": "web-01",
  "event_type": "authentication",
  "severity": "medium",
  "timestamp": "2026-06-04T15:20:00Z",
  "message": "Failed SSH login",
  "metadata": {
    "username": "admin",
    "ip": "10.0.0.15"
  }
}
```

### Validation Rules

| Field        | Requirement |
|--------------|-------------|
| `source`     | Required |
| `host`       | Required |
| `event_type` | Required |
| `timestamp`  | Required, RFC3339 format |
| `severity`   | Required, one of: `low`, `medium`, `high`, `critical` |

Invalid or missing fields are rejected with `400 Bad Request`.

## Tech Stack

- **Language:** Go
- **HTTP Framework:** [Gin](https://github.com/gin-gonic/gin)
- **Message Queue:** [NATS JetStream](https://docs.nats.io/nats-concepts/jetstream)
- **Auth:** API keys (`Authorization: Bearer <key>`)
- **Containerization:** Docker / Docker Compose

## Running Locally

### Prerequisites

- Go 1.26+
- Docker & Docker Compose

### Start NATS

```bash
docker compose up -d nats
```

### Run the Gateway

```bash
go run ./cmd
```

The gateway listens on `:8080`.

### Send a Test Event

```bash
curl -X POST http://localhost:8080/api/v1/events \
  -H "Authorization: Bearer sk_live_abc123" \
  -H "Content-Type: application/json" \
  -d '{
    "source": "linux-agent",
    "host": "web-01",
    "event_type": "authentication",
    "severity": "medium",
    "timestamp": "2026-06-04T15:20:00Z",
    "message": "Failed SSH login",
    "metadata": {
      "username": "admin",
      "ip": "10.0.0.15"
    }
  }'
```

Expected response:

```json
{ "status": "accepted" }
```

### Verify the Message Reached NATS

```bash
curl http://localhost:8222/jsz?streams=true
```

Look for `"messages": 1` (or however many requests sent) under the `EVENTS` stream.

## Request/Response Reference

| Scenario | Status Code |
|----------|-------------|
| Valid event, valid API key | `202 Accepted` |
| Missing/invalid required field | `400 Bad Request` |
| Invalid severity value | `400 Bad Request` |
| Missing or malformed `Authorization` header | `401 Unauthorized` |
| Invalid API key | `401 Unauthorized` |
| Rate limit exceeded | `429 Too Many Requests` |
| Failed to publish to queue | `500 Internal Server Error` |

## Project Structure

```
.
├── cmd/
│   └── main.go              # Application entrypoint
├── internal/
│   ├── domain/               # Core data models (Event)
│   ├── handler/               # HTTP handlers
│   ├── middleware/            # Auth & rate limiting
│   ├── queue/                  # NATS JetStream publisher
│   └── validator/              # Event validation logic
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── go.sum
```

## Roadmap


- [ ] Event Storage schema & migrations
- [ ] `/health` endpoint
- [ ] Persistent NATS storage volume
- [ ] Move API keys to PostgreSQL-backed store
- [ ] Structured logging
- [ ] Unit & integration tests