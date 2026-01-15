package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Command-line flags
var (
	section     = flag.String("section", "nyp-phil", "Section to import: nyp-phil, phil-wash, bos-nyp, full-nec")
	outputFile  = flag.String("output", "osm_data.json", "Output file for raw OSM data")
	sqlFile     = flag.String("sql", "import_blocks.sql", "Output SQL file for signal blocks")
	dryRun      = flag.Bool("dry-run", false, "Download and process but don't generate SQL")
	verbose     = flag.Bool("verbose", false, "Verbose logging")
)

// BoundingBox represents a geographic area
type BoundingBox struct {
	MinLat float64
	MinLon float64
	MaxLat float64
	MaxLon float64
}

// Section definitions for NEC
var sections = map[string]BoundingBox{
	"nyp-phil": {
		MinLat: 39.95,  // Philadelphia 30th Street
		MinLon: -75.18,
		MaxLat: 40.75,  // Penn Station NYC
		MaxLon: -74.0,
	},
	"phil-wash": {
		MinLat: 38.89,  // Union Station DC
		MinLon: -77.01,
		MaxLat: 39.95,  // Philadelphia 30th Street
		MaxLon: -75.15,
	},
	"bos-nyp": {
		MinLat: 40.75,  // Penn Station NYC
		MinLon: -74.0,
		MaxLat: 42.37,  // Boston South Station
		MaxLon: -71.05,
	},
	"full-nec": {
		MinLat: 38.89,  // DC
		MinLon: -77.01,
		MaxLat: 42.37,  // Boston
		MaxLon: -71.05,
	},
}

// OSM data structures
type OSMResponse struct {
	Elements []OSMElement `json:"elements"`
}

type OSMElement struct {
	Type string             `json:"type"`
	ID   int64              `json:"id"`
	Lat  float64            `json:"lat,omitempty"`
	Lon  float64            `json:"lon,omitempty"`
	Nodes []int64           `json:"nodes,omitempty"`
	Tags map[string]string  `json:"tags,omitempty"`
}

type OSMNode struct {
	ID  int64
	Lat float64
	Lon float64
}

type OSMWay struct {
	ID      int64
	Nodes   []OSMNode
	Tags    map[string]string
}

type OSMSignal struct {
	ID       int64
	Location OSMNode
	SignalType string
}

type OSMSwitch struct {
	ID       int64
	Location OSMNode
	SwitchType string
}

type OSMPlatform struct {
	ID       int64
	Location OSMNode
	Name     string
	Ref      string
	Length   float64
}

// SignalBlock represents an inferred block
type SignalBlock struct {
	BlockID       string
	ControlPoint  string
	StartMilepost float64
	EndMilepost   float64
	Length        float64
	Track         string
	Direction     string
	MaxSpeed      float64
	StartLat      float64
	StartLon      float64
	EndLat        float64
	EndLon        float64
}

func main() {
	flag.Parse()

	bbox, ok := sections[*section]
	if !ok {
		log.Fatalf("Unknown section: %s. Available: nyp-phil, phil-wash, bos-nyp, full-nec", *section)
	}

	log.Printf("Importing Northeast Corridor section: %s", *section)
	log.Printf("Bounding box: (%.2f,%.2f) to (%.2f,%.2f)", bbox.MinLat, bbox.MinLon, bbox.MaxLat, bbox.MaxLon)

	// Build Overpass query
	query := buildOverpassQuery(bbox)
	
	// Fetch data from Overpass API
	log.Println("Fetching data from Overpass API...")
	osmData, err := fetchOverpassData(query)
	if err != nil {
		log.Fatalf("Failed to fetch OSM data: %v", err)
	}

	// Save raw data
	if err := saveRawData(osmData, *outputFile); err != nil {
		log.Fatalf("Failed to save raw data: %v", err)
	}
	log.Printf("Saved raw OSM data to %s (%d elements)", *outputFile, len(osmData.Elements))

	// Parse into structured format
	ways, signals, switches, platforms := parseOSMData(osmData)
	log.Printf("Parsed: %d track segments, %d signals, %d switches, %d platforms", 
		len(ways), len(signals), len(switches), len(platforms))

	// Filter for main line tracks only
	mainTracks := filterMainLineTracks(ways)
	log.Printf("Filtered to %d main line track segments", len(mainTracks))

	// Infer signal blocks
	blocks := inferSignalBlocks(mainTracks, signals, bbox)
	log.Printf("Inferred %d signal blocks", len(blocks))

	if *verbose {
		for i, block := range blocks {
			if i < 5 { // Show first 5
				log.Printf("  Block %s: MP %.2f-%.2f (%.2f km, max %d km/h)", 
					block.BlockID, block.StartMilepost, block.EndMilepost, 
					block.Length, int(block.MaxSpeed))
			}
		}
	}

	// Generate SQL
	if !*dryRun {
		if err := generateSQL(blocks, platforms, *sqlFile); err != nil {
			log.Fatalf("Failed to generate SQL: %v", err)
		}
		log.Printf("Generated SQL import script: %s", *sqlFile)
	}

	log.Println("Import complete!")
	log.Printf("\nNext steps:")
	log.Printf("  1. Review the generated SQL: %s", *sqlFile)
	log.Printf("  2. Run: psql -d amtrk_infra -f %s", *sqlFile)
	log.Printf("  3. Verify: SELECT COUNT(*) FROM signal_blocks;")
}

