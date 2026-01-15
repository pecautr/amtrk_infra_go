# Machine Learning Models Documentation

## Overview

The ML Engine provides predictive capabilities to enhance optimization decisions. It consists of four primary models that work together to forecast delays, predict conflicts, and recommend optimal dispatch decisions.

## Model Architecture

```
Historical Data → Feature Extraction → ML Models → Predictions → Optimization Input
```

## 1. Delay Predictor

### Purpose
Forecast how current delays will evolve over the next hour.

### Input Features
- Current delay (seconds)
- Time of day (hour, 0-23)
- Day of week (0-6)
- Weather conditions (temperature, precipitation, wind)
- Network load (number of active trains)
- Train type (Acela, Regional, etc.)
- Station occupancy rate
- Recent delay trend (last 15 minutes)

### Output
- Predicted delay in 60 minutes
- Confidence score (0-1)
- Contributing factors with impact weights

### Model Type
**Gradient Boosting Regressor** (XGBoost or LightGBM)

### Training Data
Historical trip data with features:
```
[trip_id, timestamp, current_delay, hour, dow, weather, ..., actual_delay_60min]
```

### Performance Metrics
- MAE (Mean Absolute Error): Target < 120 seconds
- RMSE (Root Mean Squared Error): Target < 180 seconds
- Accuracy within ±5 minutes: Target > 85%

### Training Code (Python)
```python
import pandas as pd
from xgboost import XGBRegressor
from sklearn.model_selection import train_test_split

# Load training data
data = pd.read_csv('historical_delays.csv')

features = ['current_delay', 'hour', 'dow', 'weather_temp', 
            'weather_precip', 'network_load', 'train_type_encoded']
target = 'delay_60min'

X = data[features]
y = data[target]

X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2)

# Train model
model = XGBRegressor(
    n_estimators=200,
    max_depth=8,
    learning_rate=0.05,
    subsample=0.8
)
model.fit(X_train, y_train)

# Evaluate
predictions = model.predict(X_test)
mae = mean_absolute_error(y_test, predictions)
print(f"MAE: {mae:.2f} seconds")

# Save model
model.save_model('delay_predictor.pkl')
```

## 2. Conflict Predictor

### Purpose
Identify potential conflicts between trains before they occur.

### Input Features
- Distance between trains (meters)
- Relative speed difference (km/h)
- Track segment capacity
- Platform availability
- Scheduled arrival time difference
- Historical conflict rate for segment
- Current delays of involved trains
- Time of day
- Weather conditions

### Output
- Conflict probability (0-1)
- Estimated conflict time
- Conflict severity (1-10)
- Prevention cost (delay seconds needed to avoid)

### Model Type
**Random Forest Classifier** for probability estimation

### Training Data
Historical pairs of trains with labels:
```
[train_a_id, train_b_id, distance, speed_diff, ..., conflict_occurred (0/1)]
```

### Performance Metrics
- Precision: Target > 80% (avoid false alarms)
- Recall: Target > 95% (catch all real conflicts)
- F1 Score: Target > 0.85
- ROC AUC: Target > 0.90

### Training Code (Python)
```python
from sklearn.ensemble import RandomForestClassifier

features = ['distance', 'speed_diff', 'track_capacity', 'platform_avail',
            'time_diff', 'historical_conflict_rate', 'delay_a', 'delay_b']
target = 'conflict_occurred'

X = data[features]
y = data[target]

X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2)

model = RandomForestClassifier(
    n_estimators=300,
    max_depth=12,
    min_samples_split=20,
    class_weight='balanced'
)
model.fit(X_train, y_train)

# Evaluate
y_pred_proba = model.predict_proba(X_test)[:, 1]
roc_auc = roc_auc_score(y_test, y_pred_proba)
print(f"ROC AUC: {roc_auc:.3f}")

joblib.dump(model, 'conflict_predictor.pkl')
```

## 3. Dispatch Advisor

### Purpose
Recommend optimal hold times to minimize total network delay.

### Approach
**Reinforcement Learning** - Learn from outcomes of past decisions

### State Space
- Current network state (active trains, positions, delays)
- Predicted conflicts
- Station capacities
- Time of day

