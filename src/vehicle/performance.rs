// Vehicle performance profiles

use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct TrainsetType {
    pub name: String,
    pub max_speed_kmh: f64,
    pub max_acceleration: f64,  // m/s²
    pub max_deceleration: f64,  // m/s²
    pub length_m: f64,
    pub tare_weight_kg: f64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PerformanceProfile {
    pub trainset: TrainsetType,
    // TODO: Acceleration curves, braking curves, power curves
}

pub struct VehicleRegistry {
    profiles: std::collections::HashMap<String, PerformanceProfile>,
}

impl VehicleRegistry {
    pub fn new() -> Self {
        let mut profiles = std::collections::HashMap::new();

        // Acela
        profiles.insert(
            "Acela".to_string(),
            PerformanceProfile {
                trainset: TrainsetType {
                    name: "Acela".to_string(),
                    max_speed_kmh: 240.0,
                    max_acceleration: 0.8,
                    max_deceleration: 1.2,
                    length_m: 202.0,
                    tare_weight_kg: 560_000.0,
                },
            },
        );

        // ACS-64 (Electric locomotive)
        profiles.insert(
            "ACS-64".to_string(),
            PerformanceProfile {
                trainset: TrainsetType {
                    name: "ACS-64".to_string(),
                    max_speed_kmh: 201.0,
                    max_acceleration: 0.6,
                    max_deceleration: 1.0,
                    length_m: 21.3,
                    tare_weight_kg: 96_000.0,
                },
            },
        );

        Self { profiles }
    }

    pub fn get(&self, name: &str) -> Option<&PerformanceProfile> {
        self.profiles.get(name)
    }
}
