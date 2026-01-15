# Migration from Go to Rust + Python

## Why the Switch?

After initial exploration with Go, we've transitioned to a **Rust + Python hybrid architecture** for better optimization and ML capabilities.

### Decision Rationale

| Aspect | Go (Original) | Rust + Python (New) | Winner |
|--------|--------------|---------------------|---------|
| **Optimization Libraries** | Limited (requires C bindings) | good_lp (native) + Python (PuLP, OR-Tools) | ✅ Rust/Python |
| **ML Ecosystem** | Via gRPC to Python | Native Python training → ONNX inference in Rust | ✅ Rust/Python |
| **Performance** | Fast (GC pauses) | Faster (no GC, zero-cost abstractions) | ✅ Rust |
| **Type Safety** | Good | Excellent (borrow checker catches more bugs) | ✅ Rust |
| **Protobuf** | Code generation | Native prost (better integration) | ✅ Rust |
| **Spatial Indexing** | Manual | rstar R-tree (mature library) | ✅ Rust |
| **Development Speed** | Fast | Medium (learning curve) | ⚖️ Tie |
| **Memory Safety** | Runtime panics possible | Compile-time guarantees | ✅ Rust |
| **Deployment** | Single binary | Single binary + Python for training | ⚖️ Tie |

**Conclusion**: Rust + Python provides **better library support** for optimization and ML while maintaining production-grade performance and safety.

## What We Kept

The following concepts and designs transfer directly:

### Data Model
- ✅ VehiclePosition, TripUpdate, ServiceAlert (GTFS-RT compatible)
- ✅ NetworkState, ActiveTrip, Conflict, Recommendation
- ✅ Signal blocks, interlockings, track topology

### Architecture
- ✅ Hourly optimization loop
- ✅ 8-step cycle: Ingest → Update → Predict → Detect → Optimize → Publish → Update Arrivals → Record
- ✅ PostgreSQL + TimescaleDB for persistence
- ✅ GTFS-RT as primary data source
- ✅ Multi-objective optimization (delay, passengers, OTP)

### Documentation
- ✅ All docs (architecture.md, track_modeling.md, data_sources.md, ml_models.md)
- ✅ OSM extraction strategy
- ✅ Database schema
- ✅ Control point definitions (BOS, NYP, PHL, WAS, etc.)

## What Changed

### 1. GTFS-RT Parsing

**Go (Old)**:
```go
// Manual protobuf handling
feed := &transit_realtime.FeedMessage{}
proto.Unmarshal(bytes, feed)
```

**Rust (New)**:
```rust
// prost - cleaner, faster
use prost::Message;
let feed: proto::FeedMessage = Message::decode(&bytes[..])?;
```

### 2. Spatial Indexing

**Go (Old)**:
```go
// Manual distance calculations in loops
for _, block := range blocks {
    dist := haversineDistance(lat, lon, block.Lat, block.Lon)
    if dist < maxDistance {
        nearby = append(nearby, block)
    }
}
```

**Rust (New)**:
```rust
// R-tree for O(log n) queries
track_network.find_blocks_near(lat, lon, max_distance_km)
// Uses rstar spatial index - much faster for large datasets
```

### 3. Optimization

**Go (Old)**:
```go
// Placeholder heuristic (would need Gurobi CGO bindings)
func solveMIP() []Recommendation {
    // Greedy heuristic
}
```

**Rust (New)**:
```rust
// Native Rust MIP with HiGHS solver
use good_lp::*;
let solution = vars.minimise(objective)
    .using(default_solver)
    .solve()?;
```

### 4. ML Integration

**Go (Old)**:
```go
// gRPC to Python service (network overhead)
mlClient.PredictDelay(ctx, &request)
```

**Rust (New)**:
```rust
// In-process ONNX inference (zero overhead)
let model = tract_onnx::onnx().model_for_path("model.onnx")?;
let result = model.run(tvec![input.into()])?;
```

**Python (Training)**:
```python
# Train in Python, export to ONNX
model = xgb.XGBRegressor(...)
model.fit(X_train, y_train)
onnx_model = to_onnx(model, X_train[:1])
```

