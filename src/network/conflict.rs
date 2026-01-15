// Conflict detection logic

use crate::{datamodel::*, track::TrackNetwork};
use chrono::{DateTime, Utc, Duration};
use std::collections::HashMap;
use tracing::debug;

pub struct ConflictDetector {
    track_network: TrackNetwork,
}

impl ConflictDetector {
    pub fn new(track_network: TrackNetwork) -> Self {
        Self { track_network }
    }

    /// Detect conflicts between active trips
    pub fn detect_conflicts(&self, trips: &[ActiveTrip]) -> Vec<Conflict> {
        let mut conflicts = Vec::new();

        // Build block occupancy map
        let mut block_occupancy: HashMap<String, Vec<&ActiveTrip>> = HashMap::new();
        for trip in trips {
            if let Some(ref block_id) = trip.current_block_id {
                block_occupancy.entry(block_id.clone())
                    .or_insert_with(Vec::new)
                    .push(trip);
            }
        }

        // Check for same-block conflicts
        for (block_id, occupants) in block_occupancy.iter() {
            if occupants.len() > 1 {
                conflicts.push(Conflict {
                    id: uuid::Uuid::new_v4(),
                    conflict_type: ConflictType::Following,
                    involved_trips: occupants.iter().map(|t| t.trip_id.clone()).collect(),
                    block_id: block_id.clone(),
                    estimated_time: Utc::now(),
                    severity: 0.9,
                    estimated_impact_secs: 300,
                });
            }
        }

        debug!("Detected {} conflicts", conflicts.len());
        conflicts
    }
}
