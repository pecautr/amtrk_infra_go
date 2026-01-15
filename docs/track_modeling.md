# Track Modeling and Signal Block Implementation

This document outlines how to model the Northeast Corridor's track infrastructure in a format suitable for dispatching optimization.

## Track Diagram Elements

### Signal Blocks
The fundamental unit of track occupancy and conflict detection.

```go
type SignalBlock struct {
    BlockID         string    // e.g., "HAROLD-CP1-Block-1"
    ControlPoint    string    // Interlocking or CP name
    StartMilepost   float64   // Beginning milepost
    EndMilepost     float64   // Ending milepost
    Length          float64   // Meters
    Track           string    // Track number/name (e.g., "Track 1", "Main 2")
    Direction       Direction // Northbound, Southbound, Both
    MaxSpeed        float64   // km/h
    SignalType      string    // Absolute, Permissive, etc.
    Electrified     bool
    Capacity        int       // Simultaneous occupancy (usually 1)
    NextBlocks      []string  // Adjacent blocks (track circuits)
    PreviousBlocks  []string
}

type Direction int
const (
    Northbound Direction = iota
    Southbound
    Bidirectional
)
```

### Interlockings (Control Points)
Complex junction areas where tracks converge, diverge, or cross.

```go
type Interlocking struct {
    InterlockingID  string    // e.g., "HAROLD", "DIVIDE", "PRINCE"
    Name            string    // Common name
    Location        string    // Station or landmark
    Milepost        float64
    Latitude        float64
    Longitude       float64
    Tracks          []Track   // All tracks through interlocking
    Routes          []Route   // Possible routing configurations
    Signals         []Signal  // Controlling signals
    OperatingRules  []Rule    // Special restrictions
}

type Route struct {
    RouteID         string
    FromTrack       string
    ToTrack         string
    ThroughTracks   []string  // Intermediate tracks
    ConflictingRoutes []string // Cannot be set simultaneously
    LockingTime     int       // Seconds route stays locked
    MaxSpeed        float64   // km/h through this route
}
```

### Crossovers and Switches
Points where trains can change tracks.

```go
type Crossover struct {
    CrossoverID     string
    Milepost        float64
    FromTrack       string
    ToTrack         string
    Type            string    // "universal", "single", "double"
    MaxSpeed        float64   // Speed restriction through crossover
    ControlledBy    string    // Interlocking ID
    NormalPosition  string    // Default alignment
}
```

### Sidings and Yard Tracks
Storage and passing tracks.

```go
type Siding struct {
    SidingID        string
    Name            string
    StartMilepost   float64
    EndMilepost     float64
    Length          float64   // Meters
    Capacity        int       // Number of trains
    MaxTrainLength  float64   // Meters
    Electrified     bool
    ConnectsTo      []string  // Main line connections
}
```

## Schematic vs. Geographic Modeling

### Schematic Model (for optimization)
- **Purpose**: Operational logic and conflict detection
- **Structure**: Graph-based network
- **Key**: Connectivity and routing rules
- **Not concerned with**: Exact geographic curves

```go
// Example: Track connectivity graph
type TrackNetwork struct {
    Nodes map[string]*NetworkNode // CPs, stations, junctions
    Edges map[string]*TrackSegment // Blocks between nodes
}

type NetworkNode struct {
    NodeID      string
    Type        string    // "interlocking", "station", "junction"
    Milepost    float64
    Connections []string  // Connected node IDs
}

type TrackSegment struct {
    SegmentID   string
    FromNode    string
    ToNode      string
    Blocks      []SignalBlock
    TotalLength float64
    RunningTime int       // Typical travel time (seconds)
}
```

### Geographic Model (for visualization)
- **Purpose**: Display on maps
- **Structure**: GIS LineStrings
- **Key**: Accurate lat/lon positions
- **Used for**: User interfaces, analysis

## Data Integration Workflow

### Phase 1: Schematic Network
1. Identify all control points (interlockings)
2. Map signal blocks between CPs
3. Define possible routes through each CP
4. Establish track connectivity graph

### Phase 2: Operating Rules
1. Speed restrictions by block
2. Route conflict matrices
3. Headway requirements
4. Special restrictions (weather, maintenance)

### Phase 3: Geographic Overlay
1. Map schematic elements to lat/lon
2. Useful for visualization only
3. Not required for optimization logic

## Example: Modeling a Section

### Penn Station Area (Conceptual)

