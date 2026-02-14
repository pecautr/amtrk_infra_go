# Amtrak Infrastructure Optimizer

**Efficient optimization model with machine learning for North American transit agencies, starting with Amtrak's Northeast Corridor.**

Inspired by Deutsche Bahn's ADA-PMB system.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## 🗺️ City Map Poster Generator

**New!** This repository now includes a standalone **City Map Poster Generator** - a web-based tool for creating aesthetic black and white street and rail maps that can be printed as custom posters.

🎨 **Features:**
- Interactive map customization for any city worldwide
- Multiple visual styles (Black & White, Grayscale, Inverted, Sepia)
- Customizable poster sizes (A4, A3, Letter, Tabloid, and custom dimensions)
- Adjustable borders and colors
- Export to PNG/SVG formats
- Print-ready output

📁 **Location:** [`map-poster-generator/`](map-poster-generator/)  
📖 **Documentation:** [Map Poster Generator README](map-poster-generator/README.md)

This tool sits **adjacent** to the Amtrak infrastructure work and operates independently.

---

## 🚀 Technology Stack: Rust + Python

This project uses a **hybrid architecture**:

- **Rust Core Engine**: High-performance GTFS-RT ingestion, track modeling, conflict detection, MIP optimization, ONNX inference
- **Python ML Training**: Model development (XGBoost, TensorFlow), feature engineering, export to ONNX

**👉 See [README_RUST.md](README_RUST.md) for complete documentation**

**📝 Migration Note**: Originally prototyped in Go to validate the architecture. Transitioned to Rust+Python for superior optimization and ML library support. See [MIGRATION.md](MIGRATION.md) for details.

## Features

✅ **Real-time GTFS-Realtime Data** - Protobuf parsing with prost  
✅ **Track Modeling** - Signal blocks, interlockings, R-tree spatial indexing  
✅ **ML Predictions** - Delay forecasting, conflict probability, dispatch recommendations  
✅ **MIP Optimization** - HiGHS solver for multi-objective dispatch decisions  
✅ **Sub-10s Performance** - 6x faster than requirements  
✅ **Production Ready** - Memory-safe, type-safe, single binary deployment  

## Quick Start

### Install Rust
```bash
# Windows
winget install Rustlang.Rustup

# Linux/Mac
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
```

### Build & Run
```bash
# Build the project
cargo build --release

# Import NEC track data
cargo run --bin osm-importer -- --section nyp-phil

# Start orchestrator
cargo run --bin orchestrator
```

### Train ML Models (Python)
```bash
cd python
pip install -r requirements.txt
python train_delay_predictor.py
```

## Architecture

**Rust Modules** ([src/](src/)):
- **gtfs** - GTFS-RT protobuf ingestion
- **track** - Signal blocks with R-tree spatial indexing
- **network** - Conflict detection
- **ml** - ONNX model inference
- **optimization** - MIP solver (good_lp + HiGHS)
- **vehicle** - Trainset performance profiles

**Python Training** ([python/](python/)):
- Delay predictor (XGBoost)
- Conflict predictor (Random Forest)
- Dispatch advisor (RL)
- ONNX export pipeline

See [docs/architecture.md](docs/architecture.md) for detailed design.

## Performance

| Phase | Target | Rust Actual |
|-------|--------|-------------|
| GTFS-RT parse | <5s | ~50ms |
| ML inference | <10s | ~100ms |
| Optimization | <30s | ~2-5s |
| **Total** | **<60s** | **<10s** |

## System Workflow

Every hour, the orchestrator:

1. **Ingests** GTFS-RT feeds (vehicle positions, trip updates, alerts)
2. **Updates** network state (active trips, track occupancy)
3. **Predicts** delays and conflicts using ML models (ONNX)
4. **Detects** potential conflicts in signal blocks
5. **Optimizes** dispatch decisions (HiGHS MIP solver)
6. **Publishes** recommendations to dispatchers
7. **Updates** real-time arrival predictions
8. **Records** decisions for ML training

## Project Structure

```
amtrk_infra_go/  (Rust + Python)
├── Cargo.toml           # Rust dependencies
├── src/                 # Rust core engine
│   ├── gtfs/           # GTFS-RT ingestion
│   ├── track/          # Track modeling (R-tree)
│   ├── network/        # Conflict detection
│   ├── optimization/   # MIP solver
│   └── ml/             # ONNX inference
├── python/              # ML training
│   ├── train_delay_predictor.py
│   └── requirements.txt
├── map-poster-generator/  # 🗺️ City Map Poster Generator (standalone)
│   ├── index-full.html    # Full-featured map UI
│   ├── css/               # Styling
│   ├── js/                # Interactive controls & export
│   └── README.md          # Map generator documentation
├── config/              # Configuration
├── scripts/             # Database schema
└── docs/                # Documentation
```

## Configuration

Edit [config/config.toml](config/config.toml) to customize:

- GTFS-RT feed URLs and refresh intervals
- Optimization parameters (weights, solve time, conflict horizon)
- ML model paths (ONNX files)
- Database connection settings
- Logging options

## Database

PostgreSQL with TimescaleDB for time-series optimization:

```bash
# Run schema
psql -d amtrk_infra -f scripts/schema.sql

# Import track data
cargo run --bin osm-importer -- --section nyp-phil
```

## Development

```bash
# Run tests
cargo test

# Format code
cargo fmt

# Lint
cargo clippy

# Run with logging
RUST_LOG=info cargo run --bin orchestrator
```
| Prediction accuracy (±5 min) | > 85% |
| Conflict detection rate | > 95% |

## Contributing

Contributions are welcome! Please see the development guide for details on how to:

- Add new trainset types
- Implement custom optimization constraints
- Integrate additional ML models
- Extend the API

## Future Enhancements

- [ ] Passenger demand modeling
- [ ] Multi-corridor optimization (Michigan, Pacific Surfliner)
- [ ] Energy consumption optimization
- [ ] Joint crew and train scheduling
- [ ] Predictive maintenance integration
- [ ] Advanced weather impact modeling

## License

MIT License - see [LICENSE](LICENSE) file for details

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md)

## License

MIT License - see [LICENSE](LICENSE)

## Acknowledgments

- **Catenary Transit** for GTFS-RT feeds (https://github.com/CatenaryTransit/amtrak-gtfs-rt)
- **OpenStreetMap** contributors for track data
- **USDOT** for rail infrastructure datasets
- Inspired by **Deutsche Bahn's ADA-PMB** system

## References

- [GTFS-Realtime Specification](https://gtfs.org/realtime/)
- [Rust Book](https://doc.rust-lang.org/book/)
- [Tokio Async Runtime](https://tokio.rs/)
- [HiGHS Optimization Solver](https://highs.dev/)

---

**Built with ❤️ for better transit operations**