## Migration Checklist

### Phase 1: Core Infrastructure ✅
- [x] Set up Rust project (Cargo.toml)
- [x] Core data models
- [x] Error handling
- [x] Configuration management
- [x] GTFS-RT module (prost)
- [x] Track modeling with R-tree
- [x] Conflict detection
- [x] Optimization (good_lp)
- [x] ML inference (tract-onnx)
- [x] Vehicle performance profiles

### Phase 2: Data & Integration 🚧
- [ ] Download GTFS-RT .proto file
- [ ] Implement OSM importer (Rust version)
- [ ] Database persistence (sqlx)
- [ ] GTFS-RT feed parsing (real data)
- [ ] Track network loading from database

### Phase 3: Python ML 🚧
- [ ] Python environment setup
- [ ] Historical data extraction
- [ ] Delay predictor training
- [ ] Conflict predictor training
- [ ] Dispatch advisor (RL) training
- [ ] ONNX export pipeline

### Phase 4: Orchestration 📅
- [ ] Main orchestrator binary
- [ ] Hourly execution loop
- [ ] API server (Axum)
- [ ] Health checks
- [ ] Metrics & monitoring

### Phase 5: Deployment 📅
- [ ] Docker image
- [ ] docker-compose.yml
- [ ] Kubernetes manifests
- [ ] CI/CD pipeline
- [ ] Production configuration

## File Mapping

| Go File | Rust Equivalent | Status |
|---------|-----------------|--------|
| `pkg/datamodel/*.go` | `src/datamodel.rs` | ✅ Ported |
| `pkg/vehicle/*.go` | `src/vehicle/` | ✅ Ported |
| `pkg/optimization/*.go` | `src/optimization/` | ✅ Improved |
| `pkg/ml/*.go` | `src/ml/` + `python/` | ✅ Improved |
| `pkg/ingestion/*.go` | `src/gtfs/` | ✅ Improved |
| `cmd/orchestrator/main.go` | `src/bin/orchestrator.rs` | ⏳ TODO |
| `scripts/schema.sql` | `scripts/schema.sql` | ✅ Kept |
| `config/config.yaml` | `config/config.toml` | ⏳ TODO |
| `cmd/osm_importer/main.go` | `src/bin/osm_importer.rs` | ⏳ TODO |

## Getting Started with Rust

If you're new to Rust coming from Go:

### Install Rust
```bash
# Linux/Mac
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh

# Windows
# Download from https://rustup.rs/
```

### Key Differences from Go

| Go | Rust | Notes |
|----|------|-------|
| `defer` | Drop trait | Automatic cleanup |
| Goroutines | Tokio tasks | Similar async model |
| Channels | `mpsc`, `oneshot` | Similar patterns |
| `interface{}` | Generics + traits | More type-safe |
| `error` return | `Result<T, E>` | Explicit error handling |
| GC | Ownership | Manual but safe |

### Learning Resources

- **The Rust Book**: https://doc.rust-lang.org/book/
- **Rust by Example**: https://doc.rust-lang.org/rust-by-example/
- **Tokio Tutorial**: https://tokio.rs/tokio/tutorial
- **Axum Examples**: https://github.com/tokio-rs/axum/tree/main/examples

## Next Steps

1. **Install Rust**: Follow instructions above
2. **Build the project**: `cargo build`
3. **Run tests**: `cargo test`
4. **Review Rust code**: Start with `src/datamodel.rs` (familiar types)
5. **Read README_RUST.md**: Full architecture guide
6. **Set up Python**: `cd python; pip install -r requirements.txt`
7. **Train a model**: `python train_delay_predictor.py`

## Questions?

See detailed documentation:
- [README_RUST.md](README_RUST.md) - Full Rust architecture
- [docs/architecture.md](docs/architecture.md) - System design
- [docs/development.md](docs/development.md) - Build & deploy

## Acknowledgment

The Go implementation served as excellent prototyping and validated the architecture. Rust + Python provides the production-grade foundation we need for optimization and ML at scale.
