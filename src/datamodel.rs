// Core data models

use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use uuid::Uuid;

/// Vehicle position from GTFS-Realtime
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct VehiclePosition {
    pub id: Uuid,
    pub vehicle_id: String,
    pub trip_id: String,
    pub route_id: String,
    pub timestamp: DateTime<Utc>,
    pub latitude: f64,
    pub longitude: f64,
    pub speed: Option<f64>,
    pub heading: Option<f64>,
    pub current_stop_sequence: Option<i32>,
    pub current_status: Option<String>,
}

/// Trip update with stop time predictions
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct TripUpdate {
    pub trip_id: String,
    pub route_id: String,
    pub timestamp: DateTime<Utc>,
    pub delay: i32,  // seconds
    pub stop_time_updates: Vec<StopTimeUpdate>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct StopTimeUpdate {
    pub stop_id: String,
    pub stop_sequence: i32,
    pub arrival_delay: Option<i32>,
    pub departure_delay: Option<i32>,
    pub arrival_time: Option<DateTime<Utc>>,
    pub departure_time: Option<DateTime<Utc>>,
}

/// Service alert
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ServiceAlert {
    pub id: String,
    pub cause: String,
    pub effect: String,
    pub header_text: String,
    pub description_text: String,
    pub affected_routes: Vec<String>,
    pub active_period_start: DateTime<Utc>,
    pub active_period_end: Option<DateTime<Utc>>,
}

/// Network state snapshot
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct NetworkState {
    pub timestamp: DateTime<Utc>,
    pub active_trips: Vec<ActiveTrip>,
    pub conflicts: Vec<Conflict>,
    pub recommendations: Vec<Recommendation>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ActiveTrip {
    pub trip_id: String,
    pub route_id: String,
    pub vehicle_id: String,
    pub current_position: (f64, f64),  // (lat, lon)
    pub current_block_id: Option<String>,
    pub current_delay: i32,
    pub predicted_delay: i32,
    pub status: TripStatus,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum TripStatus {
    Scheduled,
    InProgress,
    Delayed,
    OnTime,
    Completed,
    Cancelled,
}

/// Detected conflict between trains
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Conflict {
    pub id: Uuid,
    pub conflict_type: ConflictType,
    pub involved_trips: Vec<String>,
    pub block_id: String,
    pub estimated_time: DateTime<Utc>,
    pub severity: f64,  // 0.0 - 1.0
    pub estimated_impact_secs: i32,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum ConflictType {
    HeadToHead,
    Following,
    Crossing,
    PlatformOccupancy,
}

/// Dispatch recommendation from optimization
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Recommendation {
    pub id: Uuid,
    pub trip_id: String,
    pub action_type: ActionType,
    pub hold_duration_secs: Option<i32>,
    pub target_block: Option<String>,
    pub expected_benefit: f64,
    pub confidence: f64,
    pub reason: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum ActionType {
    HoldAtStation,
    ExpediteDispatch,
    ChangeTrack,
    SlowApproach,
    NoAction,
}
