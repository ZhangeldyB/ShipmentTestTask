# Shipment gRPC Microservice

A Go gRPC microservice for tracking shipment lifecycle in a logistics system, backed by MongoDB.

Supports multiple transport modes (truck, air, sea, rail) through a flexible domain model. Designed for Clean Architecture + DDD with full separation between the domain, application, and infrastructure layers.

---

## How to Run

### With Docker Compose (recommended)

No external database needed — MongoDB runs as a sidecar container.

```bash
docker-compose up --build
# or via Makefile:
make up
```

The service starts on port `50051`. MongoDB data is persisted in a named Docker volume.

### With an external MongoDB (e.g. Atlas)

```bash
export DATABASE_URL="mongodb+srv://user:pass@cluster.mongodb.net"
export GRPC_PORT=50051          # optional, default 50051
go run ./cmd/server
# or: make build && ./bin/server
```

---

## How to Run Tests

Unit tests cover the domain and application layers. No database or running service is required — all I/O is replaced by in-memory mocks.

```bash
make test
# or: go test ./...

# With race detector (recommended before merging):
make test-race
# or: go test -race ./...
```

---

## Architecture Overview

The project follows **Clean Architecture** with a strict inward dependency rule:

```
cmd/server/
└── main.go              ← wires everything together

internal/
├── domain/              ← pure business logic, zero external deps
│   ├── shipment.go      ← Shipment aggregate root
│   ├── event.go         ← ShipmentEvent (append-only history)
│   ├── status.go        ← Status type + state machine transitions
│   ├── transport.go     ← TransportMode value object (TRUCK/AIR/SEA/RAIL)
│   ├── carrier.go       ← CarrierInfo value object (mode-agnostic operator data)
│   └── errors.go        ← typed domain errors
│
├── app/                 ← use cases, depends only on domain + repo interface
│   ├── repository.go    ← ShipmentRepository interface (defined here, not in domain)
│   ├── create_shipment.go
│   ├── get_shipment.go
│   ├── add_status_event.go
│   └── get_events.go
│
└── infra/
    ├── grpc/            ← gRPC transport adapter
    │   ├── handler.go   ← maps proto requests → use case calls
    │   ├── mapper.go    ← domain ↔ proto conversion
    │   └── server.go    ← grpc.Server with keepalive / concurrency tuning
    └── mongo/           ← MongoDB persistence adapter
        ├── dto.go       ← bson-tagged persistence structs (isolated from domain)
        ├── mapper.go    ← domain ↔ mongo doc conversion
        └── repository.go← implements app.ShipmentRepository
```

**Dependency direction:** `infra → app → domain`

The domain layer has no struct tags, no framework imports, and no knowledge of how or where it is stored or served.

---

## gRPC API

gRPC reflection is enabled — use `grpcurl` or Postman to explore without a generated client.

```bash
grpcurl -plaintext localhost:50051 list
```

| RPC | Description |
|-----|-------------|
| `CreateShipment` | Create a shipment in PENDING status |
| `GetShipment` | Fetch shipment by ID |
| `AddStatusEvent` | Advance the shipment to a new status |
| `GetShipmentEvents` | List the full status history for a shipment |

### Status machine

```
PENDING → ASSIGNED → PICKED_UP → IN_TRANSIT → DELIVERED
        ↘          ↘                        ↘
        CANCELLED   CANCELLED               FAILED
```

Terminal states (DELIVERED, FAILED, CANCELLED) are final — no further transitions allowed.

### Example session (grpcurl)

**Create a truck shipment:**
```bash
grpcurl -plaintext -d '{
  "origin": "Almaty",
  "destination": "Astana",
  "amount": 1000,
  "carrier_revenue": 700,
  "transport_mode": "TRANSPORT_MODE_TRUCK",
  "operator_name": "Aibek Dzhaksybekov",
  "operator_phone": "+77001234567",
  "unit_identifier": "KZ-001-AA"
}' localhost:50051 shipment.v1.ShipmentService/CreateShipment
```

**Create an air shipment:**
```bash
grpcurl -plaintext -d '{
  "origin": "Almaty",
  "destination": "Dubai",
  "amount": 5000,
  "carrier_revenue": 3500,
  "transport_mode": "TRANSPORT_MODE_AIR",
  "operator_name": "Pilot Seitkali",
  "operator_phone": "+77009876543",
  "unit_identifier": "KC-401"
}' localhost:50051 shipment.v1.ShipmentService/CreateShipment
```

**Advance to ASSIGNED:**
```bash
grpcurl -plaintext -d '{
  "shipment_id": "<id from above>",
  "new_status": "ASSIGNED",
  "note": "Driver confirmed"
}' localhost:50051 shipment.v1.ShipmentService/AddStatusEvent
```

**Get full event history:**
```bash
grpcurl -plaintext -d '{"shipment_id": "<id>"}' \
  localhost:50051 shipment.v1.ShipmentService/GetShipmentEvents
```

---

## Design Decisions

### 1. Transport-mode-agnostic carrier model
Instead of driver-specific fields (`driver_name`, `unit_number`), the domain uses a `CarrierInfo` value object with generic fields (`OperatorName`, `OperatorPhone`, `UnitIdentifier`). Adding a new transport mode (e.g. drone, pipeline) requires only a new constant in `transport.go` — no other domain code changes.

### 2. Domain free of infrastructure concerns
The `Shipment` and `ShipmentEvent` structs have no `bson`, `json`, or `db` tags. MongoDB DTOs live in `infra/mongo/dto.go` and are mapped to/from domain objects by `infra/mongo/mapper.go`. This means the persistence layer can be swapped (Postgres, Redis, etc.) without touching the domain.

### 3. Repository interface in the app layer
`ShipmentRepository` is defined in `internal/app/`, not in `internal/domain/`. The domain has no concept of persistence at all. The app layer knows it needs to persist things, but not how.

### 4. Append-only event log
Every status change writes a `ShipmentEvent`. The `Shipment` document holds the current state for fast lookups; the events collection holds the full audit trail. These are two separate writes — not a transaction — which is acceptable because the shipment document is the source of truth and events are supplementary.

### 5. gRPC tuning for high load
The server is configured with `MaxConcurrentStreams(1000)` and keepalive parameters so idle connections stay warm and the server recycles long-lived connections gracefully instead of holding them open forever.

### 6. MongoDB connection pooling
The client is configured with `MaxPoolSize=100` / `MinPoolSize=10` so the pool never drops to zero between traffic bursts, and never grows unbounded under spike load.

---

## Assumptions

- **Carrier management is handled by another service.** This service does not store full carrier profiles — it only records the operator name, phone, and unit identifier as provided at the time of shipment creation. A separate carrier/driver service would be responsible for carrier data management.

- **Multiple transport modes exist.** Logistics systems commonly mix truck, air, sea, and rail legs. The domain was modeled to be mode-agnostic from the start rather than truck-centric.

- **No authentication/authorization.** This is a backend internal service. Auth (e.g. mTLS between services) is assumed to be handled at the infrastructure level (API gateway, service mesh).

- **Atlas or any MongoDB-compatible URI is acceptable.** The service only uses standard MongoDB driver calls and works with any MongoDB 5+ instance. The docker-compose setup provides a local Mongo 7 container for development.

- **Single database `shipment_service`** with two collections: `shipments` and `shipment_events`. Indexes are created automatically on startup via `EnsureIndexes`.

- **No saga/distributed transaction.** The service writes the shipment update and the event as two separate operations. In a production system with strict consistency requirements, an outbox pattern or change streams could be added.
