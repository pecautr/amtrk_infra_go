-- Database schema for Amtrak NEC Optimization System

-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "timescaledb";

-- Trips table
CREATE TABLE trips (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    trip_id VARCHAR(50) NOT NULL,
    route_id VARCHAR(50) NOT NULL,
    vehicle_id VARCHAR(50),
    trainset_type VARCHAR(50),
    direction VARCHAR(20),
    status VARCHAR(20),
    scheduled_start TIMESTAMP NOT NULL,
    scheduled_end TIMESTAMP NOT NULL,
    actual_start TIMESTAMP,
    actual_end TIMESTAMP,
    total_delay_seconds INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_trips_trip_id ON trips(trip_id);
CREATE INDEX idx_trips_scheduled_start ON trips(scheduled_start);

-- Vehicle positions (time-series data)
CREATE TABLE vehicle_positions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    vehicle_id VARCHAR(50) NOT NULL,
    trip_id VARCHAR(50) NOT NULL,
    route_id VARCHAR(50) NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    speed DOUBLE PRECISION,
    heading DOUBLE PRECISION,
    current_stop VARCHAR(50),
    stop_status INTEGER,
    occupancy INTEGER
);

-- Convert to TimescaleDB hypertable
SELECT create_hypertable('vehicle_positions', 'timestamp', if_not_exists => TRUE);

CREATE INDEX idx_vehicle_positions_vehicle_id ON vehicle_positions(vehicle_id, timestamp DESC);
CREATE INDEX idx_vehicle_positions_trip_id ON vehicle_positions(trip_id, timestamp DESC);

-- Trip updates
CREATE TABLE trip_updates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    trip_id VARCHAR(50) NOT NULL,
    route_id VARCHAR(50) NOT NULL,
    vehicle_id VARCHAR(50),
    timestamp TIMESTAMP NOT NULL,
    delay_seconds INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

SELECT create_hypertable('trip_updates', 'timestamp', if_not_exists => TRUE);

-- Stop time updates
CREATE TABLE stop_time_updates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    trip_update_id UUID REFERENCES trip_updates(id),
    stop_sequence INTEGER NOT NULL,
    stop_id VARCHAR(50) NOT NULL,
    arrival_delay INTEGER,
    departure_delay INTEGER,
    arrival_time TIMESTAMP,
    departure_time TIMESTAMP,
    schedule_relationship INTEGER
);

-- Service alerts
CREATE TABLE service_alerts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    timestamp TIMESTAMP NOT NULL,
    cause VARCHAR(100),
    effect VARCHAR(100),
    header_text TEXT,
    description TEXT,
    active_start TIMESTAMP,
    active_end TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Affected routes/stops for alerts
CREATE TABLE alert_affected_entities (
    alert_id UUID REFERENCES service_alerts(id),
    entity_type VARCHAR(20), -- 'route' or 'stop'
    entity_id VARCHAR(50),
    PRIMARY KEY (alert_id, entity_type, entity_id)
);

-- Conflicts
CREATE TABLE conflicts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    conflict_id VARCHAR(100) UNIQUE NOT NULL,
    detected_at TIMESTAMP NOT NULL,
    conflict_time TIMESTAMP NOT NULL,
    conflict_type INTEGER NOT NULL,
    location VARCHAR(100),
    severity INTEGER,
    estimated_impact_seconds INTEGER,
    resolution TEXT,
    resolved_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_conflicts_detected_at ON conflicts(detected_at);
CREATE INDEX idx_conflicts_conflict_time ON conflicts(conflict_time);

-- Involved trips in conflicts
CREATE TABLE conflict_trips (
    conflict_id UUID REFERENCES conflicts(id),
    trip_id VARCHAR(50) NOT NULL,
    PRIMARY KEY (conflict_id, trip_id)
);

