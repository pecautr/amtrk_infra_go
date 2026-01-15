# Data Sources and Infrastructure Data

This document catalogs the various data sources used for the Amtrak NEC optimization system.

## Real-Time Operational Data

### Catenary Transit Amtrak GTFS-RT
- **URL**: https://github.com/CatenaryTransit/amtrak-gtfs-rt
- **Feed**: https://gtfs.catenarymaps.org/gtfs-rt/amtrak
- **Data Type**: GTFS-Realtime (Protocol Buffers)
- **Update Frequency**: ~30 seconds
- **Contents**:
  - Vehicle positions (lat/lon, speed, heading)
  - Trip updates (delays, arrival/departure predictions)
  - Service alerts (disruptions, special conditions)
- **Coverage**: All Amtrak services nationwide
- **Limitations**: No historical data (requires local persistence)
- **License**: Open source, community-maintained

## Infrastructure and Network Topology Data

### USDOT North American Rail Network
Federal rail infrastructure datasets from the U.S. Department of Transportation.

#### 1. North American Rail Network Nodes
- **URL**: https://data-usdot.opendata.arcgis.com/datasets/usdot::north-american-rail-network-nodes/about
- **Data Type**: GIS Point Features (Shapefile, GeoJSON, CSV)
- **Contents**:
  - Rail junction points
  - Station locations
  - Yard locations
  - Interchange points
  - Node attributes (type, name, connections)
- **Coverage**: United States, Canada, Mexico
- **Source**: Bureau of Transportation Statistics (BTS)
- **Update Frequency**: Quarterly (last updated November 2025)
- **Use Case**: Network topology, junction modeling

#### 2. North American Rail Lines
- **URL**: https://data-usdot.opendata.arcgis.com/datasets/2d932ad3623640f1affe28656577037e_0/explore
- **Data Type**: GIS Line Features (Shapefile, GeoJSON)
- **Contents**:
  - Rail line segments
  - Track ownership (Class I, Regional, Amtrak, etc.)
  - Number of tracks
  - Operational status
  - Line attributes (electrification, speed limits)
- **Coverage**: United States, Canada, Mexico
- **Source**: Bureau of Transportation Statistics (BTS)
- **Update Frequency**: Annual
- **Use Case**: Route planning, track capacity analysis, infrastructure constraints

### Amtrak Explorer Route Data
- **URL**: https://trriley.github.io/amtrak-explorer/
- **GitHub**: https://github.com/trriley/amtrak-explorer
- **Data Type**: GeoJSON route visualizations
- **Contents**:
  - Amtrak route geometries
  - Station locations
  - Route names and services
- **Coverage**: All Amtrak routes
- **Limitations**: May not include individual track details (mainline vs. siding)
- **Use Case**: Route visualization, general network understanding

## Additional Federal Data Sources

### Federal Railroad Administration (FRA) Data
- **Portal**: https://railroads.dot.gov/
- **Highway-Rail Grade Crossings**: https://safetydata.fra.dot.gov/OfficeofSafety/PublicSite/Crossing/Crossing.aspx
- **Railroad Inspections**: Track condition, safety data
- **Use Case**: Safety constraints, infrastructure condition

### National Transportation Atlas Database (NTAD)
- **URL**: https://geodata.bts.gov/
- **Contents**:
  - Multimodal transportation networks
  - Intermodal facilities
  - Railroad infrastructure
- **Format**: Various GIS formats
- **Use Case**: Comprehensive infrastructure analysis

## NEC-Specific Infrastructure Data

### Track Configuration
For detailed NEC track layouts, consider:

1. **Amtrak Infrastructure Access Manuals** (if publicly available)
2. **State DOT GIS Portals**:
   - Massachusetts: https://www.mass.gov/massgis
   - Connecticut: https://portal.ct.gov/DEEP/GIS/GIS-Data
   - New York: https://gis.ny.gov/
   - New Jersey: https://njogis-newjersey.opendata.arcgis.com/
   - Pennsylvania: https://www.pasda.psu.edu/
   - Maryland: https://data.imap.maryland.gov/
   - Washington DC: https://opendata.dc.gov/