### Action Space
- Hold duration: 0, 60, 120, 180, 300, 600 seconds
- Per train decision

### Reward Function
```python
def calculate_reward(decision):
    # Negative reward for delays caused
    delay_penalty = -sum(train.delay for train in network.trains)
    
    # Negative reward for conflicts
    conflict_penalty = -500 * len(network.conflicts)
    
    # Positive reward for resolving conflicts
    resolved_reward = 200 * len(decision.resolved_conflicts)
    
    # Penalty for excessive holds
    hold_penalty = -0.5 * decision.hold_duration
    
    return delay_penalty + conflict_penalty + resolved_reward + hold_penalty
```

### Model Type
**Deep Q-Network (DQN)** with experience replay

### Training Code (Python)
```python
import torch
import torch.nn as nn
from collections import deque
import random

class DQN(nn.Module):
    def __init__(self, state_size, action_size):
        super(DQN, self).__init__()
        self.fc1 = nn.Linear(state_size, 128)
        self.fc2 = nn.Linear(128, 128)
        self.fc3 = nn.Linear(128, action_size)
    
    def forward(self, x):
        x = torch.relu(self.fc1(x))
        x = torch.relu(self.fc2(x))
        return self.fc3(x)

# Initialize
state_size = 50  # Network state features
action_size = 6  # Hold durations
model = DQN(state_size, action_size)
optimizer = torch.optim.Adam(model.parameters(), lr=0.001)
replay_buffer = deque(maxlen=10000)

# Training loop
for episode in range(1000):
    state = env.reset()
    total_reward = 0
    
    for step in range(100):
        # Epsilon-greedy action selection
        if random.random() < epsilon:
            action = random.randint(0, action_size-1)
        else:
            with torch.no_grad():
                q_values = model(torch.FloatTensor(state))
                action = q_values.argmax().item()
        
        # Take action, observe reward
        next_state, reward, done = env.step(action)
        replay_buffer.append((state, action, reward, next_state, done))
        
        # Train on batch
        if len(replay_buffer) > 32:
            batch = random.sample(replay_buffer, 32)
            # ... DQN update logic
        
        state = next_state
        total_reward += reward
        
        if done:
            break
    
    print(f"Episode {episode}: Reward = {total_reward}")

torch.save(model.state_dict(), 'dispatch_advisor.pt')
```

## 4. Performance Analyzer

### Purpose
Continuous learning from historical dispatch decisions to improve system parameters.

### Analysis Tasks

1. **Decision Quality Assessment**
   - Compare predicted vs. actual outcomes
   - Identify patterns in successful/unsuccessful decisions
   - Adjust model confidence scores

2. **Parameter Tuning**
   - Optimize hold time recommendations
   - Calibrate conflict probability thresholds
   - Adjust optimization objective weights

3. **Anomaly Detection**
   - Identify unusual delay patterns
   - Detect degraded vehicle performance
   - Flag data quality issues

### Implementation
```python
class PerformanceAnalyzer:
    def analyze_decision(self, decision, outcome):
        # Calculate success metrics
        delay_reduction = outcome.initial_delay - outcome.final_delay
        conflicts_resolved = len(decision.resolved_conflicts)
        passenger_impact = outcome.total_passenger_minutes_delayed
        
        success_score = (
            0.4 * normalize(delay_reduction) +
            0.3 * conflicts_resolved / max(len(decision.conflicts), 1) +
            0.3 * (1 - normalize(passenger_impact))
        )
        
        # Store for retraining
        self.training_data.append({
            'features': decision.features,
            'action': decision.action,
            'success_score': success_score
        })
        
        # Update running statistics
        self.update_statistics(success_score)
        
        return success_score
```

## Feature Engineering

### Temporal Features
```python
def extract_temporal_features(timestamp):
    return {
        'hour': timestamp.hour,
        'minute': timestamp.minute,
        'day_of_week': timestamp.weekday(),
        'is_weekend': timestamp.weekday() >= 5,
        'is_rush_hour': timestamp.hour in [7, 8, 9, 17, 18, 19],
        'is_holiday': check_holiday(timestamp.date())
    }
```

