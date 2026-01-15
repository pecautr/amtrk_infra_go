# Development Guide

## Getting Started

### Prerequisites

- Go 1.21 or later
- PostgreSQL 14+ (optional, for persistence)
- Git

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/amtrk_infra_go.git
cd amtrk_infra_go

# Download dependencies
go mod download

# Verify installation
go test ./...
```

### Project Structure

```
amtrk_infra_go/
├── cmd/
│   └── orchestrator/        # Main application entry point
│       └── main.go
├── pkg/
│   ├── datamodel/           # Core data structures
│   │   ├── vehicle_position.go
│   │   ├── network.go
│   │   └── historical.go
│   ├── ingestion/           # GTFS-RT data ingestion
│   │   └── gtfs_rt.go
│   ├── optimization/        # Optimization engine
│   │   ├── optimizer.go
│   │   └── mip_solver.go
│   ├── ml/                  # Machine learning components
│   │   ├── engine.go
│   │   └── training.go
│   └── vehicle/             # Vehicle performance modeling
│       ├── performance.go
│       └── registry.go
├── config/                  # Configuration files
│   └── config.yaml
├── data/                    # Data files (gitignored)
│   ├── raw/
│   └── processed/
├── models/                  # ML models (gitignored except .gitkeep)
├── docs/                    # Documentation
│   └── architecture.md
├── scripts/                 # Utility scripts
├── go.mod
├── go.sum
└── README.md
```

## Building the Application

### Development Build

```bash
go build -o bin/orchestrator cmd/orchestrator/main.go
```

### Run

```bash
# With default config
./bin/orchestrator

# With custom config
AMTRAK_API_KEY=your_key ./bin/orchestrator
```

### Production Build

```bash
# Build with optimizations
go build -ldflags="-s -w" -o bin/orchestrator cmd/orchestrator/main.go

# Build for Linux (from any OS)
GOOS=linux GOARCH=amd64 go build -o bin/orchestrator-linux cmd/orchestrator/main.go
```

## Testing

### Run All Tests

```bash
go test ./...
```

### Run with Coverage

```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Run Specific Package

```bash
go test ./pkg/optimization/...
```

### Benchmark Tests

```bash
go test -bench=. ./pkg/optimization/
```

## Configuration

The system uses a YAML configuration file (`config/config.yaml`) with the following sections:

### GTFS-RT Feeds

```yaml
gtfs_rt:
  feeds:
    vehicle_positions: "https://api-v3.amtrak.com/gtfsrt/v1/vehiclepositions"
    trip_updates: "https://api-v3.amtrak.com/gtfsrt/v1/tripupdates"
    service_alerts: "https://api-v3.amtrak.com/gtfsrt/v1/alerts"
  refresh_interval: 30s
  api_key: "" # Set via AMTRAK_API_KEY env var
```

### Optimization Parameters

```yaml
optimization:
  max_solve_time: 30s
  conflict_horizon: 2h
  minimum_headway: 180
  objective_weights:
    total_delay: 0.4
    max_delay: 0.3
    passenger_impact: 0.2
```

### Environment Variables

```bash
# Amtrak API credentials
export AMTRAK_API_KEY=your_api_key

# Database connection
export DB_PASSWORD=your_db_password

# Logging
export LOG_LEVEL=debug
```

## Development Workflow

### 1. Create a Feature Branch

```bash
git checkout -b feature/new-conflict-detector
```

### 2. Make Changes

Follow Go conventions:
- Package names: lowercase, single word
- Exported names: Start with capital letter
- Unexported: Start with lowercase letter
- Error handling: Always check and return errors

### 3. Write Tests

```go
func TestConflictDetector(t *testing.T) {
    detector := NewConflictDetector()
    
    // Test scenario
    state := createTestNetworkState()
    conflicts := detector.Detect(state)
    
    if len(conflicts) != 2 {
        t.Errorf("Expected 2 conflicts, got %d", len(conflicts))
    }
}
```

### 4. Format Code

```bash
# Format all files
go fmt ./...

# Run linter
golangci-lint run
```

### 5. Commit Changes

```bash
git add .
git commit -m "feat: add advanced conflict detection for platform conflicts"
```

## Debugging

### Enable Debug Logging

```yaml
logging:
  level: "debug"
```

### Inspect Network State

Add logging in orchestrator:

```go
log.WithFields(logrus.Fields{
    "active_trips": len(o.networkState.ActiveTrips),
    "conflicts": len(o.networkState.Conflicts),
}).Debug("Network state snapshot")
```

### Profile Performance

```bash
# CPU profiling
go run cmd/orchestrator/main.go -cpuprofile=cpu.prof

# Memory profiling
go run cmd/orchestrator/main.go -memprofile=mem.prof

# Analyze
go tool pprof cpu.prof
```

## Common Tasks

### Add a New Trainset Type

1. Edit `pkg/vehicle/registry.go`
2. Add to `GetDefaultNECTrainsets()`:

```go
{
    TypeID:       "avelia-liberty",
    Name:         "Avelia Liberty",
    Manufacturer: "Alstom",
    MaxSpeed:     320.0,
    ServiceSpeed: 260.0,
    Length:       200.0,
    Capacity:     386,
    PowerType:    "Electric",
}
```

3. Add performance profile to `GetDefaultPerformanceProfiles()`

### Add a New Optimization Constraint

1. Define constraint type in `pkg/optimization/optimizer.go`:

```go
const (
    // ... existing
    PassengerComfort ConstraintType = iota
)
```

2. Implement in `pkg/optimization/mip_solver.go`:

```go
func (s *MIPSolver) addPassengerComfortConstraint(model *MIPModel) {
    // Constraint logic
}
```

### Integrate a New ML Model

1. Define model in `pkg/ml/engine.go`
2. Add prediction method
3. Update `Predict()` to call new model
4. Add training logic in `pkg/ml/training.go`

## API Development

### Adding a New Endpoint

Create `pkg/api/server.go`:

```go
package api

import (
    "github.com/go-chi/chi/v5"
    "net/http"
)

func (s *Server) RegisterRoutes(r chi.Router) {
    r.Get("/api/v1/state", s.handleGetState)
}

func (s *Server) handleGetState(w http.ResponseWriter, r *http.Request) {
    // Handler logic
}
```

## Performance Optimization Tips

1. **Parallel Processing**: Use goroutines for independent operations
2. **Caching**: Cache frequently accessed data (trainset profiles, station configs)
3. **Connection Pooling**: Reuse HTTP clients and database connections
4. **Batch Updates**: Accumulate and batch database writes
5. **Profiling**: Use pprof to identify bottlenecks

## Troubleshooting

### Issue: Optimization times out

**Solution**: Reduce `conflict_horizon` or `max_solve_time` in config

### Issue: ML predictions low quality

**Solution**: Increase `min_training_data` and retrain models

### Issue: GTFS-RT feed errors

**Solution**: Check API key and network connectivity

## Resources

- [Go Documentation](https://golang.org/doc/)
- [GTFS-Realtime Specification](https://gtfs.org/realtime/)
- [Amtrak API](https://www.amtrak.com/developer)
- [Mixed Integer Programming](https://en.wikipedia.org/wiki/Integer_programming)

## Getting Help

- File issues on GitHub
- Check `docs/architecture.md` for system design
- Review test files for usage examples
