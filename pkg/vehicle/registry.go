package vehicle

import (
	"encoding/json"
	"os"
)

// TrainsetRegistry manages all known trainset types
type TrainsetRegistry struct {
	trainsets map[string]*TrainsetType
	profiles  map[string]*PerformanceProfile
}

// NewTrainsetRegistry creates a new registry
func NewTrainsetRegistry() *TrainsetRegistry {
	return &TrainsetRegistry{
		trainsets: make(map[string]*TrainsetType),
		profiles:  make(map[string]*PerformanceProfile),
	}
}

// LoadFromFile loads trainset definitions from JSON
func (tr *TrainsetRegistry) LoadFromFile(filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	var trainsets []TrainsetType
	if err := json.Unmarshal(data, &trainsets); err != nil {
		return err
	}

	for i := range trainsets {
		tr.trainsets[trainsets[i].TypeID] = &trainsets[i]
	}

	return nil
}

// GetTrainset retrieves trainset by type ID
func (tr *TrainsetRegistry) GetTrainset(typeID string) *TrainsetType {
	return tr.trainsets[typeID]
}

// GetPerformanceProfile retrieves performance profile for a trainset type
func (tr *TrainsetRegistry) GetPerformanceProfile(typeID string) *PerformanceProfile {
	return tr.profiles[typeID]
}

// RegisterTrainset adds a new trainset type
func (tr *TrainsetRegistry) RegisterTrainset(trainset *TrainsetType, profile *PerformanceProfile) {
	tr.trainsets[trainset.TypeID] = trainset
	if profile != nil {
		tr.profiles[trainset.TypeID] = profile
	}
}

// GetDefaultNECTrainsets returns typical NEC trainset configurations
func GetDefaultNECTrainsets() []*TrainsetType {
	return []*TrainsetType{
		{
			TypeID:       "acela",
			Name:         "Acela Express",
			Manufacturer: "Alstom",
			MaxSpeed:     240.0, // km/h (150 mph)
			ServiceSpeed: 215.0, // typical
			Length:       201.5, // meters
			Capacity:     304,
			PowerType:    "Electric",
		},
		{
			TypeID:       "acs64-amfleet",
			Name:         "ACS-64 with Amfleet",
			Manufacturer: "Siemens",
			MaxSpeed:     200.0, // km/h (125 mph)
			ServiceSpeed: 177.0,
			Length:       180.0,
			Capacity:     450,
			PowerType:    "Electric",
		},
		{
			TypeID:       "acs64-viewliner",
			Name:         "ACS-64 with Viewliner",
			Manufacturer: "Siemens",
			MaxSpeed:     200.0,
			ServiceSpeed: 177.0,
			Length:       200.0,
			Capacity:     250,
			PowerType:    "Electric",
		},
	}
}

// GetDefaultPerformanceProfiles returns default profiles for NEC trains
func GetDefaultPerformanceProfiles() map[string]*PerformanceProfile {
	return map[string]*PerformanceProfile{
		"acela": {
			TrainsetTypeID: "acela",
			AccelerationCurve: []AccelerationPoint{
				{Speed: 0, Acceleration: 1.2},
				{Speed: 50, Acceleration: 1.0},
				{Speed: 100, Acceleration: 0.7},
				{Speed: 150, Acceleration: 0.4},
				{Speed: 200, Acceleration: 0.2},
				{Speed: 240, Acceleration: 0.0},
			},
			BrakingCurve: []BrakingPoint{
				{Speed: 240, Deceleration: 0.8, BrakingType: "service"},
				{Speed: 150, Deceleration: 1.0, BrakingType: "service"},
				{Speed: 100, Deceleration: 1.2, BrakingType: "service"},
				{Speed: 50, Deceleration: 1.3, BrakingType: "service"},
				{Speed: 0, Deceleration: 0.0, BrakingType: "service"},
			},
			GradeResistance: 0.02,
			CurveResistance: 0.015,
			AerodynamicDrag: 0.012,
		},
		"acs64-amfleet": {
			TrainsetTypeID: "acs64-amfleet",
			AccelerationCurve: []AccelerationPoint{
				{Speed: 0, Acceleration: 1.0},
				{Speed: 50, Acceleration: 0.8},
				{Speed: 100, Acceleration: 0.5},
				{Speed: 150, Acceleration: 0.3},
				{Speed: 200, Acceleration: 0.0},
			},
			BrakingCurve: []BrakingPoint{
				{Speed: 200, Deceleration: 0.7, BrakingType: "service"},
				{Speed: 100, Deceleration: 1.0, BrakingType: "service"},
				{Speed: 50, Deceleration: 1.2, BrakingType: "service"},
				{Speed: 0, Deceleration: 0.0, BrakingType: "service"},
			},
			GradeResistance: 0.025,
			CurveResistance: 0.018,
			AerodynamicDrag: 0.015,
		},
	}
}
