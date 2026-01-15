# OSM Importer - Northeast Corridor Track Data

This tool imports real-world track data from OpenStreetMap for the Northeast Corridor.

## Quick Start

```bash
# Build the importer
cd cmd/osm_importer
go build

# Import NYC to Philadelphia section (demo)
./osm_importer -section nyp-phil

# Review the generated SQL
cat import_blocks.sql

# Import into database (when ready)
psql -d amtrk_infra -f import_blocks.sql
```

## Available Sections

- `nyp-phil` - Penn Station NYC to Philadelphia 30th St (default, ~90 miles)
- `phil-wash` - Philadelphia to Washington DC (~135 miles)
- `bos-nyp` - Boston to Penn Station NYC (~230 miles)
- `full-nec` - Complete Boston to DC corridor (~450 miles)

## Options

```bash
./osm_importer [flags]

Flags:
  -section string
        Section to import (default "nyp-phil")
        Options: nyp-phil, phil-wash, bos-nyp, full-nec
        
  -output string
        Output file for raw OSM data (default "osm_data.json")
        
  -sql string
        Output SQL file for signal blocks (default "import_blocks.sql")
        
  -dry-run
        Download and process but don't generate SQL
        
  -verbose
        Verbose logging
```

## What It Does

1. **Fetches OSM Data**: Queries Overpass API for NEC railway infrastructure:
   - Main rail tracks (passenger service)
   - High-speed rail tracks
   - Signals (when available)
   - Switches/turnouts
   - Platforms (Amtrak stations)

2. **Processes Data**:
   - Filters to main line tracks only (excludes yards, sidings)
   - Calculates mileposts from Penn Station NYC (MP 0.0)
   - Infers signal blocks from track segments (~2km each)
   - Assigns control points based on proximity to major stations
   - Determines track direction (NB/SB)

3. **Generates SQL**:
   - Creates INSERT statements for `signal_blocks` table
   - Includes milepost, length, max speed, direction
   - Ready to load into PostgreSQL database

## Example Output

```sql
INSERT INTO signal_blocks (block_id, control_point, start_milepost, end_milepost, length, track, direction, max_speed, electrified, capacity) VALUES
  ('BLK-123456-0', 'NYP', 0.000, 2.134, 2.13, 'Main', 'SB', 201.0, true, 1);
  ('BLK-123456-1', 'NYP', 2.134, 4.289, 2.16, 'Main', 'SB', 201.0, true, 1);
  ...
```

## Control Points

The importer recognizes these major stations as control points:

- **BOS** - Boston South Station
- **BBY** - Back Bay
- **RTE** - Route 128
- **PVD** - Providence
- **NYP** - Penn Station NYC (Milepost 0.0)
- **NWK** - Newark Penn Station
- **TRE** - Trenton
- **PHL** - Philadelphia 30th Street
- **WIL** - Wilmington
- **BAL** - Baltimore Penn Station
- **BWI** - BWI Airport Station
- **WAS** - Washington Union Station

## Data Quality Notes

### Strengths
- ✅ Real geographic track layout
- ✅ Accurate station locations
- ✅ Max speed limits (where tagged)
- ✅ Track electrification status
- ✅ Platform data

### Limitations
- ⚠️ Signal positions incomplete (inferred from track segments)
- ⚠️ Block boundaries approximated (~2km segments)
- ⚠️ Control point assignments based on proximity
- ⚠️ Track numbers may not match railroad designations

### Improvements
This is a **starting point** for track modeling. Refine with:
- Official track charts (via FOIA)
- Actual signal positions
- Real control point codes
- Interlocking diagrams

## Integration with Optimization System

Once imported, the signal blocks are available for:

1. **Conflict Detection**: Check if two trains occupy the same block
2. **Routing**: Plan paths through the network
3. **Running Time Calculation**: Estimate travel times based on block lengths and speeds
4. **Capacity Analysis**: Identify bottlenecks

See `docs/track_modeling.md` for full details on using this data.

## Troubleshooting

**Overpass API timeout?**
- Try a smaller section (nyp-phil instead of full-nec)
- Wait a few minutes and retry (API has rate limits)

**No data returned?**
- Check bounding box coordinates
- Verify OpenStreetMap has data for your area
- Try with `-verbose` flag for details

**Database errors?**
- Ensure schema.sql has been run first
- Check table `signal_blocks` exists
- Verify column names match

## Next Steps

After importing:

1. **Visualize**: Query blocks and plot on map
2. **Validate**: Compare with known station distances
3. **Enhance**: Add interlocking data for major junctions
4. **Test**: Run conflict detection with sample trains

## Example Workflow

```bash
# 1. Build
cd cmd/osm_importer
go build

# 2. Import demo section
./osm_importer -section nyp-phil -verbose

# 3. Check the data
cat osm_data.json | jq '.elements | length'
cat import_blocks.sql | grep INSERT | wc -l

# 4. Load into database
psql -d amtrk_infra -f import_blocks.sql

# 5. Verify
psql -d amtrk_infra -c "SELECT COUNT(*), AVG(length), MIN(start_milepost), MAX(end_milepost) FROM signal_blocks;"

# 6. View blocks near Penn Station
psql -d amtrk_infra -c "SELECT block_id, control_point, start_milepost, end_milepost, length FROM signal_blocks WHERE control_point = 'NYP' ORDER BY start_milepost LIMIT 10;"
```

## References

- OpenStreetMap: https://www.openstreetmap.org/
- Overpass API: https://overpass-api.de/
- Railway tagging: https://wiki.openstreetmap.org/wiki/Railways
- Track modeling docs: `docs/track_modeling.md`
