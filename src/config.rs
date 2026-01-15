// Configuration management

use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct Config {
    pub database: DatabaseConfig,
    pub gtfs_rt: GtfsRtConfig,
    pub optimization: OptimizationConfig,
    pub ml: MlConfig,
    pub api: ApiConfig,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct DatabaseConfig {
    pub url: String,
    pub max_connections: u32,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct GtfsRtConfig {
    pub vehicle_positions_url: String,
    pub trip_updates_url: String,
    pub service_alerts_url: String,
    pub refresh_interval_secs: u64,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct OptimizationConfig {
    pub max_solve_time_secs: u64,
    pub conflict_horizon_hours: f64,
    pub minimum_headway_secs: u64,
    pub weights: Weights,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct Weights {
    pub total_delay: f64,
    pub max_delay: f64,
    pub passenger_impact: f64,
    pub on_time_performance: f64,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct MlConfig {
    pub delay_predictor_path: String,
    pub conflict_predictor_path: String,
    pub dispatch_advisor_path: String,
    pub performance_analyzer_path: String,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct ApiConfig {
    pub host: String,
    pub port: u16,
}

impl Config {
    pub fn load() -> Result<Self, config::ConfigError> {
        let config = config::Config::builder()
            .add_source(config::File::with_name("config/config"))
            .add_source(config::Environment::with_prefix("AMTRK"))
            .build()?;

        config.try_deserialize()
    }
}