```go
// East River Tunnels approach
blocks := []SignalBlock{
    {
        BlockID:       "HAROLD-WX-1",
        ControlPoint:  "HAROLD",
        StartMilepost: 8.5,
        EndMilepost:   8.8,
        Track:         "Track 1",
        Direction:     Eastbound,
        MaxSpeed:      50.0, // km/h
        Capacity:      1,
    },
    {
        BlockID:       "HAROLD-TUNNEL-1", 
        ControlPoint:  "HAROLD",
        StartMilepost: 8.8,
        EndMilepost:   9.2,
        Track:         "Tunnel 1",
        Direction:     Eastbound,
        MaxSpeed:      40.0,
        Capacity:      1,
    },
}

// Harold Interlocking
harold := Interlocking{
    InterlockingID: "HAROLD",
    Name:          "Harold Interlocking",
    Location:      "Sunnyside Yard",
    Milepost:      8.5,
    Routes: []Route{
        {
            RouteID:   "HAROLD-1-to-TUNNEL1",
            FromTrack: "Track 1",
            ToTrack:   "Tunnel 1",
            ConflictingRoutes: []string{"HAROLD-2-to-TUNNEL1"},
            MaxSpeed:  50.0,
        },
        {
            RouteID:   "HAROLD-2-to-TUNNEL1",
            FromTrack: "Track 2", 
            ToTrack:   "Tunnel 1",
            ConflictingRoutes: []string{"HAROLD-1-to-TUNNEL1"},
            MaxSpeed:  40.0,
        },
    },
}
```

## Building the Model

### Step 1: Data Collection
From available sources (see [data_sources.md](data_sources.md)):
1. Get mile post locations from timetables
2. Extract interlocking names and locations
3. Map track numbers from station diagrams
4. Identify signal block boundaries (if available)

### Step 2: Inference
Where exact data unavailable:
1. **Estimate block lengths**: Typically 1-3 miles on main line
2. **Assume bidirectional**: Unless known single-track
3. **Standard crossover speeds**: 30-50 km/h typical
4. **Use station spacing**: As proxy for interlocking locations

### Step 3: Validation
1. Cross-reference with GTFS stop sequences
2. Verify against employee timetables (if available)
3. Test with real train positions from GTFS-RT
4. Community feedback from rail professionals

## Conflict Detection Using Schematic Model

```go
func (n *TrackNetwork) DetectConflict(trip1, trip2 *Trip) *Conflict {
    // Get next blocks for each trip
    blocks1 := n.GetUpcomingBlocks(trip1, 10) // Next 10 blocks
    blocks2 := n.GetUpcomingBlocks(trip2, 10)
    
    // Check for block conflicts
    for _, b1 := range blocks1 {
        for _, b2 := range blocks2 {
            if b1.BlockID == b2.BlockID {
                // Calculate arrival times
                eta1 := calculateBlockArrival(trip1, b1)
                eta2 := calculateBlockArrival(trip2, b2)
                
                // Check overlap
                if timeRangesOverlap(eta1, eta2, MinimumHeadway) {
                    return &Conflict{
                        Type:     BlockConflict,
                        Location: b1.BlockID,
                        Time:     earlierTime(eta1, eta2),
                        Trips:    []string{trip1.ID, trip2.ID},
                    }
                }
            }
        }
    }
    
    return nil
}
```

## Data Storage Schema

```sql
-- Signal blocks table
CREATE TABLE signal_blocks (
    block_id VARCHAR(50) PRIMARY KEY,
    control_point VARCHAR(50),
    start_milepost DECIMAL(10,3),
    end_milepost DECIMAL(10,3),
    length DECIMAL(10,2),
    track VARCHAR(20),
    direction VARCHAR(20),
    max_speed DECIMAL(6,2),
    signal_type VARCHAR(50),
    electrified BOOLEAN,
    capacity INTEGER DEFAULT 1
);

-- Interlockings table
CREATE TABLE interlockings (
    interlocking_id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100),
    location VARCHAR(100),
    milepost DECIMAL(10,3),
    latitude DECIMAL(10,6),
    longitude DECIMAL(11,6)
);

-- Routes through interlockings
CREATE TABLE interlocking_routes (
    route_id VARCHAR(50) PRIMARY KEY,
    interlocking_id VARCHAR(50) REFERENCES interlockings(interlocking_id),
    from_track VARCHAR(20),
    to_track VARCHAR(20),
    max_speed DECIMAL(6,2),
    locking_time INTEGER
);

-- Route conflicts
CREATE TABLE route_conflicts (
    route_id VARCHAR(50) REFERENCES interlocking_routes(route_id),
    conflicts_with VARCHAR(50) REFERENCES interlocking_routes(route_id),
    PRIMARY KEY (route_id, conflicts_with)
);

-- Block adjacency (track connectivity)
CREATE TABLE block_adjacency (
    from_block VARCHAR(50) REFERENCES signal_blocks(block_id),
    to_block VARCHAR(50) REFERENCES signal_blocks(block_id),
    via_route VARCHAR(50) REFERENCES interlocking_routes(route_id),
    PRIMARY KEY (from_block, to_block)
);
```

