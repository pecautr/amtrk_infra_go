// Map initialization and core functionality

// Default map center (Cincinnati, OH)
const DEFAULT_CENTER = [39.1031, -84.5120];
const DEFAULT_ZOOM = 13;

// Initialize the map
let map;
let currentTileLayer;
let streetColor = '#000000';
let backgroundColor = '#ffffff';

// Initialize map on page load
document.addEventListener('DOMContentLoaded', function() {
    initializeMap();
    updateCoordinatesDisplay();
});

function initializeMap() {
    // Create the map instance
    map = L.map('map', {
        center: DEFAULT_CENTER,
        zoom: DEFAULT_ZOOM,
        zoomControl: true,
        attributionControl: true
    });

    // Add custom zoom control position
    map.zoomControl.setPosition('topright');

    // Initialize with black and white tile layer
    updateTileLayer('black-white');

    // Update coordinates and zoom display on map movement
    map.on('move', updateCoordinatesDisplay);
    map.on('zoom', updateCoordinatesDisplay);
    map.on('moveend', updateCoordinatesDisplay);
}

function updateTileLayer(style) {
    // Remove existing tile layer if present
    if (currentTileLayer) {
        map.removeLayer(currentTileLayer);
    }

    // Remove all style classes
    const mapContainer = document.getElementById('map');
    mapContainer.classList.remove('black-white', 'grayscale', 'inverted', 'sepia');

    // Base OSM tile layer for streets and railways
    let tileUrl = 'https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png';
    
    // For better print quality, use alternative tile servers
    // Uncomment to use different tile providers:
    // tileUrl = 'https://{s}.tile.openstreetmap.fr/osmfr/{z}/{x}/{y}.png';
    // tileUrl = 'https://tiles.wmflabs.org/bw-mapnik/{z}/{x}/{y}.png'; // Built-in B&W

    currentTileLayer = L.tileLayer(tileUrl, {
        attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
        maxZoom: 19,
        minZoom: 3
    });

    currentTileLayer.addTo(map);

    // Apply CSS filter for different styles
    const leafletContainer = document.querySelector('.leaflet-container');
    
    switch(style) {
        case 'black-white':
            leafletContainer.classList.add('black-white');
            break;
        case 'grayscale':
            leafletContainer.classList.add('grayscale');
            break;
        case 'inverted':
            leafletContainer.classList.add('inverted');
            break;
        case 'sepia':
            leafletContainer.classList.add('sepia');
            break;
    }
}

function updateCoordinatesDisplay() {
    if (!map) return;

    const center = map.getCenter();
    const zoom = map.getZoom();

    // Update coordinate display
    document.getElementById('current-coords').textContent = 
        `${center.lat.toFixed(4)}, ${center.lng.toFixed(4)}`;
    
    // Update zoom display
    document.getElementById('current-zoom').textContent = zoom;
}

// Search for a city using Nominatim API (OpenStreetMap)
async function searchCity(query) {
    if (!query || query.trim() === '') return;

    try {
        const response = await fetch(
            `https://nominatim.openstreetmap.org/search?format=json&q=${encodeURIComponent(query)}&limit=5`
        );
        const results = await response.json();
        
        displaySearchResults(results);
    } catch (error) {
        console.error('Error searching for city:', error);
        alert('Failed to search for location. Please try again.');
    }
}

function displaySearchResults(results) {
    const resultsContainer = document.getElementById('search-results');
    
    if (!results || results.length === 0) {
        resultsContainer.innerHTML = '<div class="search-result-item">No results found</div>';
        resultsContainer.classList.add('active');
        return;
    }

    resultsContainer.innerHTML = '';
    
    results.forEach(result => {
        const item = document.createElement('div');
        item.className = 'search-result-item';
        item.textContent = result.display_name;
        item.addEventListener('click', () => {
            flyToLocation(parseFloat(result.lat), parseFloat(result.lon));
            resultsContainer.classList.remove('active');
            resultsContainer.innerHTML = '';
        });
        resultsContainer.appendChild(item);
    });

    resultsContainer.classList.add('active');
}

function flyToLocation(lat, lon, zoom = 13) {
    map.flyTo([lat, lon], zoom, {
        duration: 1.5,
        easeLinearity: 0.25
    });
}

// Apply custom colors to map (for future enhancement with custom layers)
function applyCustomColors(streetCol, bgCol) {
    streetColor = streetCol;
    backgroundColor = bgCol;
    
    // This is a placeholder for future implementation with custom vector tiles
    // For now, the CSS filters handle the basic color transformations
    console.log(`Colors updated: Streets=${streetColor}, Background=${backgroundColor}`);
}

// Get current map bounds for export
function getMapBounds() {
    return map.getBounds();
}

// Get current map view for export
function getMapView() {
    return {
        center: map.getCenter(),
        zoom: map.getZoom(),
        bounds: map.getBounds()
    };
}

// Export the map object for use in other scripts
window.mapInstance = map;