-- Recommendations
CREATE TABLE recommendations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    recommendation_id VARCHAR(100) UNIQUE NOT NULL,
    generated_at TIMESTAMP NOT NULL,
    action_type INTEGER NOT NULL,
    trip_id VARCHAR(50) NOT NULL,
    stop_id VARCHAR(50),
    hold_duration_seconds INTEGER,
    alternative_track VARCHAR(50),
    priority INTEGER,
    expected_benefit_seconds INTEGER,
    confidence DOUBLE PRECISION,
    applied BOOLEAN DEFAULT FALSE,
    applied_at TIMESTAMP,
    outcome_delay_seconds INTEGER,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_recommendations_generated_at ON recommendations(generated_at);
CREATE INDEX idx_recommendations_trip_id ON recommendations(trip_id);

-- Historical dispatch decisions (for ML training)
CREATE TABLE dispatch_decisions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    timestamp TIMESTAMP NOT NULL,
    network_state_snapshot JSONB,
    decision_action VARCHAR(50) NOT NULL,
    trip_id VARCHAR(50) NOT NULL,
    conflict_context TEXT,
    
    -- Outcomes
    actual_delay_seconds INTEGER,
    network_impact_seconds INTEGER,
    success_score DOUBLE PRECISION,
    
    -- Context features
    time_of_day INTEGER,
    day_of_week INTEGER,
    train_type VARCHAR(50),
    current_load INTEGER,
    weather_condition VARCHAR(50),
    
    created_at TIMESTAMP DEFAULT NOW()
);

SELECT create_hypertable('dispatch_decisions', 'timestamp', if_not_exists => TRUE);

CREATE INDEX idx_dispatch_decisions_trip_id ON dispatch_decisions(trip_id);
CREATE INDEX idx_dispatch_decisions_success_score ON dispatch_decisions(success_score);

-- Performance metrics (aggregated)
CREATE TABLE performance_metrics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    period_start TIMESTAMP NOT NULL,
    period_end TIMESTAMP NOT NULL,
    
    total_trips INTEGER,
    on_time_trips INTEGER,
    delayed_trips INTEGER,
    cancelled_trips INTEGER,
    
    average_delay_seconds DOUBLE PRECISION,
    max_delay_seconds INTEGER,
    total_delay_minutes INTEGER,
    
    conflicts_detected INTEGER,
    conflicts_resolved INTEGER,
    
    otp_percentage DOUBLE PRECISION,
    
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_performance_metrics_period ON performance_metrics(period_start, period_end);

-- Trainset registry
CREATE TABLE trainsets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    type_id VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    manufacturer VARCHAR(100),
    max_speed DOUBLE PRECISION,
    service_speed DOUBLE PRECISION,
    length DOUBLE PRECISION,
    capacity INTEGER,
    power_type VARCHAR(50)
);

-- Vehicle metrics (track individual vehicle performance)
CREATE TABLE vehicle_metrics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    vehicle_id VARCHAR(50) NOT NULL,
    trainset_type_id VARCHAR(50) REFERENCES trainsets(type_id),
    timestamp TIMESTAMP NOT NULL,
    
    current_speed DOUBLE PRECISION,
    current_acceleration DOUBLE PRECISION,
    power_draw DOUBLE PRECISION,
    
    average_speed DOUBLE PRECISION,
    max_speed_achieved DOUBLE PRECISION,
    
    performance_factor DOUBLE PRECISION,
    maintenance_status VARCHAR(50),
    last_maintenance TIMESTAMP,
    
    total_distance DOUBLE PRECISION,
    total_trips INTEGER,
    average_delay_seconds DOUBLE PRECISION
);

SELECT create_hypertable('vehicle_metrics', 'timestamp', if_not_exists => TRUE);

CREATE INDEX idx_vehicle_metrics_vehicle_id ON vehicle_metrics(vehicle_id, timestamp DESC);

-- Stations
CREATE TABLE stations (
    id VARCHAR(10) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    platforms INTEGER,
    capacity INTEGER,
    region VARCHAR(50)
);