## Extracting Data from OpenRailwayMap/OpenStreetMap

While OSM provides geographic (lat/lon) data rather than operational schematics, we can extract and transform it into our signal block model.

### Using Overpass API

Query OSM for railway infrastructure along the NEC:

```python
# Example: Extract NEC signals and tracks between NYP and WAS
import requests
import json

overpass_url = "http://overpass-api.de/api/interpreter"
query = """
[out:json][timeout:60];
(
  way["railway"="rail"]["usage"="main"]["name"~"Northeast Corridor",i]
    (bbox:38.8,-77.1,40.8,-73.9);
  node["railway"="signal"](bbox:38.8,-77.1,40.8,-73.9);
  node["railway"="switch"](bbox:38.8,-77.1,40.8,-73.9);
  node["railway"="buffer_stop"](bbox:38.8,-77.1,40.8,-73.9);
  node["railway"="platform"]["operator"~"Amtrak",i](bbox:38.8,-77.1,40.8,-73.9);
);
out body;
>;
out skel qt;
"""

response = requests.get(overpass_url, params={'data': query})
osm_data = response.json()

# Save for processing
with open('nec_osm_data.json', 'w') as f:
    json.dump(osm_data, f, indent=2)
```

### Transformation Pipeline (Go Implementation)

**Step 1: Parse OSM Track Segments**
```go
type OSMTrackSegment struct {
    WayID       int64
    Nodes       []OSMNode  // Ordered lat/lon points
    MaxSpeed    int        // From "maxspeed" tag
    Electrified bool       // From "electrified" tag
    Tracks      int        // From "tracks" tag (number of parallel tracks)
    Usage       string     // "main", "branch", "industrial"
    Name        string     // Route name
}

type OSMNode struct {
    ID  int64
    Lat float64
    Lon float64
}

func ParseOSMData(data []byte) ([]OSMTrackSegment, []OSMSignal, []OSMSwitch, error) {
    // Parse Overpass JSON response
    // Extract ways, nodes with railway tags
    // Return structured data
}
```

**Step 2: Infer Signal Blocks from Signal Spacing**
```go
func InferSignalBlocks(trackSegments []OSMTrackSegment, signals []OSMSignal) []SignalBlock {
    blocks := []SignalBlock{}
    
    // Sort signals along track by milepost
    sortedSignals := sortSignalsByLocation(signals, trackSegments)
    
    // Create block between each pair of consecutive signals
    for i := 0; i < len(sortedSignals)-1; i++ {
        sig1 := sortedSignals[i]
        sig2 := sortedSignals[i+1]
        
        distance := calculateDistance(sig1.Location, sig2.Location)
        
        block := SignalBlock{
            BlockID:        fmt.Sprintf("BLK-%s-%s", sig1.ID, sig2.ID),
            ControlPoint:   sig1.ControlPoint,
            StartMilepost:  sig1.Milepost,
            EndMilepost:    sig2.Milepost,
            Length:         distance,
            MaxSpeed:       getMaxSpeedBetween(sig1, sig2, trackSegments),
        }
        blocks = append(blocks, block)
    }
    
    return blocks
}

// Haversine distance between two lat/lon points
func calculateDistance(p1, p2 OSMNode) float64 {
    const earthRadius = 6371.0 // km
    
    lat1Rad := p1.Lat * math.Pi / 180
    lat2Rad := p2.Lat * math.Pi / 180
    deltaLat := (p2.Lat - p1.Lat) * math.Pi / 180
    deltaLon := (p2.Lon - p1.Lon) * math.Pi / 180
    
    a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
         math.Cos(lat1Rad)*math.Cos(lat2Rad)*
         math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
    c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
    
    return earthRadius * c
}
```

**Step 3: Build Interlockings from Switch Clusters**
```go
func BuildInterlockings(switches []OSMSwitch, signals []OSMSignal, threshold float64) []Interlocking {
    // Cluster switches within threshold distance (e.g., 500m)
    clusters := clusterSwitchesByProximity(switches, threshold)
    
    interlockings := []Interlocking{}
    for _, cluster := range clusters {
        // Find all signals within cluster area
        clusterSignals := findSignalsNearCluster(cluster, signals, threshold)
        
        interlocking := Interlocking{
            ID:       fmt.Sprintf("INTERLOCKING-%d", len(interlockings)),
            Name:     inferInterlockingName(cluster),
            Location: calculateCentroid(cluster),
            Routes:   inferRoutesFromSwitches(cluster, clusterSignals),
        }
        
        interlockings = append(interlockings, interlocking)
    }
    
    return interlockings
}
```

