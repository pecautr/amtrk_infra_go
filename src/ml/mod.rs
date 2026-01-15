// ML inference module (ONNX models trained in Python)

use crate::{datamodel::*, error::Result, config::MlConfig};
use std::collections::HashMap;
use tract_onnx::prelude::*;
use tracing::{debug, info};

pub struct MlEngine {
    delay_predictor: Option<SimplePlan<TypedFact, Box<dyn TypedOp>, Graph<TypedFact, Box<dyn TypedOp>>>>,
    conflict_predictor: Option<SimplePlan<TypedFact, Box<dyn TypedOp>, Graph<TypedFact, Box<dyn TypedOp>>>>,
    config: MlConfig,
}

pub struct MlPredictions {
    pub delay_predictions: HashMap<String, i32>,  // trip_id -> predicted delay (secs)
    pub conflict_predictions: HashMap<String, ConflictPrediction>,
}

#[derive(Debug, Clone)]
pub struct ConflictPrediction {
    pub conflict_id: uuid::Uuid,
    pub probability: f64,
    pub severity: f64,
}

impl MlEngine {
    pub fn new(config: MlConfig) -> Result<Self> {
        info!("Loading ML models...");

        // Load ONNX models (trained in Python)
        let delay_predictor = if std::path::Path::new(&config.delay_predictor_path).exists() {
            debug!("Loading delay predictor from {}", config.delay_predictor_path);
            let model = tract_onnx::onnx()
                .model_for_path(&config.delay_predictor_path)
                .map_err(|e| crate::Error::ML(format!("Failed to load delay predictor: {}", e)))?
                .into_optimized()
                .map_err(|e| crate::Error::ML(format!("Failed to optimize model: {}", e)))?
                .into_runnable()
                .map_err(|e| crate::Error::ML(format!("Failed to make model runnable: {}", e)))?;
            Some(model)
        } else {
            debug!("Delay predictor not found at {}, using fallback", config.delay_predictor_path);
            None
        };

        let conflict_predictor = if std::path::Path::new(&config.conflict_predictor_path).exists() {
            debug!("Loading conflict predictor from {}", config.conflict_predictor_path);
            let model = tract_onnx::onnx()
                .model_for_path(&config.conflict_predictor_path)
                .map_err(|e| crate::Error::ML(format!("Failed to load conflict predictor: {}", e)))?
                .into_optimized()
                .map_err(|e| crate::Error::ML(format!("Failed to optimize model: {}", e)))?
                .into_runnable()
                .map_err(|e| crate::Error::ML(format!("Failed to make model runnable: {}", e)))?;
            Some(model)
        } else {
            debug!("Conflict predictor not found at {}, using fallback", config.conflict_predictor_path);
            None
        };

        Ok(Self {
            delay_predictor,
            conflict_predictor,
            config,
        })
    }

    /// Predict delays and conflicts using ML models
    pub async fn predict(
        &self,
        active_trips: &[ActiveTrip],
        detected_conflicts: &[Conflict],
    ) -> Result<MlPredictions> {
        info!("Running ML predictions for {} trips, {} conflicts",
            active_trips.len(), detected_conflicts.len());

        let delay_predictions = self.predict_delays(active_trips).await?;
        let conflict_predictions = self.predict_conflict_severity(detected_conflicts).await?;

        Ok(MlPredictions {
            delay_predictions,
            conflict_predictions,
        })
    }

    async fn predict_delays(&self, trips: &[ActiveTrip]) -> Result<HashMap<String, i32>> {
        let mut predictions = HashMap::new();

        if let Some(ref model) = self.delay_predictor {
            for trip in trips {
                // Extract features (simplified - in production would use more features)
                let features = vec![
                    trip.current_delay as f32,
                    trip.predicted_delay as f32,
                    // Add: time of day, day of week, route, weather, etc.
                ];

                let input = tract_ndarray::Array2::from_shape_vec(
                    (1, features.len()),
                    features
                ).map_err(|e| crate::Error::ML(format!("Failed to create input tensor: {}", e)))?;

                let result = model.run(tvec!(input.into()))
                    .map_err(|e| crate::Error::ML(format!("Model inference failed: {}", e)))?;

                let predicted_delay = result[0]
                    .to_array_view::<f32>()
                    .map_err(|e| crate::Error::ML(format!("Failed to extract output: {}", e)))?
                    [[0, 0]] as i32;

                predictions.insert(trip.trip_id.clone(), predicted_delay);
            }
        } else {
            // Fallback: simple heuristic (current delay + 10%)
            for trip in trips {
                predictions.insert(
                    trip.trip_id.clone(),
                    (trip.current_delay as f32 * 1.1) as i32
                );
            }
        }

        debug!("Predicted delays for {} trips", predictions.len());
        Ok(predictions)
    }

    async fn predict_conflict_severity(
        &self,
        conflicts: &[Conflict],
    ) -> Result<HashMap<String, ConflictPrediction>> {
        let mut predictions = HashMap::new();

        if let Some(ref model) = self.conflict_predictor {
            for conflict in conflicts {
                // Extract features
                let features = vec![
                    conflict.severity as f32,
                    conflict.estimated_impact_secs as f32,
                    conflict.involved_trips.len() as f32,
                    // Add: block characteristics, historical data, etc.
                ];

                let input = tract_ndarray::Array2::from_shape_vec(
                    (1, features.len()),
                    features
                ).map_err(|e| crate::Error::ML(format!("Failed to create input tensor: {}", e)))?;

                let result = model.run(tvec!(input.into()))
                    .map_err(|e| crate::Error::ML(format!("Model inference failed: {}", e)))?;

                let output = result[0]
                    .to_array_view::<f32>()
                    .map_err(|e| crate::Error::ML(format!("Failed to extract output: {}", e)))?;

                predictions.insert(
                    conflict.id.to_string(),
                    ConflictPrediction {
                        conflict_id: conflict.id,
                        probability: output[[0, 0]] as f64,
                        severity: output[[0, 1]] as f64,
                    },
                );
            }
        } else {
            // Fallback: use existing severity
            for conflict in conflicts {
                predictions.insert(
                    conflict.id.to_string(),
                    ConflictPrediction {
                        conflict_id: conflict.id,
                        probability: conflict.severity,
                        severity: conflict.severity,
                    },
                );
            }
        }

        debug!("Predicted conflict severity for {} conflicts", predictions.len());
        Ok(predictions)
    }
}
