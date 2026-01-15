# Amtrak NEC Optimization System - Architecture

## System Overview

This document describes the architecture of the Amtrak Northeast Corridor (NEC) optimization system, inspired by Deutsche Bahn's ADA-PMB platform.

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                     GTFS-Realtime Feeds                         │
│  (Vehicle Positions | Trip Updates | Service Alerts)            │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Ingestion Layer                              │
│  - GTFS-RT Parser                                               │
│  - Data Validation                                              │
│  - 30-second refresh cycle                                      │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│                  Network State Manager                          │
│  - Active trips tracking                                        │
│  - Track occupancy monitoring                                   │
│  - Station status management                                    │
│  - Real-time state synchronization                              │
└─────────┬──────────────────────────────┬────────────────────────┘
          │                              │
          ▼                              ▼
┌─────────────────────┐        ┌─────────────────────────────────┐
│   ML Engine         │        │  Conflict Detector              │
│                     │        │                                 │
│ - Delay Predictor   │        │ - Track conflicts               │
│ - Conflict Predictor│        │ - Platform conflicts            │
│ - Dispatch Advisor  │        │ - Headway violations            │
│ - Performance       │        │ - Timing conflicts              │
│   Analyzer          │        │                                 │
└─────────┬───────────┘        └────────────┬────────────────────┘
          │                                 │
          │        ┌────────────────────────┘
          │        │
          ▼        ▼
