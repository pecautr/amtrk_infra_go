package datamodel

import "time"

// NetworkState represents the current state of the entire NEC network
type NetworkState struct {
	Timestamp       time.Time         `json:"timestamp"`
	ActiveTrips     map[string]*Trip  `json:"active_trips"`
	TrackOccupancy  map[string]*TrackSegment `json:"track_occupancy"`
	StationStatus   map[string]*StationState `json:"station_status"`
	Conflicts       []Conflict        `json:"conflicts"`
	Recommendations []Recommendation  `json:"recommendations"`
}

// Trip represents a complete train journey
type Trip struct {
	TripID          string           `json:"trip_id"`
	RouteID         string           `json:"route_id"`
	VehicleID       string           `json:"vehicle_id"`
	TrainsetType    string           `json:"trainset_type"` // Acela, Regional, etc.
	Direction       string           `json:"direction"`
	Status          TripStatus       `json:"status"`
	CurrentPosition *VehiclePosition `json:"current_position"`
	ScheduledStops  []ScheduledStop  `json:"scheduled_stops"`
	ActualStops     []ActualStop     `json:"actual_stops"`
	CurrentDelay    int              `json:"current_delay"` // seconds
	PredictedDelay  int              `json:"predicted_delay"` // seconds, ML prediction
}

type TripStatus int

const (
	TripScheduled TripStatus = iota
	TripActive
	TripDelayed
	TripCancelled
	TripCompleted
)

// ScheduledStop represents the planned stop at a station
type ScheduledStop struct {
	StopSequence      int       `json:"stop_sequence"`
	StopID            string    `json:"stop_id"`
	StopName          string    `json:"stop_name"`
	ScheduledArrival  time.Time `json:"scheduled_arrival"`
	ScheduledDeparture time.Time `json:"scheduled_departure"`
	Track             string    `json:"track"`
	PlatformID        string    `json:"platform_id"`
}

// ActualStop represents what actually happened/is happening at a stop
type ActualStop struct {
	StopSequence    int       `json:"stop_sequence"`
	StopID          string    `json:"stop_id"`
	ActualArrival   time.Time `json:"actual_arrival"`
	ActualDeparture time.Time `json:"actual_departure"`
	Track           string    `json:"track"`
	DwellTime       int       `json:"dwell_time"` // seconds
}

// TrackSegment represents a section of track between two points
type TrackSegment struct {
	SegmentID      string    `json:"segment_id"`
	StartLocation  string    `json:"start_location"`
	EndLocation    string    `json:"end_location"`
	Length         float64   `json:"length"` // meters
	MaxSpeed       float64   `json:"max_speed"` // km/h
	CurrentOccupant string   `json:"current_occupant"` // TripID or empty
	Capacity       int       `json:"capacity"` // max trains
	TrackType      string    `json:"track_type"` // main, siding, etc.
}

// StationState represents current conditions at a station
type StationState struct {
	StationID       string              `json:"station_id"`
	StationName     string              `json:"station_name"`
	PlatformStates  map[string]*Platform `json:"platform_states"`
	Capacity        int                 `json:"capacity"` // simultaneous trains
	CurrentLoad     int                 `json:"current_load"`
}

// Platform represents a single platform at a station
type Platform struct {
	PlatformID     string    `json:"platform_id"`
	Track          string    `json:"track"`
	OccupiedBy     string    `json:"occupied_by"` // TripID or empty
	OccupiedUntil  time.Time `json:"occupied_until"`
	AvailableFrom  time.Time `json:"available_from"`
}

// Conflict represents a predicted or actual conflict in the network
type Conflict struct {
	ConflictID      string    `json:"conflict_id"`
	DetectedAt      time.Time `json:"detected_at"`
	ConflictTime    time.Time `json:"conflict_time"`
	ConflictType    ConflictType `json:"conflict_type"`
	InvolvedTrips   []string  `json:"involved_trips"`
	Location        string    `json:"location"` // Track segment or station
	Severity        int       `json:"severity"` // 1-10
	EstimatedImpact int       `json:"estimated_impact"` // delay in seconds
	Resolution      string    `json:"resolution"` // empty if unresolved
}

type ConflictType int

const (
	TrackConflict ConflictType = iota
	PlatformConflict
	CrossingConflict
	MaintenanceConflict
	SpeedRestriction
)

// Recommendation represents a dispatching decision
type Recommendation struct {
	RecommendationID string    `json:"recommendation_id"`
	GeneratedAt      time.Time `json:"generated_at"`
	ActionType       ActionType `json:"action_type"`
	TripID           string    `json:"trip_id"`
	StopID           string    `json:"stop_id,omitempty"`
	HoldDuration     int       `json:"hold_duration,omitempty"` // seconds
	AlternativeTrack string    `json:"alternative_track,omitempty"`
	Priority         int       `json:"priority"` // 1-10
	ExpectedBenefit  int       `json:"expected_benefit"` // delay reduction in seconds
	Confidence       float64   `json:"confidence"` // 0-1
}

type ActionType int

const (
	HoldAtStation ActionType = iota
	ExpediteDispatch
	ChangeTrack
	RerouteVia
	CancelStop
	AdjustSpeed
)