-- Insert NEC stations
INSERT INTO stations (id, name, platforms, capacity, region) VALUES
    ('BOS', 'Boston South Station', 12, 8, 'NEC'),
    ('BBY', 'Boston Back Bay', 5, 3, 'NEC'),
    ('RTE', 'Route 128', 4, 3, 'NEC'),
    ('PVD', 'Providence', 4, 3, 'NEC'),
    ('NHV', 'New Haven', 6, 4, 'NEC'),
    ('STM', 'Stamford', 6, 4, 'NEC'),
    ('NYP', 'New York Penn Station', 21, 12, 'NEC'),
    ('EWR', 'Newark Penn Station', 5, 4, 'NEC'),
    ('MET', 'Metropark', 4, 2, 'NEC'),
    ('TRE', 'Trenton', 4, 3, 'NEC'),
    ('PHI', 'Philadelphia 30th Street', 16, 10, 'NEC'),
    ('WIL', 'Wilmington', 4, 3, 'NEC'),
    ('BAL', 'Baltimore Penn Station', 6, 4, 'NEC'),
    ('BWI', 'BWI Marshall Airport', 4, 3, 'NEC'),
    ('WAS', 'Washington Union Station', 22, 14, 'NEC');

-- Create views for common queries

-- Active trips view
CREATE VIEW active_trips_view AS
SELECT 
    t.trip_id,
    t.route_id,
    t.vehicle_id,
    t.trainset_type,
    t.status,
    t.scheduled_start,
    t.total_delay_seconds,
    vp.latitude,
    vp.longitude,
    vp.speed,
    vp.current_stop
FROM trips t
LEFT JOIN LATERAL (
    SELECT * FROM vehicle_positions 
    WHERE vehicle_id = t.vehicle_id 
    ORDER BY timestamp DESC 
    LIMIT 1
) vp ON true
WHERE t.status IN ('active', 'delayed')
AND t.scheduled_start >= NOW() - INTERVAL '2 hours'
AND (t.actual_end IS NULL OR t.actual_end >= NOW() - INTERVAL '30 minutes');

-- Recent conflicts view
CREATE VIEW recent_conflicts_view AS
SELECT 
    c.conflict_id,
    c.detected_at,
    c.conflict_time,
    c.conflict_type,
    c.location,
    c.severity,
    c.resolution,
    array_agg(ct.trip_id) as involved_trips
FROM conflicts c
JOIN conflict_trips ct ON c.id = ct.conflict_id
WHERE c.detected_at >= NOW() - INTERVAL '1 hour'
GROUP BY c.id, c.conflict_id, c.detected_at, c.conflict_time, 
         c.conflict_type, c.location, c.severity, c.resolution;

-- Functions

-- Update trip delay
CREATE OR REPLACE FUNCTION update_trip_delay()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE trips
    SET total_delay_seconds = NEW.delay_seconds,
        updated_at = NOW()
    WHERE trip_id = NEW.trip_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_trip_delay
AFTER INSERT ON trip_updates
FOR EACH ROW
EXECUTE FUNCTION update_trip_delay();

-- Calculate OTP
CREATE OR REPLACE FUNCTION calculate_otp(start_time TIMESTAMP, end_time TIMESTAMP)
RETURNS DOUBLE PRECISION AS $$
DECLARE
    total INTEGER;
    on_time INTEGER;
BEGIN
    SELECT COUNT(*) INTO total
    FROM trips
    WHERE scheduled_start >= start_time AND scheduled_start < end_time
    AND status = 'completed';
    
    IF total = 0 THEN
        RETURN 0;
    END IF;
    
    SELECT COUNT(*) INTO on_time
    FROM trips
    WHERE scheduled_start >= start_time AND scheduled_start < end_time
    AND status = 'completed'
    AND total_delay_seconds <= 300; -- 5 minutes tolerance
    
    RETURN (on_time::DOUBLE PRECISION / total::DOUBLE PRECISION) * 100.0;
END;
$$ LANGUAGE plpgsql;

COMMENT ON DATABASE amtrak_nec IS 'Amtrak Northeast Corridor Optimization System';
