// Track modeling module

use geo_types::{Coord, Point};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

/// Signal block - fundamental unit for conflict detection
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SignalBlock {
    pub block_id: String,
    pub control_point: Option<String>,
    pub start_milepost: f64,
    pub end_milepost: f64,
    pub length: f64,  // km
    pub track: String,
    pub direction: Direction,
    pub max_speed: f64,  // km/h
    pub capacity: u8,  // Usually 1 train
    pub electrified: bool,
    pub start_location: Point<f64>,
    pub end_location: Point<f64>,
    pub next_blocks: Vec<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq)]
pub enum Direction {
    Northbound,
    Southbound,
    Eastbound,
    Westbound,
    Bidirectional,
}

/// Interlocking - complex junction area
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Interlocking {
    pub interlocking_id: String,
    pub name: String,
    pub location: Point<f64>,
    pub milepost: f64,
    pub routes: Vec<Route>,
}

/// Route through an interlocking
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Route {
    pub route_id: String,
    pub from_block: String,
    pub to_block: String,
    pub max_speed: f64,
    pub locking_time_secs: u32,
    pub conflicts_with: Vec<String>,  // Other route IDs
}

/// Track network with spatial indexing
pub struct TrackNetwork {
    blocks: HashMap<String, SignalBlock>,
    interlockings: HashMap<String, Interlocking>,
    // R-tree for spatial queries (find blocks near position)
    spatial_index: rstar::RTree<BlockEntry>,
}

#[derive(Debug, Clone)]
struct BlockEntry {
    block_id: String,
    center: Point<f64>,
}

impl rstar::RTreeObject for BlockEntry {
    type Envelope = rstar::AABB<[f64; 2]>;

    fn envelope(&self) -> Self::Envelope {
        rstar::AABB::from_point([self.center.x(), self.center.y()])
    }
}

impl rstar::PointDistance for BlockEntry {
    fn distance_2(&self, point: &[f64; 2]) -> f64 {
        let dx = self.center.x() - point[0];
        let dy = self.center.y() - point[1];
        dx * dx + dy * dy
    }
}

impl TrackNetwork {
    pub fn new() -> Self {
        Self {
            blocks: HashMap::new(),
            interlockings: HashMap::new(),
            spatial_index: rstar::RTree::new(),
        }
    }

    pub fn add_block(&mut self, block: SignalBlock) {
        let center = Point::new(
            (block.start_location.x() + block.end_location.x()) / 2.0,
            (block.start_location.y() + block.end_location.y()) / 2.0,
        );

        let entry = BlockEntry {
            block_id: block.block_id.clone(),
            center,
        };

        self.spatial_index.insert(entry);
        self.blocks.insert(block.block_id.clone(), block);
    }

    pub fn add_interlocking(&mut self, interlocking: Interlocking) {
        self.interlockings.insert(interlocking.interlocking_id.clone(), interlocking);
    }

    pub fn get_block(&self, block_id: &str) -> Option<&SignalBlock> {
        self.blocks.get(block_id)
    }

    /// Find blocks near a geographic position
    pub fn find_blocks_near(&self, lat: f64, lon: f64, max_distance_km: f64) -> Vec<&SignalBlock> {
        let query_point = [lon, lat];
        
        self.spatial_index
            .nearest_neighbor_iter(&query_point)
            .take(10)  // Top 10 nearest
            .filter_map(|entry| {
                let block = self.blocks.get(&entry.block_id)?;
                let distance = haversine_distance(lat, lon, entry.center.y(), entry.center.x());
                if distance <= max_distance_km {
                    Some(block)
                } else {
                    None
                }
            })
            .collect()
    }

    /// Get next blocks in the route
    pub fn get_next_blocks(&self, current_block_id: &str) -> Vec<&SignalBlock> {
        if let Some(block) = self.blocks.get(current_block_id) {
            block.next_blocks.iter()
                .filter_map(|id| self.blocks.get(id))
                .collect()
        } else {
            Vec::new()
        }
    }
}

/// Calculate distance between two lat/lon points (Haversine formula)
pub fn haversine_distance(lat1: f64, lon1: f64, lat2: f64, lon2: f64) -> f64 {
    const EARTH_RADIUS_KM: f64 = 6371.0;

    let lat1_rad = lat1.to_radians();
    let lat2_rad = lat2.to_radians();
    let delta_lat = (lat2 - lat1).to_radians();
    let delta_lon = (lon2 - lon1).to_radians();

    let a = (delta_lat / 2.0).sin().powi(2)
        + lat1_rad.cos() * lat2_rad.cos() * (delta_lon / 2.0).sin().powi(2);
    let c = 2.0 * a.sqrt().atan2((1.0 - a).sqrt());

    EARTH_RADIUS_KM * c
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_haversine_distance() {
        // NYP to PHL (approx 90 miles = 145 km)
        let dist = haversine_distance(40.7505, -73.9934, 39.9566, -75.1822);
        assert!((dist - 145.0).abs() < 5.0);
    }
}