3. **OpenStreetMap (OSM)**:
   - **URL**: https://www.openstreetmap.org/
   - **OpenRailwayMap**: https://www.openrailwaymap.org/ (specialized railway view)
   - **Railway Tags**: https://wiki.openstreetmap.org/wiki/Railways
   - **Extract Tool**: https://overpass-turbo.eu/
   - **Contents**: Detailed track layouts, including:
     - Main tracks
     - Sidings and yard tracks
     - Crossovers and turnouts
     - **Platforms** (with lengths and positions - excellent coverage)
     - Signals (partial coverage - if mapped)
     - Switches (good coverage)
     - Electrification status
   - **Quality**: Varies by location; NEC generally well-mapped
   - **Use Case**: 
     - Platform/station layout modeling (very detailed)
     - Track topology extraction
     - Inferring signal blocks from signal positions
     - Cross-validating other data sources
   - **See**: `docs/track_modeling.md` for transformation pipeline (geographic → operational model)

### Recommended OSM Query for NEC Tracks
```overpass
[out:json][timeout:60];
// Northeast Corridor bounding box
(
  way["railway"="rail"]["usage"="main"](40.0,-75.0,43.0,-70.0);
  way["railway"="rail"]["usage"="branch"](40.0,-75.0,43.0,-70.0);
  node["railway"="station"](40.0,-75.0,43.0,-70.0);
  node["railway"="signal"](40.0,-75.0,43.0,-70.0);
  way["railway"="rail"]["service"="siding"](40.0,-75.0,43.0,-70.0);
);
out geom;
```

## Signal and Control System Data

### Positive Train Control (PTC) Documentation
- **FRA PTC Resources**: https://railroads.dot.gov/train-control/ptc
- May include signal block boundaries and control point locations
- Implementation plans sometimes publicly available

### Interlocking and Control Points
For the NEC's complex interlockings (e.g., "HAROLD", "PRINCE", "DIVIDE"), detailed diagrams may be available through:
- Engineering publications
- Safety investigation reports (NTSB)
- Academic case studies
- Railroad signal engineering conferences

## Historical Performance Data

Since GTFS-RT feeds don't provide historical data, consider:

### 1. Local Data Collection
- **Method**: Persist incoming GTFS-RT feeds to database
- **Storage**: TimescaleDB for time-series efficiency
- **Retention**: Recommend 2+ years for ML training
- **Schema**: See `scripts/schema.sql`

### 2. Alternative Sources
- **Amtrak Track-A-Train Historical Data**: May require official request
- **Transit App APIs**: Some third-party apps archive historical data
- **Academic Datasets**: Search for published transit research datasets

### 3. Freedom of Information Act (FOIA)
- Request historical operational data from Amtrak
- May include:
  - On-time performance records
  - Delay causes and durations
  - Dispatch decisions
  - Incident reports

## Data Integration Plan

### Phase 1: Real-Time Operations (Current)
- ✅ GTFS-RT ingestion from Catenary Transit
- ✅ Database persistence for historical analysis
- 🔄 Real-time network state tracking

### Phase 2: Infrastructure Integration
- [ ] Download USDOT rail nodes and lines
- [ ] Import to PostGIS database
- [ ] Match GTFS stops to infrastructure nodes
- [ ] Build track segment mappings

### Phase 3: Detailed Track Modeling
- [ ] Extract OSM railway data for NEC
- [ ] Identify main tracks, sidings, crossovers
- [ ] Model platform assignments
- [ ] Add speed restrictions from FRA data

### Phase 4: Enhanced Constraints
- [ ] Grade crossing locations (FRA)
- [ ] Bridge and tunnel restrictions
- [ ] Maintenance windows
- [ ] Crew base locations

## Data Quality Considerations

### Known Limitations
1. **USDOT Rail Lines**: May not include all sidings and yard tracks
2. **OpenStreetMap**: Coverage varies; requires validation
3. **GTFS-RT**: No historical data; requires collection
4. **Track Numbers**: Platform/track assignments may not be in open data

### Validation Strategy
1. Cross-reference multiple sources
2. Field validation where possible
3. Use official Amtrak timetables for verification
4. Community contributions and corrections

## License and Usage

- **Federal Data (USDOT, FRA)**: Public domain, no restrictions
- **OpenStreetMap**: ODbL license, attribution required
- **GTFS-RT (Catenary Transit)**: Open source, check repository for specifics
- **Amtrak Explorer**: Check repository license

## Contributing Data

If you have access to additional data sources:
1. Document in this file
2. Add processing scripts to `scripts/data_processing/`
3. Update database schema if needed
4. Submit pull request

## References

- Bureau of Transportation Statistics: https://www.bts.gov/
- Federal Railroad Administration: https://railroads.dot.gov/
- OpenStreetMap Railway Mapping: https://wiki.openstreetmap.org/wiki/Railways
- GTFS-Realtime Specification: https://gtfs.org/realtime/