func buildOverpassQuery(bbox BoundingBox) string {
	return fmt.Sprintf(`
[out:json][timeout:120];
(
  // Main rail tracks (passenger service)
  way["railway"="rail"]["usage"="main"]["service"!="yard"]["service"!="siding"]
    (%.6f,%.6f,%.6f,%.6f);
  
  // High-speed rail tracks
  way["railway"="rail"]["highspeed"="yes"]
    (%.6f,%.6f,%.6f,%.6f);
  
  // Signals
  node["railway"="signal"]
    (%.6f,%.6f,%.6f,%.6f);
  
  // Switches/turnouts
  node["railway"="switch"]
    (%.6f,%.6f,%.6f,%.6f);
  
  // Platforms (Amtrak stations)
  node["railway"="platform"]["operator"~"Amtrak",i]
    (%.6f,%.6f,%.6f,%.6f);
  way["railway"="platform"]["operator"~"Amtrak",i]
    (%.6f,%.6f,%.6f,%.6f);
  
  // Station nodes
  node["railway"="station"]["operator"~"Amtrak",i]
    (%.6f,%.6f,%.6f,%.6f);
);
out body;
>;
out skel qt;
`,
		bbox.MinLat, bbox.MinLon, bbox.MaxLat, bbox.MaxLon, // rail tracks
		bbox.MinLat, bbox.MinLon, bbox.MaxLat, bbox.MaxLon, // high-speed
		bbox.MinLat, bbox.MinLon, bbox.MaxLat, bbox.MaxLon, // signals
		bbox.MinLat, bbox.MinLon, bbox.MaxLat, bbox.MaxLon, // switches
		bbox.MinLat, bbox.MinLon, bbox.MaxLat, bbox.MaxLon, // platform nodes
		bbox.MinLat, bbox.MinLon, bbox.MaxLat, bbox.MaxLon, // platform ways
		bbox.MinLat, bbox.MinLon, bbox.MaxLat, bbox.MaxLon, // stations
	)
}

func fetchOverpassData(query string) (*OSMResponse, error) {
	url := "https://overpass-api.de/api/interpreter"
	
	resp, err := http.Post(url, "application/x-www-form-urlencoded", 
		strings.NewReader("data="+query))
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Overpass API returned %d: %s", resp.StatusCode, string(body))
	}

	var osmData OSMResponse
	if err := json.NewDecoder(resp.Body).Decode(&osmData); err != nil {
		return nil, fmt.Errorf("JSON decode failed: %w", err)
	}

	return &osmData, nil
}

