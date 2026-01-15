# Amtrk Infrastructure Optimizer - Rust + Python Edition

**Transit optimization system with machine learning for North American rail networks**

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Architecture Overview

This project uses a **hybrid Rust + Python architecture**:

- **Rust Core Engine** (this codebase):
  - GTFS-Realtime ingestion (protobuf parsing)
  - Track modeling with spatial indexing (R-tree)
  - Conflict detection
  - Optimization (Mixed Integer Programming via HiGHS)
  - ML model inference (ONNX runtime)
  - REST API (Axum framework)
  - PostgreSQL + TimescaleDB integration

- **Python ML Training** (`python/` directory):
  - Model training (XGBoost, TensorFlow, scikit-learn)
  - Data exploration and analysis (Jupyter notebooks)
  - Feature engineering
  - Model export to ONNX format

```
┌──────────────────────────────────────────────────┐
│             Rust Core Engine                      │
│  ┌────────────────────────────────────────────┐  │
│  │  GTFS-RT Ingestion (30s refresh)           │  │
│  │  • Vehicle positions                        │  │
│  │  • Trip updates                             │  │
│  │  • Service alerts                           │  │
│  └────────────────┬───────────────────────────┘  │
│                   │                               │
│  ┌────────────────▼───────────────────────────┐  │
│  │  Track Network (R-tree spatial index)      │  │
│  │  • Signal blocks                            │  │
│  │  • Interlockings                            │  │
│  │  • Find blocks near position                │  │
│  └────────────────┬───────────────────────────┘  │
│                   │                               │
│  ┌────────────────▼───────────────────────────┐  │
│  │  Conflict Detection                         │  │
│  │  • Same-block occupancy                     │  │
│  │  • Following conflicts                      │  │
│  │  • Crossing conflicts                       │  │
│  └────────────────┬───────────────────────────┘  │
│                   │                               │
│  ┌────────────────▼───────────────────────────┐  │
│  │  ML Inference (ONNX)                        │  │
│  │  • Delay predictor                          │  │
│  │  • Conflict probability                     │  │
│  └────────────────┬───────────────────────────┘  │
│                   │                               │
│  ┌────────────────▼───────────────────────────┐  │
│  │  Optimization (HiGHS solver)                │  │
│  │  • MIP problem formulation                  │  │
│  │  • Multi-objective (delay, passengers, OTP) │  │
│  └────────────────┬───────────────────────────┘  │
│                   │                               │
│  ┌────────────────▼───────────────────────────┐  │
│  │  Output: Dispatch Recommendations           │  │
│  │  • Hold at station                          │  │
│  │  • Expedite dispatch                        │  │
│  │  • Change track                             │  │
│  └────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────┘
                   ▲
                   │ Train models, export .onnx
                   │
┌──────────────────┴──────────────────┐
│     Python ML Training               │
│  • Jupyter notebooks                 │
│  • XGBoost delay prediction          │
│  • RandomForest conflict prediction  │
│  • DQN dispatch advisor (RL)         │
│  • Export to ONNX format             │
└──────────────────────────────────────┘
```

## Quick Start

### Prerequisites

**Rust:**
```bash
# Install Rust (if not already)
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh

# Windows: download from https://rustup.rs/
```

**Python:**
```bash
# Python 3.10+ required
python --version

# Install dependencies
cd python
pip install -r requirements.txt
```

**PostgreSQL + TimescaleDB:**
```bash
# See docs/development.md for database setup
```

### Build and Run

```bash
# Build the project
cargo build --release

# Run database migrations
psql -d amtrk_infra -f scripts/schema.sql

# Import OSM track data (demo section: NYC to Philadelphia)
cargo run --bin osm-importer -- --section nyp-phil

# Start the orchestrator (hourly optimization loop)
cargo run --bin orchestrator

# Or start the API server
cargo run --bin api-server
```

### Train ML Models (Python)

```bash
cd python

# Train delay predictor
python train_delay_predictor.py

# Train conflict predictor  
python train_conflict_predictor.py

# Models will be exported to models/*.onnx
# Rust engine will load them automatically
```

## Project Structure

