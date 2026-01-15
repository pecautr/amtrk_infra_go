# Amtrak NEC Optimization System

A comprehensive train dispatching optimization system for Amtrak's Northeast Corridor, inspired by Deutsche Bahn's ADA-PMB platform. This system uses machine learning and mathematical optimization to minimize delays and conflicts in real-time.

## Features

✅ **Real-time GTFS-Realtime Data Ingestion** - Fetches live vehicle positions, trip updates, and service alerts  
✅ **Machine Learning Predictions** - Forecasts delays, predicts conflicts, recommends optimal dispatch decisions  
✅ **Mathematical Optimization** - Mixed Integer Programming to minimize total network delays  
✅ **Vehicle Performance Modeling** - Physics-based travel time calculations for different trainsets  
✅ **Hourly Execution Cycle** - Automated conflict detection and resolution recommendations  
✅ **Real-time Arrival Updates** - Publishes revised predictions to passenger information systems  

## Architecture

The system consists of several key components:

- **Data Ingestion Layer** ([pkg/ingestion](pkg/ingestion)) - GTFS-RT feed processing
- **Data Model** ([pkg/datamodel](pkg/datamodel)) - Network state representation
- **Vehicle Performance** ([pkg/vehicle](pkg/vehicle)) - Trainset characteristics and physics models
- **ML Engine** ([pkg/ml](pkg/ml)) - Delay prediction, conflict prediction, dispatch advisor
- **Optimization Engine** ([pkg/optimization](pkg/optimization)) - MIP solver for dispatch decisions
- **Orchestrator** ([cmd/orchestrator](cmd/orchestrator)) - Main execution loop

See [docs/architecture.md](docs/architecture.md) for detailed architecture diagrams.

## Quick Start

### Prerequisites

- Go 1.21 or later
- PostgreSQL 14+ (optional, for data persistence)

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/amtrk_infra_go.git
cd amtrk_infra_go

# Run initialization script
bash scripts/init.sh

# Set environment variables
export AMTRAK_API_KEY=your_api_key
export DB_PASSWORD=your_db_password

# Build the application
go build -o bin/orchestrator cmd/orchestrator/main.go

# Run
./bin/orchestrator
```

## Configuration

Edit [config/config.yaml](config/config.yaml) to customize:

- GTFS-RT feed URLs and refresh intervals
- Optimization parameters (weights, solve time, conflict horizon)
- ML model paths and retraining intervals
- Database connection settings
- Logging and monitoring options

## Documentation

- **[Architecture](docs/architecture.md)** - System design and component interactions
- **[Development Guide](docs/development.md)** - How to build, test, and extend the system
- **[ML Models](docs/ml_models.md)** - Machine learning model details and training procedures

## Project Structure

```
amtrk_infra_go/
├── cmd/orchestrator/       # Main application
├── pkg/
│   ├── datamodel/          # Core data structures
│   ├── ingestion/          # GTFS-RT data fetching
│   ├── optimization/       # MIP optimization engine
│   ├── ml/                 # Machine learning components
│   └── vehicle/            # Vehicle performance modeling
├── config/                 # Configuration files
├── docs/                   # Documentation
├── scripts/                # Utility scripts
└── models/                 # ML model files
```

## Key Technologies

- **Language**: Go 1.21
- **Optimization**: Mixed Integer Programming
- **ML Framework**: Python (scikit-learn, XGBoost, TensorFlow)
- **Database**: PostgreSQL with TimescaleDB
- **Data Format**: GTFS-Realtime (Protocol Buffers)

## System Workflow

Every hour, the orchestrator:

1. **Ingests** GTFS-RT feeds (vehicle positions, trip updates, alerts)
2. **Updates** network state (active trips, track occupancy, station status)
3. **Predicts** delays and conflicts using ML models
4. **Detects** potential conflicts in the next hour
5. **Optimizes** dispatch decisions to minimize total delay
6. **Publishes** recommendations to dispatchers
7. **Updates** real-time arrival predictions
8. **Records** decisions for ML training

## Performance Targets

| Metric | Target |
|--------|--------|
| Data ingestion latency | < 5s |
| ML prediction time | < 10s |
| Optimization solve time | < 30s |
| Total cycle time | < 60s |
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

TBD

## References

- [GTFS-Realtime Specification](https://gtfs.org/realtime/)
- [Amtrak Developer API](https://www.amtrak.com/developer)
- Deutsche Bahn ADA-PMB System (inspiration)

## Contact

For questions or support, please file an issue on GitHub.

---

**Built with ❤️ for better transit operations**
