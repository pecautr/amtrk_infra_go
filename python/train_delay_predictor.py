"""
Delay Predictor Model Training

Trains an XGBoost model to predict train delays based on:
- Current delay
- Time of day
- Day of week
- Route
- Weather conditions
- Historical patterns

Exports to ONNX format for Rust inference.
"""

import pandas as pd
import numpy as np
from sklearn.model_selection import train_test_split
from sklearn.metrics import mean_absolute_error, mean_squared_error, r2_score
import xgboost as xgb
from skl2onnx import to_onnx
from skl2onnx.common.data_types import FloatTensorType
import joblib
import sys
from pathlib import Path

# Add parent directory to path
sys.path.append(str(Path(__file__).parent))

def load_historical_data():
    """
    Load historical trip delay data from database.
    
    Expected schema:
    - trip_id: str
    - timestamp: datetime
    - current_delay: int (seconds)
    - route_id: str
    - hour_of_day: int (0-23)
    - day_of_week: int (0-6, Monday=0)
    - weather_condition: str
    - actual_future_delay: int (target - delay 30min later)
    """
    
    # TODO: Connect to PostgreSQL and load real data
    # For now, generate synthetic data for demonstration
    
    np.random.seed(42)
    n_samples = 10000
    
    data = pd.DataFrame({
        'current_delay': np.random.randint(-300, 1800, n_samples),  # -5min to +30min
        'hour_of_day': np.random.randint(0, 24, n_samples),
        'day_of_week': np.random.randint(0, 7, n_samples),
        'route_encoded': np.random.randint(0, 10, n_samples),  # Route ID encoded
        'weather_temp_f': np.random.normal(60, 20, n_samples),
        'is_weekend': np.random.randint(0, 2, n_samples),
        'is_rush_hour': np.random.randint(0, 2, n_samples),
        # Target: future delay (with some correlation to current)
        'future_delay': None
    })
    
    # Simulate future delay (correlated with current delay + some noise)
    data['future_delay'] = (
        data['current_delay'] * 1.1 +  # Delays tend to increase slightly
        np.random.normal(0, 180, n_samples) +  # Random noise
        data['is_rush_hour'] * 120  # Rush hour adds delay
    ).astype(int)
    
    return data

def engineer_features(df):
    """Add engineered features"""
    
    # Normalize delays to minutes
    df['current_delay_min'] = df['current_delay'] / 60
    
    # Cyclical time encoding
    df['hour_sin'] = np.sin(2 * np.pi * df['hour_of_day'] / 24)
    df['hour_cos'] = np.cos(2 * np.pi * df['hour_of_day'] / 24)
    df['day_sin'] = np.sin(2 * np.pi * df['day_of_week'] / 7)
    df['day_cos'] = np.cos(2 * np.pi * df['day_of_week'] / 7)
    
    return df

def train_delay_predictor(save_path='../models/delay_predictor.onnx'):
    """Train XGBoost delay prediction model"""
    
    print("Loading historical data...")
    df = load_historical_data()
    df = engineer_features(df)
    
    # Features
    feature_cols = [
        'current_delay', 'hour_of_day', 'day_of_week', 'route_encoded',
        'weather_temp_f', 'is_weekend', 'is_rush_hour',
        'hour_sin', 'hour_cos', 'day_sin', 'day_cos', 'current_delay_min'
    ]
    
    X = df[feature_cols]
    y = df['future_delay']
    
    # Split data
    X_train, X_test, y_train, y_test = train_test_split(
        X, y, test_size=0.2, random_state=42
    )
    
    print(f"Training set: {len(X_train)} samples")
    print(f"Test set: {len(X_test)} samples")
    
    # Train XGBoost model
    print("\nTraining XGBoost model...")
    model = xgb.XGBRegressor(
        n_estimators=200,
        max_depth=6,
        learning_rate=0.1,
        subsample=0.8,
        colsample_bytree=0.8,
        objective='reg:squarederror',
        random_state=42,
        n_jobs=-1
    )
    
    model.fit(
        X_train, y_train,
        eval_set=[(X_test, y_test)],
        early_stopping_rounds=20,
        verbose=False
    )
    
    # Evaluate
    y_pred = model.predict(X_test)
    mae = mean_absolute_error(y_test, y_pred)
    rmse = np.sqrt(mean_squared_error(y_test, y_pred))
    r2 = r2_score(y_test, y_pred)
    
    print(f"\n=== Model Performance ===")
    print(f"MAE:  {mae:.2f} seconds ({mae/60:.2f} minutes)")
    print(f"RMSE: {rmse:.2f} seconds ({rmse/60:.2f} minutes)")
    print(f"R²:   {r2:.4f}")
    
    # Feature importance
    importance = pd.DataFrame({
        'feature': feature_cols,
        'importance': model.feature_importances_
    }).sort_values('importance', ascending=False)
    
    print(f"\n=== Feature Importance ===")
    print(importance.to_string(index=False))
    
    # Export to ONNX
    print(f"\nExporting to ONNX: {save_path}")
    Path(save_path).parent.mkdir(parents=True, exist_ok=True)
    
    initial_type = [('float_input', FloatTensorType([None, len(feature_cols)]))]
    onnx_model = to_onnx(model, X_train[:1].values.astype(np.float32), target_opset=12)
    
    with open(save_path, "wb") as f:
        f.write(onnx_model.SerializeToString())
    
    # Also save as joblib for Python inference
    joblib_path = save_path.replace('.onnx', '.joblib')
    joblib.dump(model, joblib_path)
    
    print(f"✓ ONNX model saved: {save_path}")
    print(f"✓ Joblib model saved: {joblib_path}")
    
    return model

if __name__ == "__main__":
    model = train_delay_predictor()
    print("\n✓ Training complete!")