func saveRawData(data *OSMResponse, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func parseOSMData(data *OSMResponse) ([]OSMWay, []OSMSignal, []OSMSwitch, []OSMPlatform) {
	// First pass: index all nodes by ID
	nodeIndex := make(map[int64]OSMNode)
	for _, elem := range data.Elements {
		if elem.Type == "node" {
			nodeIndex[elem.ID] = OSMNode{
				ID:  elem.ID,
				Lat: elem.Lat,
				Lon: elem.Lon,
			}
		}
	}

	var ways []OSMWay
	var signals []OSMSignal
	var switches []OSMSwitch
	var platforms []OSMPlatform

	// Second pass: process ways and tagged nodes
	for _, elem := range data.Elements {
		switch elem.Type {
		case "node":
			if elem.Tags["railway"] == "signal" {
				signals = append(signals, OSMSignal{
					ID:       elem.ID,
					Location: nodeIndex[elem.ID],
					SignalType: elem.Tags["railway:signal:main"],
				})
			} else if elem.Tags["railway"] == "switch" {
				switches = append(switches, OSMSwitch{
					ID:       elem.ID,
					Location: nodeIndex[elem.ID],
					SwitchType: elem.Tags["railway:switch"],
				})
			} else if elem.Tags["railway"] == "platform" {
				platforms = append(platforms, OSMPlatform{
					ID:       elem.ID,
					Location: nodeIndex[elem.ID],
					Name:     elem.Tags["name"],
					Ref:      elem.Tags["ref"],
				})
			}

		case "way":
			if elem.Tags["railway"] == "rail" {
				// Resolve node IDs to coordinates
				var nodes []OSMNode
				for _, nodeID := range elem.Nodes {
					if node, ok := nodeIndex[nodeID]; ok {
						nodes = append(nodes, node)
					}
				}
				
				if len(nodes) > 0 {
					ways = append(ways, OSMWay{
						ID:    elem.ID,
						Nodes: nodes,
						Tags:  elem.Tags,
					})
				}
			}
		}
	}

	return ways, signals, switches, platforms
}

func filterMainLineTracks(ways []OSMWay) []OSMWay {
	var mainTracks []OSMWay
	for _, way := range ways {
		// Include if:
		// - usage=main OR highspeed=yes
		// - NOT a yard or siding (unless explicitly main)
		usage := way.Tags["usage"]
		service := way.Tags["service"]
		highspeed := way.Tags["highspeed"]
		
		if usage == "main" || highspeed == "yes" {
			if service != "yard" && service != "siding" {
				mainTracks = append(mainTracks, way)
			}
		}
	}
	return mainTracks
}

func inferSignalBlocks(ways []OSMWay, signals []OSMSignal, bbox BoundingBox) []SignalBlock {
	blocks := []SignalBlock{}

	// Use Penn Station NYC as reference point (milepost 0)
	origin := OSMNode{Lat: 40.7505, Lon: -73.9934}

	for _, way := range ways {
		if len(way.Nodes) < 2 {
			continue
		}

		// Get max speed from tags (default 125 mph for NEC)
		maxSpeed := 201.0 // 125 mph in km/h
		if speedStr, ok := way.Tags["maxspeed"]; ok {
			if strings.HasSuffix(speedStr, "mph") {
				mph := 0.0
				fmt.Sscanf(speedStr, "%f", &mph)
				maxSpeed = mph * 1.60934
			} else {
				fmt.Sscanf(speedStr, "%f", &maxSpeed)
			}
		}

		// Get track number/direction if available
		track := way.Tags["railway:track_ref"]
		if track == "" {
			track = "Main"
		}

		// Create blocks between consecutive waypoints
		// In absence of actual signal data, use ~2km segments
		const targetBlockLength = 2.0 // km

		for i := 0; i < len(way.Nodes)-1; i++ {
			node1 := way.Nodes[i]
			node2 := way.Nodes[i+1]

			distance := haversineDistance(node1, node2)
			
			// Calculate mileposts relative to origin
			mp1 := haversineDistance(origin, node1)
			mp2 := haversineDistance(origin, node2)

			// Determine direction based on milepost increase
			direction := "NB" // Northbound (toward Boston)
			if mp2 < mp1 {
				direction = "SB" // Southbound (toward DC)
				mp1, mp2 = mp2, mp1 // Swap for consistent ordering
			}

			blockID := fmt.Sprintf("BLK-%d-%d", way.ID, i)

			block := SignalBlock{
				BlockID:       blockID,
				ControlPoint:  "", // To be filled in later
				StartMilepost: mp1,
				EndMilepost:   mp2,
				Length:        distance,
				Track:         track,
				Direction:     direction,
				MaxSpeed:      maxSpeed,
				StartLat:      node1.Lat,
				StartLon:      node1.Lon,
				EndLat:        node2.Lat,
				EndLon:        node2.Lon,
			}

			blocks = append(blocks, block)
		}
	}

	// Sort blocks by milepost
	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i].StartMilepost < blocks[j].StartMilepost
	})

	// Assign control point names based on proximity to known stations
	assignControlPoints(blocks)

	return blocks
}