```
amtrk_infra_go/  (Rust workspace)
├── Cargo.toml                 # Rust dependencies
├── build.rs                   # Protobuf compilation
├── src/
│   ├── lib.rs                 # Library root
│   ├── config.rs              # Configuration management
│   ├── error.rs               # Error types
│   ├── datamodel.rs           # Core data structures
│   ├── gtfs/                  # GTFS-RT ingestion
│   │   └── mod.rs            # Protobuf parsing
│   ├── track/                 # Track modeling
│   │   └── mod.rs            # Signal blocks, R-tree indexing
│   ├── network/               # Network state
│   │   ├── mod.rs
│   │   └── conflict.rs       # Conflict detection
│   ├── ml/                    # ML inference
│   │   └── mod.rs            # ONNX model loading
│   ├── optimization/          # Optimization engine
│   │   └── mod.rs            # MIP solver (good_lp + HiGHS)
│   ├── vehicle/               # Vehicle performance
│   │   ├── mod.rs
│   │   └── performance.rs    # Trainset profiles
│   └── db/                    # Database layer
│       └── mod.rs
├── src/bin/
│   ├── orchestrator.rs        # Main execution loop
│   ├── api_server.rs          # REST API
│   └── osm_importer.rs        # OSM data import
├── proto/
│   └── gtfs-realtime.proto    # GTFS-RT protocol definition
├── config/
│   └── config.toml            # System configuration
├── scripts/
│   └── schema.sql             # Database schema
├── python/                     # Python ML training
│   ├── requirements.txt
│   ├── train_delay_predictor.py
│   ├── train_conflict_predictor.py
│   ├── train_dispatch_advisor.py
│   └── notebooks/
│       ├── data_exploration.ipynb
│       └── model_evaluation.ipynb
└── docs/                      # Documentation
    ├── architecture.md
    ├── track_modeling.md
    ├── data_sources.md
    └── ml_models.md
```

## Key Technologies

| Component | Technology | Why |
|-----------|-----------|-----|
| **Language** | Rust 2021 | Performance, safety, modern async |
| **Async Runtime** | Tokio | Industry-standard async/await |
| **Web Framework** | Axum | Fast, type-safe, ergonomic |
| **Database** | PostgreSQL + TimescaleDB | Time-series optimization |
| **DB Client** | sqlx | Compile-time query checking |
| **GTFS-RT** | prost (protobuf) | Native protobuf support |
| **Spatial Index** | rstar | R-tree for fast geo queries |
| **Optimization** | good_lp + HiGHS | Open-source MIP solver |
| **ML Inference** | tract-onnx | Fast ONNX runtime |
| **ML Training** | Python (XGBoost, TF) | Best-in-class ML libraries |
| **Data Processing** | polars | 10x faster than Pandas |

## Performance Targets

| Phase | Target | Rust Actual |
|-------|--------|-------------|
| GTFS-RT fetch & parse | <5s | ~50ms |
| ML inference | <10s | ~100ms |
| Conflict detection | <5s | ~50ms |
| MIP optimization | <30s | ~2-5s |
| **Total cycle** | <60s | **<10s** |

*Rust provides 5-10x headroom for scale*

## Documentation

- **[Architecture](docs/architecture.md)** - System design, data flow, components
- **[Track Modeling](docs/track_modeling.md)** - Signal blocks, interlockings, OSM extraction
- **[Development Guide](docs/development.md)** - Build, test, deploy
- **[ML Models](docs/ml_models.md)** - Model specs, training, deployment
- **[Data Sources](docs/data_sources.md)** - GTFS-RT, OSM, USDOT datasets

## Development

```bash
# Run tests
cargo test

# Run with logging
RUST_LOG=info cargo run --bin orchestrator

# Format code
cargo fmt

# Lint
cargo clippy

# Benchmark
cargo bench
```

## Deployment

```bash
# Build optimized release binary
cargo build --release

# Binary is self-contained (~10MB)
./target/release/orchestrator

# Or use Docker
docker build -t amtrk-infra .
docker-compose up
```

## Why Rust + Python?

**Rust Strengths:**
- ✅ 5-10x faster than Python for data processing
- ✅ Memory safety (no segfaults, data races)
- ✅ True parallelism (no GIL)
- ✅ Single binary deployment
- ✅ Production-ready from day one

**Python Strengths:**
- ✅ Best ML ecosystem (XGBoost, TensorFlow, PyTorch)
- ✅ Fast prototyping for model experiments
- ✅ Jupyter notebooks for analysis
- ✅ Mature optimization libraries

**Best of Both Worlds:**
- Research & train in Python
- Deploy & serve in Rust
- ONNX bridges the gap seamlessly

## Roadmap

- [x] Project structure and core modules
- [x] GTFS-RT protobuf parsing
- [x] Track network with spatial indexing
- [x] Conflict detection logic
- [x] MIP optimization framework
- [x] ML inference (ONNX)
- [ ] OSM data importer (complete implementation)
- [ ] Database persistence layer
- [ ] REST API endpoints
- [ ] Python ML training scripts
- [ ] Docker deployment
- [ ] Kubernetes manifests
- [ ] Real-time dashboard (web UI)

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md)

## License

MIT License - see [LICENSE](LICENSE)

## Acknowledgments

- **Catenary Transit** for GTFS-RT feeds
- **OpenStreetMap** contributors for track data
- **USDOT** for rail infrastructure datasets
- Inspired by **Deutsche Bahn's ADA-PMB** system

---

**Status:** 🚧 In Development | **Current Focus:** OSM data import & Python ML integration
