package vehicle

import "time"

// TrainsetType represents different types of trains operating on the NEC
type TrainsetType struct {
	TypeID          string  `json:"type_id"`
	Name            string  `json:"name"` // e.g., "Acela", "ACS-64 + Amfleet", "Siemens Charger"
	Manufacturer    string  `json:"manufacturer"`
	MaxSpeed        float64 `json:"max_speed"` // km/h
	ServiceSpeed    float64 `json:"service_speed"` // typical operational speed
	Length          float64 `json:"length"` // meters
	Capacity        int     `json:"capacity"` // passenger capacity
	PowerType       string  `json:"power_type"` // Electric, Diesel, etc.
}

// PerformanceProfile contains acceleration and braking characteristics
type PerformanceProfile struct {
	TrainsetTypeID  string              `json:"trainset_type_id"`
	AccelerationCurve []AccelerationPoint `json:"acceleration_curve"`
	BrakingCurve    []BrakingPoint      `json:"braking_curve"`
	PowerCurve      []PowerPoint        `json:"power_curve"`
	
	// Environmental factors
	GradeResistance  float64 `json:"grade_resistance"` // coefficient
	CurveResistance  float64 `json:"curve_resistance"` // coefficient
	AerodynamicDrag  float64 `json:"aerodynamic_drag"` // coefficient
}

// AccelerationPoint represents acceleration capability at a given speed
type AccelerationPoint struct {
	Speed        float64 `json:"speed"` // km/h
	Acceleration float64 `json:"acceleration"` // m/s²
}

// BrakingPoint represents braking capability at a given speed
type BrakingPoint struct {
	Speed       float64 `json:"speed"` // km/h
	Deceleration float64 `json:"deceleration"` // m/s²
	BrakingType string  `json:"braking_type"` // service, emergency
}

// PowerPoint represents power consumption/generation at speed
type PowerPoint struct {
	Speed       float64 `json:"speed"` // km/h
	Power       float64 `json:"power"` // kW
	Efficiency  float64 `json:"efficiency"` // 0-1
}

// VehicleMetrics tracks real-time and historical performance of a specific vehicle
type VehicleMetrics struct {
	VehicleID       string    `json:"vehicle_id"`
	TrainsetTypeID  string    `json:"trainset_type_id"`
	LastUpdated     time.Time `json:"last_updated"`
	
	// Real-time metrics
	CurrentSpeed    float64   `json:"current_speed"` // km/h
	CurrentAccel    float64   `json:"current_accel"` // m/s²
	PowerDraw       float64   `json:"power_draw"` // kW
	
	// Performance tracking
	AverageSpeed    float64   `json:"average_speed"` // km/h over last hour
	MaxSpeedAchieved float64  `json:"max_speed_achieved"` // in current trip
	
	// Degradation factors
	PerformanceFactor float64 `json:"performance_factor"` // 0-1, 1 = nominal
	MaintenanceStatus string  `json:"maintenance_status"`
	LastMaintenance time.Time `json:"last_maintenance"`
	
	// Statistical tracking
	TotalDistance   float64   `json:"total_distance"` // km
	TotalTrips      int       `json:"total_trips"`
	AverageDelay    float64   `json:"average_delay"` // seconds per trip
}

// RunningTimeCalculator estimates travel time between points
type RunningTimeCalculator struct {
	profile *PerformanceProfile
}

// TravelTimeEstimate represents predicted travel parameters
type TravelTimeEstimate struct {
	Distance        float64   `json:"distance"` // meters
	EstimatedTime   int       `json:"estimated_time"` // seconds
	AverageSpeed    float64   `json:"average_speed"` // km/h
	MaxSpeedReached float64   `json:"max_speed_reached"` // km/h
	EnergyUsed      float64   `json:"energy_used"` // kWh
	Confidence      float64   `json:"confidence"` // 0-1
}

// NewRunningTimeCalculator creates a calculator for a trainset type
func NewRunningTimeCalculator(profile *PerformanceProfile) *RunningTimeCalculator {
	return &RunningTimeCalculator{
		profile: profile,
	}
}

// CalculateTravelTime estimates time between two points
func (rtc *RunningTimeCalculator) CalculateTravelTime(
	distance float64,
	startSpeed float64,
	endSpeed float64,
	maxAllowedSpeed float64,
	grade float64,
	performanceFactor float64,
) *TravelTimeEstimate {
	// Placeholder for physics-based calculation
	// This would implement the actual dynamics equations
	
	// Simple estimate for now
	avgSpeed := (startSpeed + endSpeed) / 2
	if avgSpeed == 0 {
		avgSpeed = maxAllowedSpeed * 0.7 * performanceFactor
	}
	
	estimatedTime := (distance / 1000.0) / avgSpeed * 3600.0 // seconds
	
	return &TravelTimeEstimate{
		Distance:        distance,
		EstimatedTime:   int(estimatedTime),
		AverageSpeed:    avgSpeed,
		MaxSpeedReached: maxAllowedSpeed * performanceFactor,
		EnergyUsed:      distance * 0.05, // simplified
		Confidence:      0.85,
	}
}