func assignControlPoints(blocks []SignalBlock) {
	// Major NEC stations as control points
	controlPoints := map[string]OSMNode{
		"BOS": {Lat: 42.3519, Lon: -71.0552},  // Boston South
		"BBY": {Lat: 42.3467, Lon: -71.0763},  // Back Bay
		"RTE": {Lat: 41.8231, Lon: -71.4128},  // Route 128
		"PVD": {Lat: 41.8296, Lon: -71.4128},  // Providence
		"NYP": {Lat: 40.7505, Lon: -73.9934},  // Penn Station NYC
		"NWK": {Lat: 40.7357, Lon: -74.1644},  // Newark Penn
		"TRE": {Lat: 40.2170, Lon: -74.7544},  // Trenton
		"PHL": {Lat: 39.9566, Lon: -75.1822},  // 30th Street Philadelphia
		"WIL": {Lat: 39.7368, Lon: -75.5521},  // Wilmington
		"BAL": {Lat: 39.3074, Lon: -76.6157},  // Baltimore Penn
		"BWI": {Lat: 39.1753, Lon: -76.6683},  // BWI Airport
		"WAS": {Lat: 38.8977, Lon: -77.0063},  // Union Station DC
	}

	for i := range blocks {
		minDist := 999999.0
		closestCP := ""

		midpoint := OSMNode{
			Lat: (blocks[i].StartLat + blocks[i].EndLat) / 2,
			Lon: (blocks[i].StartLon + blocks[i].EndLon) / 2,
		}

		for cpName, cpLoc := range controlPoints {
			dist := haversineDistance(midpoint, cpLoc)
			if dist < minDist {
				minDist = dist
				closestCP = cpName
			}
		}

		// Only assign if within 10km of station
		if minDist < 10.0 {
			blocks[i].ControlPoint = closestCP
		}
	}
}

func haversineDistance(p1, p2 OSMNode) float64 {
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

func generateSQL(blocks []SignalBlock, platforms []OSMPlatform, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Header
	fmt.Fprintf(file, "-- Signal Blocks Import for Northeast Corridor\n")
	fmt.Fprintf(file, "-- Generated: %s\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(file, "-- Source: OpenStreetMap via Overpass API\n\n")

	fmt.Fprintf(file, "BEGIN;\n\n")

	// Delete existing data (for re-imports)
	fmt.Fprintf(file, "-- Clear existing demo data\n")
	fmt.Fprintf(file, "DELETE FROM signal_blocks WHERE block_id LIKE 'BLK-%%';\n\n")

	// Insert signal blocks
	fmt.Fprintf(file, "-- Insert signal blocks\n")
	for _, block := range blocks {
		fmt.Fprintf(file, "INSERT INTO signal_blocks (block_id, control_point, start_milepost, end_milepost, length, track, direction, max_speed, electrified, capacity) VALUES\n")
		fmt.Fprintf(file, "  ('%s', %s, %.3f, %.3f, %.2f, '%s', '%s', %.1f, true, 1);\n",
			block.BlockID,
			sqlString(block.ControlPoint),
			block.StartMilepost,
			block.EndMilepost,
			block.Length,
			block.Track,
			block.Direction,
			block.MaxSpeed,
		)
	}

	fmt.Fprintf(file, "\nCOMMIT;\n\n")

	// Summary
	fmt.Fprintf(file, "-- Summary\n")
	fmt.Fprintf(file, "-- Imported %d signal blocks\n", len(blocks))
	fmt.Fprintf(file, "-- Milepost range: %.1f to %.1f\n", 
		blocks[0].StartMilepost, blocks[len(blocks)-1].EndMilepost)

	return nil
}

func sqlString(s string) string {
	if s == "" {
		return "NULL"
	}
	return fmt.Sprintf("'%s'", strings.ReplaceAll(s, "'", "''"))
}
