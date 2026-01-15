// Optimization module using good_lp

use crate::{datamodel::*, error::Result, config::OptimizationConfig};
use good_lp::*;
use std::collections::HashMap;
use tracing::{debug, info};

pub struct Optimizer {
    config: OptimizationConfig,
}

pub struct OptimizationRequest {
    pub active_trips: Vec<ActiveTrip>,
    pub predicted_conflicts: Vec<Conflict>,
    pub ml_delay_predictions: HashMap<String, i32>,  // trip_id -> predicted delay
    pub ml_conflict_probabilities: HashMap<String, f64>,  // conflict_id -> probability
}

pub struct OptimizationResult {
    pub recommendations: Vec<Recommendation>,
    pub objective_value: f64,
    pub solve_time_secs: f64,
}

impl Optimizer {
    pub fn new(config: OptimizationConfig) -> Self {
        Self { config }
    }

    pub async fn optimize(&self, request: OptimizationRequest) -> Result<OptimizationResult> {
        let start = std::time::Instant::now();

        info!("Starting optimization for {} trips, {} conflicts",
            request.active_trips.len(), request.predicted_conflicts.len());

        // Create optimization problem
        let mut vars = ProblemVariables::new();
        let trip_ids: Vec<String> = request.active_trips.iter()
            .map(|t| t.trip_id.clone())
            .collect();

        // Decision variables: hold duration for each trip (0-600 seconds)
        let hold_vars: Vec<Variable> = trip_ids.iter()
            .map(|_| vars.add(variable().min(0.0).max(600.0)))
            .collect();

        // Build objective function: minimize total delay + max delay + passenger impact
        let w = &self.config.weights;
        
        let mut objective = Expression::from(0.0);
        
        for (i, trip) in request.active_trips.iter().enumerate() {
            let current_delay = trip.current_delay as f64;
            let predicted_delay = *request.ml_delay_predictions
                .get(&trip.trip_id)
                .unwrap_or(&0) as f64;
            
            // Total delay component
            let delay_expr = (current_delay + predicted_delay) * w.total_delay;
            
            // Passenger impact (simplified - would be based on ridership data)
            let passenger_impact = 1.0 * w.passenger_impact;
            
            objective += delay_expr * hold_vars[i] + passenger_impact * hold_vars[i];
        }

        // Constraints: ensure conflicts are resolved
        for conflict in &request.predicted_conflicts {
            if conflict.involved_trips.len() >= 2 {
                let trip1_idx = trip_ids.iter().position(|id| id == &conflict.involved_trips[0]);
                let trip2_idx = trip_ids.iter().position(|id| id == &conflict.involved_trips[1]);

                if let (Some(i1), Some(i2)) = (trip1_idx, trip2_idx) {
                    // At least one train must be held long enough to avoid conflict
                    // Simplified: hold_1 + hold_2 >= minimum_headway
                    let min_headway = self.config.minimum_headway_secs as f64;
                    vars.add_constraint(hold_vars[i1] + hold_vars[i2] >= min_headway);
                }
            }
        }

        // Solve using HiGHS solver
        debug!("Solving MIP problem...");
        let solution = vars.minimise(objective)
            .using(default_solver)
            .solve()
            .map_err(|e| crate::Error::Optimization(format!("Solver failed: {:?}", e)))?;

        // Extract recommendations
        let mut recommendations = Vec::new();
        for (i, trip) in request.active_trips.iter().enumerate() {
            let hold_duration = solution.value(hold_vars[i]) as i32;

            if hold_duration > 30 {  // Only recommend if > 30 seconds
                let action_type = if hold_duration > 300 {
                    ActionType::ChangeTrack
                } else if hold_duration > 120 {
                    ActionType::HoldAtStation
                } else {
                    ActionType::SlowApproach
                };

                recommendations.push(Recommendation {
                    id: uuid::Uuid::new_v4(),
                    trip_id: trip.trip_id.clone(),
                    action_type,
                    hold_duration_secs: Some(hold_duration),
                    target_block: None,
                    expected_benefit: 0.8,  // Placeholder
                    confidence: 0.85,       // Placeholder
                    reason: format!("Avoid conflict, hold {} seconds", hold_duration),
                });
            }
        }

        let elapsed = start.elapsed();
        info!("Optimization complete: {} recommendations in {:.2}s", 
            recommendations.len(), elapsed.as_secs_f64());

        Ok(OptimizationResult {
            recommendations,
            objective_value: solution.eval(objective),
            solve_time_secs: elapsed.as_secs_f64(),
        })
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::config::Weights;

    #[test]
    fn test_basic_optimization() {
        let config = OptimizationConfig {
            max_solve_time_secs: 30,
            conflict_horizon_hours: 2.0,
            minimum_headway_secs: 180,
            weights: Weights {
                total_delay: 0.4,
                max_delay: 0.3,
                passenger_impact: 0.2,
                on_time_performance: 0.1,
            },
        };

        let optimizer = Optimizer::new(config);
        // Test with minimal problem...
    }
}