### Network Features
```python
def extract_network_features(network_state):
    return {
        'active_trains': len(network_state.active_trips),
        'delayed_trains': sum(1 for t in network_state.active_trips.values() if t.current_delay > 0),
        'avg_delay': mean([t.current_delay for t in network_state.active_trips.values()]),
        'track_utilization': len([t for t in network_state.track_occupancy.values() if t.current_occupant]) / len(network_state.track_occupancy),
        'station_congestion': max([s.current_load / s.capacity for s in network_state.station_status.values()])
    }
```

### Weather Features
```python
def extract_weather_features(weather_data):
    return {
        'temperature': weather_data.temperature,
        'precipitation': weather_data.precipitation,
        'wind_speed': weather_data.wind_speed,
        'visibility': weather_data.visibility,
        'is_adverse': weather_data.condition in ['Rain', 'Snow', 'Fog']
    }
```

## Model Deployment

### Serving Models from Go

1. **Option A: gRPC Python Server**
```python
# ml_server.py
import grpc
from concurrent import futures
import ml_service_pb2
import ml_service_pb2_grpc
import joblib

class MLService(ml_service_pb2_grpc.MLServiceServicer):
    def __init__(self):
        self.delay_model = joblib.load('delay_predictor.pkl')
    
    def PredictDelay(self, request, context):
        features = [request.current_delay, request.hour, ...]
        prediction = self.delay_model.predict([features])[0]
        return ml_service_pb2.DelayPrediction(predicted_delay=prediction)

server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
ml_service_pb2_grpc.add_MLServiceServicer_to_server(MLService(), server)
server.add_insecure_port('[::]:50051')
server.start()
```

2. **Option B: REST API**
```python
from flask import Flask, request, jsonify

app = Flask(__name__)
model = joblib.load('delay_predictor.pkl')

@app.route('/predict/delay', methods=['POST'])
def predict_delay():
    data = request.json
    features = extract_features(data)
    prediction = model.predict([features])[0]
    return jsonify({'predicted_delay': float(prediction)})
```

### Go Client
```go
// Call Python ML service
func (dp *DelayPredictor) predictViaAPI(features []float64) (int, error) {
    url := "http://localhost:5000/predict/delay"
    payload := map[string]interface{}{"features": features}
    
    resp, err := http.Post(url, "application/json", marshalJSON(payload))
    if err != nil {
        return 0, err
    }
    defer resp.Body.Close()
    
    var result struct {
        PredictedDelay float64 `json:"predicted_delay"`
    }
    json.NewDecoder(resp.Body).Decode(&result)
    
    return int(result.PredictedDelay), nil
}
```

## Model Retraining

### Trigger Conditions
- Every 24 hours (scheduled)
- After 1000 new decisions collected
- When model accuracy drops below threshold

### Retraining Process
```python
def retrain_models():
    # Load historical data from database
    conn = psycopg2.connect(DATABASE_URL)
    data = pd.read_sql("SELECT * FROM decisions WHERE timestamp > NOW() - INTERVAL '30 days'", conn)
    
    # Train delay predictor
    delay_model = train_delay_predictor(data)
    delay_model.save_model('delay_predictor.pkl')
    
    # Train conflict predictor
    conflict_model = train_conflict_predictor(data)
    joblib.dump(conflict_model, 'conflict_predictor.pkl')
    
    # Evaluate and log
    metrics = evaluate_models(delay_model, conflict_model, test_data)
    log_metrics(metrics)
    
    # Deploy if improved
    if metrics['delay_mae'] < current_metrics['delay_mae']:
        deploy_models()
```

## Monitoring

### Key Metrics
- **Prediction Accuracy**: Compare predictions vs. actual outcomes
- **Response Time**: Model inference latency
- **Recommendation Success Rate**: % of recommendations that improved outcomes

### Logging
```python
import mlflow

mlflow.log_param("model_type", "XGBoost")
mlflow.log_param("n_estimators", 200)
mlflow.log_metric("mae", mae_score)
mlflow.log_metric("rmse", rmse_score)
mlflow.log_artifact("delay_predictor.pkl")
```

## References
- [XGBoost Documentation](https://xgboost.readthedocs.io/)
- [Scikit-learn User Guide](https://scikit-learn.org/stable/user_guide.html)
- [Deep Q-Learning](https://arxiv.org/abs/1312.5602)
