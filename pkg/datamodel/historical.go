package datamodel

import "time"

// HistoricalDispatchDecision stores past dispatching decisions and outcomes
type HistoricalDispatchDecision struct {
	ID              string    `json:"id"`
	Timestamp       time.Time `json:"timestamp"`
	NetworkState    string    `json:"network_state"` // JSON snapshot
	Decision        string    `json:"decision"` // Action taken
	TripID          string    `json:"trip_id"`
	ConflictContext string    `json:"conflict_context"`
	
	// Outcomes (measured after decision)
	ActualDelay     int       `json:"actual_delay"` // seconds
	NetworkImpact   int       `json:"network_impact"` // total delay caused
	SuccessScore    float64   `json:"success_score"` // 0-1
	
	// Context features for ML
	TimeOfDay       int       `json:"time_of_day"` // hour
	DayOfWeek       int       `json:"day_of_week"`
	TrainType       string    `json:"train_type"`
	CurrentLoad     int       `json:"current_load"` // trains in network
	WeatherCondition string   `json:"weather_condition"`
}

// PerformanceMetrics tracks aggregated system performance
type PerformanceMetrics struct {
	PeriodStart     time.Time `json:"period_start"`
	PeriodEnd       time.Time `json:"period_end"`
	
	TotalTrips      int       `json:"total_trips"`
	OnTimeTrips     int       `json:"on_time_trips"`
	DelayedTrips    int       `json:"delayed_trips"`
	CancelledTrips  int       `json:"cancelled_trips"`
	
	AverageDelay    float64   `json:"average_delay"` // seconds
	MaxDelay        int       `json:"max_delay"` // seconds
	TotalDelayMinutes int     `json:"total_delay_minutes"`
	
	ConflictsDetected int     `json:"conflicts_detected"`
	ConflictsResolved int     `json:"conflicts_resolved"`
	
	OTPPercentage   float64   `json:"otp_percentage"` // On-Time Performance
}
