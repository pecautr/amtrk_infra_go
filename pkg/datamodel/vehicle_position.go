package datamodel

import (
	"time"

	"github.com/google/uuid"
)

// VehiclePosition represents real-time position data from GTFS-RT feeds
type VehiclePosition struct {
	ID          uuid.UUID `json:"id"`
	VehicleID   string    `json:"vehicle_id"`
	TripID      string    `json:"trip_id"`
	RouteID     string    `json:"route_id"`
	Timestamp   time.Time `json:"timestamp"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Speed       float64   `json:"speed"`        // meters per second
	Heading     float64   `json:"heading"`      // degrees
	CurrentStop string    `json:"current_stop"` // Current or next stop ID
	StopStatus  StopStatus `json:"stop_status"`
	Occupancy   OccupancyStatus `json:"occupancy"`
}

type StopStatus int

const (
	IncomingAt StopStatus = iota
	StoppedAt
	InTransitTo
)

type OccupancyStatus int

const (
	Empty OccupancyStatus = iota
	ManySeatsAvailable
	FewSeatsAvailable
	StandingRoomOnly
	CrushedStandingRoomOnly
	Full
	NotAcceptingPassengers
)

// TripUpdate represents schedule deviations and predictions
type TripUpdate struct {
	ID             uuid.UUID `json:"id"`
	TripID         string    `json:"trip_id"`
	RouteID        string    `json:"route_id"`
	VehicleID      string    `json:"vehicle_id"`
	Timestamp      time.Time `json:"timestamp"`
	StopTimeUpdates []StopTimeUpdate `json:"stop_time_updates"`
	Delay          int       `json:"delay"` // seconds
}

// StopTimeUpdate represents arrival/departure predictions for a stop
type StopTimeUpdate struct {
	StopSequence     int       `json:"stop_sequence"`
	StopID           string    `json:"stop_id"`
	ArrivalDelay     int       `json:"arrival_delay"`     // seconds
	DepartureDelay   int       `json:"departure_delay"`   // seconds
	ArrivalTime      time.Time `json:"arrival_time"`
	DepartureTime    time.Time `json:"departure_time"`
	ScheduleRelationship ScheduleRelationship `json:"schedule_relationship"`
}

type ScheduleRelationship int

const (
	Scheduled ScheduleRelationship = iota
	Skipped
	NoData
)

// ServiceAlert represents disruptions and special conditions
type ServiceAlert struct {
	ID          uuid.UUID `json:"id"`
	Timestamp   time.Time `json:"timestamp"`
	Cause       string    `json:"cause"`
	Effect      string    `json:"effect"`
	HeaderText  string    `json:"header_text"`
	Description string    `json:"description"`
	AffectedRoutes []string `json:"affected_routes"`
	AffectedStops  []string `json:"affected_stops"`
	ActivePeriod   []TimeRange `json:"active_period"`
}

type TimeRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}
