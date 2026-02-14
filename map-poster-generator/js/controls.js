// UI Controls and Event Handlers

document.addEventListener('DOMContentLoaded', function() {
    initializeControls();
});

function initializeControls() {
    // City Search
    const searchBtn = document.getElementById('search-btn');
    const citySearch = document.getElementById('city-search');
    
    searchBtn.addEventListener('click', () => {
        searchCity(citySearch.value);
    });

    citySearch.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') {
            searchCity(citySearch.value);
        }
    });

    // Close search results when clicking outside
    document.addEventListener('click', (e) => {
        const searchResults = document.getElementById('search-results');
        if (!e.target.closest('.control-group') || !e.target.closest('#city-search')) {
            searchResults.classList.remove('active');
        }
    });

    // Map Style Selection
    const mapStyle = document.getElementById('map-style');
    mapStyle.addEventListener('change', (e) => {
        updateTileLayer(e.target.value);
    });

    // Color Pickers
    const streetColor = document.getElementById('street-color');
    const backgroundColor = document.getElementById('background-color');
    
    streetColor.addEventListener('change', (e) => {
        applyCustomColors(e.target.value, backgroundColor.value);
    });

    backgroundColor.addEventListener('change', (e) => {
        applyCustomColors(streetColor.value, e.target.value);
        document.querySelector('.map-wrapper').style.backgroundColor = e.target.value;
    });

    // Poster Size Selection
    const posterSize = document.getElementById('poster-size');
    const customSizeControls = document.getElementById('custom-size-controls');
    
    posterSize.addEventListener('change', (e) => {
        if (e.target.value === 'custom') {
            customSizeControls.style.display = 'block';
        } else {
            customSizeControls.style.display = 'none';
        }
        updatePosterDimensions(e.target.value);
    });

    // Border Width Slider
    const borderWidth = document.getElementById('border-width');
    const borderWidthValue = document.getElementById('border-width-value');
    
    borderWidth.addEventListener('input', (e) => {
        borderWidthValue.textContent = e.target.value;
        updateBorderWidth(e.target.value);
    });

    // Feature Toggles
    const featureToggles = {
        'show-streets': true,
        'show-railways': true,
        'show-water': true,
        'show-parks': true,
        'show-buildings': false
    };

    Object.keys(featureToggles).forEach(id => {
        const checkbox = document.getElementById(id);
        checkbox.addEventListener('change', (e) => {
            featureToggles[id] = e.target.checked;
            updateMapFeatures(featureToggles);
        });
    });

    // Custom Size Inputs
    const customWidth = document.getElementById('custom-width');
    const customHeight = document.getElementById('custom-height');
    const customUnit = document.getElementById('custom-unit');

    [customWidth, customHeight, customUnit].forEach(input => {
        input.addEventListener('change', () => {
            updateCustomPosterSize(
                customWidth.value,
                customHeight.value,
                customUnit.value
            );
        });
    });
}

function updatePosterDimensions(size) {
    const dimensions = {
        'a4': { width: 210, height: 297, unit: 'mm' },
        'a3': { width: 297, height: 420, unit: 'mm' },
        'letter': { width: 8.5, height: 11, unit: 'in' },
        'tabloid': { width: 11, height: 17, unit: 'in' },
        '18x24': { width: 18, height: 24, unit: 'in' },
        '24x36': { width: 24, height: 36, unit: 'in' }
    };

    if (dimensions[size]) {
        const { width, height, unit } = dimensions[size];
        console.log(`Poster size set to: ${width} × ${height} ${unit}`);
        
        // Update map container aspect ratio
        updateMapAspectRatio(width, height);
    }
}

function updateCustomPosterSize(width, height, unit) {
    if (!width || !height) return;
    
    console.log(`Custom poster size: ${width} × ${height} ${unit}`);
    updateMapAspectRatio(width, height);
}

function updateMapAspectRatio(width, height) {
    const mapContainer = document.querySelector('.map-container');
    const aspectRatio = width / height;
    
    // Calculate new dimensions maintaining aspect ratio
    const containerWidth = mapContainer.offsetWidth;
    const newHeight = containerWidth / aspectRatio;
    
    // Don't make it too small or too large
    const minHeight = 400;
    const maxHeight = window.innerHeight - 280;
    const clampedHeight = Math.max(minHeight, Math.min(newHeight, maxHeight));
    
    mapContainer.style.height = `${clampedHeight}px`;
    
    // Force map to redraw with new dimensions
    setTimeout(() => {
        if (window.mapInstance) {
            window.mapInstance.invalidateSize();
        }
    }, 100);
}

function updateBorderWidth(width) {
    const mapPreview = document.getElementById('map-preview');
    mapPreview.style.padding = `${width}px`;
    mapPreview.style.backgroundColor = '#ffffff';
}

function updateMapFeatures(features) {
    // This is a placeholder for future implementation
    // In a full implementation, this would toggle different map layers
    console.log('Map features updated:', features);
    
    // For now, we can adjust opacity of certain elements via CSS
    const leafletContainer = document.querySelector('.leaflet-container');
    
    // Apply some basic feature visibility logic
    if (!features['show-streets']) {
        console.log('Streets hidden (not fully implemented in basic version)');
    }
    
    // Note: Full feature toggling would require custom vector tile layers
    // or additional overlay layers from specialized tile providers
}

// Preset locations for quick access
const presetLocations = {
    'cincinnati': { lat: 39.1031, lon: -84.5120, zoom: 13 },
    'new-york': { lat: 40.7128, lon: -74.0060, zoom: 12 },
    'chicago': { lat: 41.8781, lon: -87.6298, zoom: 12 },
    'san-francisco': { lat: 37.7749, lon: -122.4194, zoom: 13 },
    'boston': { lat: 42.3601, lon: -71.0589, zoom: 13 },
    'washington-dc': { lat: 38.9072, lon: -77.0369, zoom: 12 }
};

// Quick location function
function loadPresetLocation(locationKey) {
    const location = presetLocations[locationKey];
    if (location && window.mapInstance) {
        flyToLocation(location.lat, location.lon, location.zoom);
    }
}

// Export for use in console/debugging
window.loadPresetLocation = loadPresetLocation;
window.presetLocations = presetLocations;
