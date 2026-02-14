# City Map Poster Generator

A web-based application for creating aesthetic black and white street and rail maps of any city that can be printed as custom posters.

## Features

🗺️ **Interactive Map Customization**
- Search and navigate to any city in the world
- Smooth zoom and pan controls
- Real-time map updates

🎨 **Visual Customization**
- Multiple map styles (Black & White, Grayscale, Inverted, Sepia)
- Custom color selection for streets and background
- Adjustable border width for poster framing

📏 **Flexible Poster Sizing**
- Pre-defined sizes: A4, A3, Letter, Tabloid, 18×24", 24×36"
- Custom dimensions support (mm or inches)
- Automatic aspect ratio adjustment

🛤️ **Map Features**
- Toggle streets, railways, water bodies, parks, and buildings
- Optimized for printing aesthetic interior design maps
- Based on OpenStreetMap data

💾 **Export Options**
- Export as high-quality PNG (2x resolution)
- Export as SVG (vector format)
- Direct print preview and printing

## Getting Started

### Prerequisites

- A modern web browser (Chrome, Firefox, Safari, Edge)
- Internet connection (for map tiles and geocoding)

### Installation

1. Simply open `index.html` in your web browser
2. No build process or dependencies to install!

### Usage

1. **Search for Your City**
   - Enter a city name in the search box
   - Click "Search" or press Enter
   - Select from the search results

2. **Customize the Map**
   - Zoom in/out using mouse wheel or zoom controls
   - Drag to pan to your desired area
   - Choose a map style from the dropdown
   - Adjust colors using the color pickers

3. **Configure Poster Settings**
   - Select a poster size or enter custom dimensions
   - Adjust the border width using the slider
   - Toggle map features (streets, railways, etc.)

4. **Export Your Map**
   - Click "Export as PNG" for high-quality raster image
   - Click "Export as SVG" for vector format (basic)
   - Click "Print Preview" to print directly

## Technical Details

### Built With

- **Leaflet.js** - Interactive map library
- **OpenStreetMap** - Map data and tiles
- **Nominatim** - Geocoding service for city search
- **html2canvas** - PNG export functionality
- **Pure HTML/CSS/JavaScript** - No frameworks required

### File Structure

```
map-poster-generator/
├── index.html          # Main HTML file
├── css/
│   └── styles.css      # All styling
├── js/
│   ├── map.js          # Map initialization and core functions
│   ├── controls.js     # UI controls and event handlers
│   └── export.js       # Export and print functionality
└── assets/             # (Optional) Custom assets
```

### Browser Compatibility

- ✅ Chrome 90+
- ✅ Firefox 88+
- ✅ Safari 14+
- ✅ Edge 90+

## Use Cases

- **Interior Design** - Create custom city map artwork for homes and offices
- **Travel Memories** - Print maps of cities you've visited
- **Custom Gifts** - Personalized map posters for friends and family
- **Educational Materials** - Create city maps for teaching
- **Urban Planning** - Quick visualization of city layouts

## Customization Guide

### Adding Custom Map Styles

Edit the `updateTileLayer()` function in `js/map.js` to add new tile providers:

```javascript
// Example: Add a custom tile provider
let tileUrl = 'https://your-tile-server/{z}/{x}/{y}.png';
```

### Adjusting Export Quality

Modify the `scale` parameter in `js/export.js`:

```javascript
const canvas = await html2canvas(mapContainer, {
    scale: 3, // Higher = better quality, but larger file size
    // ...
});
```

### Custom Color Presets

Add preset color schemes by modifying the controls in `js/controls.js`:

```javascript
const colorPresets = {
    'classic': { street: '#000000', bg: '#ffffff' },
    'night': { street: '#ffffff', bg: '#1a1a1a' },
    // Add more presets...
};
```

## Limitations

- SVG export is basic and exports a wrapper rather than true vector tiles
- For professional vector exports, use specialized GIS tools like QGIS
- Map tile availability depends on OpenStreetMap and tile server status
- Some detailed features require custom vector tile layers

## Future Enhancements

- [ ] True vector tile support for full SVG export
- [ ] Custom layer overlays (transit lines, bike paths)
- [ ] Preset color themes
- [ ] Save/load custom map configurations
- [ ] Batch export multiple cities
- [ ] Text annotations and labels
- [ ] Additional tile providers

## Contributing

Contributions are welcome! This project is part of the Amtrak Infrastructure Optimizer repository.

## License

MIT License - See LICENSE file in the root directory

## Credits

- Map data © [OpenStreetMap](https://www.openstreetmap.org/copyright) contributors
- Built with [Leaflet.js](https://leafletjs.com/)
- Geocoding by [Nominatim](https://nominatim.org/)

## Support

For issues or questions, please open an issue in the main repository.

---

**Note:** This application sits adjacent to the Amtrak Infrastructure Optimizer and does not interfere with or depend on the transit optimization functionality.