┌─────────────────────────────────────────────────────────────────┐
│              Optimization Engine (MIP Solver)                   │
│                                                                 │
│  Objective: Minimize total delay + passenger impact            │
│  Constraints:                                                   │
│    - Track capacity                                             │
│    - Platform availability                                      │
│    - Minimum headway (3 min)                                    │
│    - Crew schedules                                             │
│    - Speed restrictions                                         │
│                                                                 │
│  Output: Dispatch recommendations                               │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│                  Recommendation Publisher                       │
│                                                                 │
│  - Dispatcher dashboard                                         │
│  - Train operator displays                                      │
│  - Passenger information systems                                │
│  - API endpoints                                                │
└─────────────────────────────────────────────────────────────────┘
```

## Component Details

### 1. Data Ingestion (`pkg/ingestion`)
- **Purpose**: Fetch and parse GTFS-Realtime feeds from Amtrak API
- **Frequency**: Every 30 seconds
- **Outputs**: Vehicle positions, trip updates, service alerts

### 2. Data Model (`pkg/datamodel`)
- **NetworkState**: Complete snapshot of NEC operations
- **VehiclePosition**: Real-time train locations and speeds
- **Trip**: Journey details with scheduled vs. actual progress
- **Conflict**: Detected or predicted conflicts
- **Recommendation**: Dispatching instructions

### 3. Vehicle Performance (`pkg/vehicle`)
- **TrainsetRegistry**: Database of train types (Acela, ACS-64, etc.)
- **PerformanceProfile**: Acceleration, braking, power curves
- **RunningTimeCalculator**: Physics-based travel time estimation
- **Metrics Tracking**: Real-time and historical performance

### 4. Optimization Engine (`pkg/optimization`)
- **Solver**: Mixed Integer Programming (MIP)
- **Objective Function**:
  ```
  Minimize: w1·TotalDelay + w2·MaxDelay + w3·PassengerImpact
            + w4·EnergyUse - w5·NetworkThroughput
  ```
- **Decision Variables**:
  - Hold durations at stations
  - Dispatch times
  - Track assignments
- **Constraints**:
  - Track capacity limits
  - Platform availability windows
  - Minimum headway between trains
  - Crew working hours
  - Speed restrictions

### 5. Machine Learning Engine (`pkg/ml`)

#### 5.1 Delay Predictor
- **Input**: Current delay, time of day, weather, network load
- **Output**: Predicted delay in 1 hour
- **Model**: Gradient Boosting Regressor
- **Training**: Historical delay patterns

#### 5.2 Conflict Predictor
- **Input**: Network state, scheduled movements
- **Output**: Conflict probability and severity
- **Model**: Random Forest Classifier
- **Training**: Historical conflict occurrences

#### 5.3 Dispatch Advisor
- **Input**: Network conditions, conflict predictions
- **Output**: Optimal hold times
- **Model**: Reinforcement Learning (DQN)
- **Training**: Outcomes of past dispatch decisions

#### 5.4 Performance Analyzer
- **Purpose**: Learn from historical decisions
- **Metrics**: Success rate, delay reduction, passenger impact
- **Feedback Loop**: Continuously improves model parameters

### 6. Orchestrator (`cmd/orchestrator`)
Main execution loop running every hour:

```
1. Ingest GTFS-RT data (5s)
2. Update network state (5s)
3. Run ML predictions (10s)
4. Detect conflicts (5s)
5. Run optimization (30s)
6. Publish recommendations (3s)
7. Update arrival times (2s)
8. Record decisions (1s)
Total: ~60s
```

## Data Flow Example

### Scenario: Train 2150 (Acela) approaching New York Penn Station

1. **T=0s**: GTFS-RT reports Train 2150 at Newark, 5 min delayed
2. **T=5s**: Network state updated with position and delay
3. **T=10s**: ML predicts:
   - Delay will grow to 8 min due to track congestion
   - 75% probability of conflict with Train 186 at Platform 11
4. **T=15s**: Conflict detector confirms platform conflict
5. **T=20s**: Optimization calculates:
   - Option A: Hold 2150 at Newark for 3 min → Total delay: 480 passenger-min
   - Option B: Reroute 186 to Platform 13 → Total delay: 320 passenger-min
   - **Select Option B**
6. **T=50s**: Publish recommendation: "Assign Train 186 to Platform 13"
7. **T=53s**: Update arrival displays: Train 2150 ETA +5 min, Train 186 on-time
8. **T=60s**: Record decision for ML training

## Performance Requirements

| Metric | Target | Critical |
|--------|--------|----------|
| Data ingestion latency | < 5s | < 10s |
| ML prediction time | < 10s | < 20s |
| Optimization solve time | < 30s | < 60s |
| Total cycle time | < 60s | < 90s |
| Prediction accuracy (±5 min) | > 85% | > 75% |
| Conflict detection rate | > 95% | > 90% |
| System availability | > 99.5% | > 99% |

## Database Schema

### Tables

1. **trips** - Historical trip records
2. **vehicle_positions** - Real-time position history
3. **conflicts** - Detected conflicts log
4. **decisions** - Dispatch decisions and outcomes
5. **performance_metrics** - Aggregated performance data
6. **trainsets** - Vehicle roster and characteristics

## API Endpoints

### Internal
- `GET /api/v1/state` - Current network state
- `GET /api/v1/conflicts` - Active conflicts
- `GET /api/v1/recommendations` - Current recommendations
- `POST /api/v1/feedback` - Manual feedback from dispatchers

### External (Public)
- `GET /api/v1/trips/{trip_id}` - Trip status
- `GET /api/v1/arrivals/{station_id}` - Station arrivals
- `GET /api/v1/delays` - System-wide delay summary

## Technology Stack

- **Language**: Go 1.21+
- **Optimization**: MIP solver (Gurobi/CPLEX or open-source alternatives)
- **ML Models**: Python scikit-learn, TensorFlow (via gRPC)
- **Database**: PostgreSQL with TimescaleDB
- **Message Queue**: Redis for real-time updates
- **Monitoring**: Prometheus + Grafana
- **Deployment**: Docker + Kubernetes

## Scalability Considerations

- **Horizontal scaling**: Separate ingestion, ML, and optimization services
- **Caching**: Redis for frequently accessed network state
- **Batch processing**: Offline ML training on historical data
- **Regional deployment**: Support for other Amtrak corridors (Michigan, Pacific Surfliner)

## Future Enhancements

1. **Passenger demand modeling**: Integrate ridership forecasts
2. **Multi-corridor optimization**: Coordinate across connected corridors
3. **Energy optimization**: Minimize electricity consumption
4. **Crew optimization**: Joint crew and train scheduling
5. **Predictive maintenance**: Integrate vehicle health data
6. **Weather integration**: Real-time weather impact modeling
