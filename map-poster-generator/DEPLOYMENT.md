# Map Poster Generator - Deployment Guide

This guide explains how to deploy and use the City Map Poster Generator.

## Quick Start

### Option 1: Web Deployment (Recommended)

Deploy to any web server with internet access:

1. **Upload Files:**
   ```bash
   # Upload the entire map-poster-generator directory to your web server
   scp -r map-poster-generator user@yourserver:/var/www/html/
   ```

2. **Access the Application:**
   - Navigate to: `https://yourserver.com/map-poster-generator/index-full.html`
   - The application will load Leaflet.js from CDN automatically

### Option 2: Local Development

For local testing with internet access:

1. **Start a Local Server:**
   ```bash
   cd map-poster-generator
   python3 -m http.server 8080
   ```

2. **Open in Browser:**
   ```
   http://localhost:8080/index-full.html
   ```

### Option 3: Offline Use

For environments without internet access:

1. **Download Leaflet.js:**
   - Visit: https://leafletjs.com/download.html
   - Download the latest stable version
   - Extract to `map-poster-generator/lib/leaflet/`

2. **Download html2canvas (optional for PNG export):**
   - Visit: https://github.com/niklasvh/html2canvas/releases
   - Download `html2canvas.min.js`
   - Place in `map-poster-generator/lib/`

3. **Update index-full.html:**
   Replace CDN links with local paths:
   ```html
   <!-- Change from: -->
   <link rel="stylesheet" href="https://unpkg.com/leaflet@1.9.4/dist/leaflet.css" />
   
   <!-- To: -->
   <link rel="stylesheet" href="lib/leaflet/leaflet.css" />
   ```

4. **Use the Application:**
   ```bash
   cd map-poster-generator
   python3 -m http.server 8080
   # Open http://localhost:8080/index-full.html
   ```

## File Structure

```
map-poster-generator/
├── index.html              # Demo page (no dependencies required)
├── index-full.html         # Full version (requires Leaflet.js)
├── README.md              # User documentation
├── DEPLOYMENT.md          # This file
├── css/
│   └── styles.css         # All styling
├── js/
│   ├── map.js            # Map initialization
│   ├── controls.js       # UI controls
│   └── export.js         # Export functionality
├── lib/                  # Optional: Local libraries
│   └── leaflet/          # Place Leaflet.js here for offline use
└── assets/               # Optional: Custom assets
```

## Configuration

### Changing Default Location

Edit `js/map.js`:

```javascript
// Change default center coordinates
const DEFAULT_CENTER = [39.1031, -84.5120]; // [latitude, longitude]
const DEFAULT_ZOOM = 13;
```

### Using Different Tile Providers

Edit `js/map.js` in the `updateTileLayer()` function:

```javascript
// Example: Use black and white tiles from Wikimedia
let tileUrl = 'https://tiles.wmflabs.org/bw-mapnik/{z}/{x}/{y}.png';

// Example: Use Stamen Toner for high contrast
let tileUrl = 'https://stamen-tiles.a.ssl.fastly.net/toner/{z}/{x}/{y}.png';
```

### Adding Preset Cities

Edit `js/controls.js`:

```javascript
const presetLocations = {
    'cincinnati': { lat: 39.1031, lon: -84.5120, zoom: 13 },
    'your-city': { lat: XX.XXXX, lon: -XX.XXXX, zoom: 12 },
    // Add more cities...
};
```

## Browser Requirements

- **Minimum Versions:**
  - Chrome 90+
  - Firefox 88+
  - Safari 14+
  - Edge 90+

- **Required Features:**
  - JavaScript enabled
  - HTML5 Canvas support
  - CSS Grid support

## Performance Optimization

### For Production Deployment:

1. **Minify CSS and JavaScript:**
   ```bash
   # Using terser and cssnano
   terser js/map.js js/controls.js js/export.js -o js/bundle.min.js
   cssnano css/styles.css css/styles.min.css
   ```

2. **Enable Compression:**
   Configure your web server to serve gzipped content:
   ```nginx
   # Nginx example
   gzip on;
   gzip_types text/css application/javascript;
   ```

3. **Use CDN with SRI:**
   The `index-full.html` already includes Subresource Integrity hashes for security.

## Troubleshooting

### Map Not Loading

**Problem:** Map tiles not appearing

**Solutions:**
- Check browser console for errors
- Verify internet connection (if using CDN)
- Check that Leaflet.js loaded correctly
- Try different tile provider URL

### Export Not Working

**Problem:** PNG/SVG export fails

**Solutions:**
- For PNG: Ensure html2canvas is loaded (or use browser screenshot tool)
- For SVG: This is a basic export; use QGIS for full vector maps
- Use Print Preview as alternative

### Search Not Finding Cities

**Problem:** City search returns no results

**Solutions:**
- Check internet connection (Nominatim requires online access)
- Try different search terms
- Use coordinates directly: `lat,lon` format

### CORS Errors

**Problem:** Cross-origin errors in console

**Solutions:**
- Serve from a web server (not `file://`)
- Use `python -m http.server` for local testing
- Configure CORS headers if self-hosting tiles

## Security Considerations

1. **Content Security Policy:**
   If using strict CSP, allow:
   ```
   script-src 'self' https://unpkg.com https://cdnjs.cloudflare.com;
   style-src 'self' https://unpkg.com;
   img-src 'self' https://*.tile.openstreetmap.org data:;
   connect-src https://nominatim.openstreetmap.org;
   ```

2. **Subresource Integrity:**
   Already included in `index-full.html` for CDN resources

3. **Rate Limiting:**
   - OpenStreetMap tiles: Respect usage policy
   - Nominatim search: Maximum 1 request per second
   - Consider caching for production use

## Production Deployment Checklist

- [ ] Test in all target browsers
- [ ] Minify CSS and JavaScript
- [ ] Enable GZIP compression
- [ ] Configure proper cache headers
- [ ] Test with slow network (3G simulation)
- [ ] Verify HTTPS is enabled
- [ ] Check mobile responsiveness
- [ ] Test print functionality
- [ ] Verify export features work
- [ ] Set up error tracking (optional)

## Advanced Usage

### Custom Map Layers

To add custom overlays (transit lines, etc.):

```javascript
// In js/map.js, after map initialization
const transitLayer = L.geoJSON(yourGeoJSONData, {
    style: { color: '#ff0000', weight: 3 }
});
transitLayer.addTo(map);
```

### Automated Batch Export

For bulk map generation, see the API documentation in README.md and use:

```javascript
// In browser console
loadPresetLocation('cincinnati');
setTimeout(() => printMap(), 2000);
```

## Support

For issues specific to this application, open an issue in the main repository.

For Leaflet.js questions, see: https://leafletjs.com/reference.html
For OpenStreetMap data, see: https://www.openstreetmap.org/

## License

This map poster generator is part of the Amtrak Infrastructure Optimizer project and is licensed under the MIT License.

Map data © OpenStreetMap contributors, licensed under ODbL.