**Step 4: Calculate Mileposts Along Route**
```go
func CalculateMileposts(wayNodes []OSMNode, originPoint OSMNode) map[int64]float64 {
    mileposts := make(map[int64]float64)
    
    cumulativeDistance := 0.0
    for i, node := range wayNodes {
        mileposts[node.ID] = cumulativeDistance
        
        if i < len(wayNodes)-1 {
            segmentDist := calculateDistance(node, wayNodes[i+1])
            cumulativeDistance += segmentDist
        }
    }
    
    return mileposts
}
```

### OSM Data Quality & Workarounds

**Strengths of OSM Data:**
- ✅ Excellent platform location and length data
- ✅ Good track geometry (geographic accuracy)
- ✅ Switch/turnout locations often tagged
- ✅ Some signal positions marked
- ✅ Electrification and track gauge usually present

**Limitations:**
- ❌ Signal aspects/types rarely complete (often just `railway=signal`)
- ❌ Signal block boundaries NOT explicitly defined
- ❌ Control point names/codes usually missing
- ❌ Interlocking operating rules not captured
- ❌ Track numbering may be inconsistent

**Practical Workarounds:**
1. **Infer block boundaries**: Use typical spacing (1.5-3km for main lines)
2. **Default block capacity**: Assume 1 train per block unless otherwise known
3. **Switch clustering**: Group nearby switches to identify interlocking areas
4. **Validate with USDOT**: Cross-reference with North American Rail Network data
5. **Use Amtrak timetables**: Verify station sequences and rough distances
6. **Progressive refinement**: Start with coarse model, add detail as better data found

### Recommended Data Fusion Approach

Combine multiple sources for best results:

```
┌─────────────────┐     ┌──────────────────┐     ┌─────────────────┐
│ OpenStreetMap   │────▶│ Track Topology   │◀────│ USDOT Rail      │
│ (via Overpass)  │     │ Builder          │     │ Network Dataset │
└─────────────────┘     └──────────────────┘     └─────────────────┘
        │                        │                         │
        │ Platform locations     │ Signal blocks           │ Ownership,
        │ Switch positions       │ Interlockings          │ line names
        │ Track geometry         │ Routes                 │ electrification
        │                        │                         │
        └────────────────────────┼─────────────────────────┘
                                 ▼
                        ┌──────────────────┐
                        │ Unified Track    │
                        │ Model Database   │
                        └──────────────────┘
```

### Implementation Script Template

```go
// cmd/osm_importer/main.go
package main

import (
    "encoding/json"
    "io/ioutil"
    "log"
    
    "github.com/pecautr/amtrk_infra_go/pkg/datamodel"
)

func main() {
    // 1. Download OSM data using Overpass API
    osmData := fetchOSMData("Northeast Corridor", bbox{38.8, -77.1, 40.8, -73.9})
    
    // 2. Parse into structured format
    tracks, signals, switches := parseOSMData(osmData)
    
    // 3. Infer signal blocks
    blocks := inferSignalBlocks(tracks, signals)
    
    // 4. Build interlockings from switch clusters
    interlockings := buildInterlockings(switches, signals, 0.5) // 500m threshold
    
    // 5. Insert into database
    db := connectDatabase()
    insertSignalBlocks(db, blocks)
    insertInterlockings(db, interlockings)
    
    log.Printf("Imported %d signal blocks and %d interlockings", 
        len(blocks), len(interlockings))
}
```

## Implementation Priority

1. **Phase 1** (MVP): 
   - Simple track segments between major stations
   - No interlocking detail
   - Single main line representation

2. **Phase 2** (Enhanced):
   - Add major interlockings (Penn, 30th St, Union Station)
   - Multiple tracks modeled
   - Basic signal blocks

3. **Phase 3** (Full):
   - Complete signal block detail
   - All interlocking routes
   - Conflict matrices
   - Real-time block occupancy

## Resources

- **Railroad Signaling Principles**: Brian Solomon's "Railway Signaling"
- **Interlocking Design**: "The Fundamentals of Railway Signaling and Interlocking"
- **PTC Documentation**: FRA website for modern signal systems
- **Track Charts**: Seek through FOIA or historical railroad societies

## Next Steps

1. Create minimal viable track model for Boston-NYC section
2. Test with real GTFS-RT data
3. Gradually add detail as data becomes available
4. Validate against known train movements
