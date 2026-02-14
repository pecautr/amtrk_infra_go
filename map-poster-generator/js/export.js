// Export and Print Functionality

document.addEventListener('DOMContentLoaded', function() {
    initializeExportControls();
});

function initializeExportControls() {
    // PNG Export
    const exportPngBtn = document.getElementById('export-png');
    exportPngBtn.addEventListener('click', exportAsPNG);

    // SVG Export
    const exportSvgBtn = document.getElementById('export-svg');
    exportSvgBtn.addEventListener('click', exportAsSVG);

    // Print
    const printBtn = document.getElementById('print-btn');
    printBtn.addEventListener('click', openPrintPreview);
}

async function exportAsPNG() {
    const exportBtn = document.getElementById('export-png');
    
    // Check if html2canvas is available
    if (typeof html2canvas === 'undefined') {
        // Show info message
        exportBtn.textContent = 'Use Browser Screenshot';
        exportBtn.disabled = true;

        alert('PNG Export:\n\n' +
              '1. Use your browser\'s built-in screenshot tool:\n' +
              '   - Chrome/Edge: Press Ctrl+Shift+P (Cmd+Shift+P on Mac), type "screenshot"\n' +
              '   - Firefox: Right-click → "Take Screenshot"\n\n' +
              '2. Or use the Print Preview button to save as PDF, then convert to PNG\n\n' +
              '3. For high-quality exports, maximize your browser window before capturing');

        // Reset button after 2 seconds
        setTimeout(() => {
            exportBtn.textContent = 'Export as PNG';
            exportBtn.disabled = false;
        }, 2000);
        return;
    }

    const mapContainer = document.getElementById('map-preview');
    
    // Show loading state
    exportBtn.textContent = 'Generating PNG...';
    exportBtn.disabled = true;

    try {
        // Use html2canvas to capture the map
        const canvas = await html2canvas(mapContainer, {
            useCORS: true,
            allowTaint: true,
            backgroundColor: document.getElementById('background-color').value,
            scale: 2, // Higher quality
            logging: false
        });

        // Convert canvas to blob and download
        canvas.toBlob(function(blob) {
            const url = URL.createObjectURL(blob);
            const link = document.createElement('a');
            const timestamp = new Date().toISOString().slice(0, 10);
            const cityName = getCityNameFromSearch() || 'map';
            
            link.href = url;
            link.download = `${cityName}-poster-${timestamp}.png`;
            document.body.appendChild(link);
            link.click();
            document.body.removeChild(link);
            URL.revokeObjectURL(url);

            // Reset button
            exportBtn.textContent = 'Export as PNG';
            exportBtn.disabled = false;
        }, 'image/png');

    } catch (error) {
        console.error('Error exporting PNG:', error);
        alert('Failed to export PNG. Please try again.');
        
        // Reset button
        exportBtn.textContent = 'Export as PNG';
        exportBtn.disabled = false;
    }
}

function exportAsSVG() {
    // SVG export for vector graphics
    const exportBtn = document.getElementById('export-svg');
    
    // Show information
    exportBtn.textContent = 'Preparing SVG...';
    exportBtn.disabled = true;

    // Note: Full SVG export of raster tiles is complex
    // This creates a basic SVG wrapper with embedded image
    createSVGExport().then(() => {
        exportBtn.textContent = 'Export as SVG (Vector)';
        exportBtn.disabled = false;
    }).catch(error => {
        console.error('Error exporting SVG:', error);
        alert('SVG export requires additional processing. Use PNG export for now.');
        exportBtn.textContent = 'Export as SVG (Vector)';
        exportBtn.disabled = false;
    });
}

async function createSVGExport() {
    // Get map view details
    const mapView = getMapView();
    const mapContainer = document.getElementById('map-preview');
    const width = mapContainer.offsetWidth;
    const height = mapContainer.offsetHeight;
    
    // Create SVG element
    const svgNS = "http://www.w3.org/2000/svg";
    const svg = document.createElementNS(svgNS, "svg");
    svg.setAttribute("width", width);
    svg.setAttribute("height", height);
    svg.setAttribute("xmlns", svgNS);
    
    // Add background
    const rect = document.createElementNS(svgNS, "rect");
    rect.setAttribute("width", width);
    rect.setAttribute("height", height);
    rect.setAttribute("fill", document.getElementById('background-color').value);
    svg.appendChild(rect);

    // Add text overlay with map info
    const text = document.createElementNS(svgNS, "text");
    text.setAttribute("x", "10");
    text.setAttribute("y", "30");
    text.setAttribute("font-family", "Arial, sans-serif");
    text.setAttribute("font-size", "14");
    text.setAttribute("fill", document.getElementById('street-color').value);
    text.textContent = `Map: ${mapView.center.lat.toFixed(4)}, ${mapView.center.lng.toFixed(4)} | Zoom: ${mapView.zoom}`;
    svg.appendChild(text);

    // Note for user
    const note = document.createElementNS(svgNS, "text");
    note.setAttribute("x", "10");
    note.setAttribute("y", height - 10);
    note.setAttribute("font-family", "Arial, sans-serif");
    note.setAttribute("font-size", "12");
    note.setAttribute("fill", "#666");
    note.textContent = "For full vector export, consider using specialized GIS tools like QGIS with OpenStreetMap data";
    svg.appendChild(note);

    // Serialize SVG to string
    const serializer = new XMLSerializer();
    const svgString = serializer.serializeToString(svg);
    
    // Create blob and download
    const blob = new Blob([svgString], { type: 'image/svg+xml' });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    const timestamp = new Date().toISOString().slice(0, 10);
    const cityName = getCityNameFromSearch() || 'map';
    
    link.href = url;
    link.download = `${cityName}-poster-${timestamp}.svg`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    URL.revokeObjectURL(url);

    // Inform user about limitations
    alert('Basic SVG exported. Note: For true vector maps, use specialized tools like QGIS. The PNG export is recommended for poster printing.');
}

function openPrintPreview() {
    // Prepare the page for printing
    const originalTitle = document.title;
    
    // Update title for print
    const cityName = getCityNameFromSearch() || 'Custom Map';
    document.title = `${cityName} - Poster Map`;

    // Add print-specific styles
    const style = document.createElement('style');
    style.id = 'print-styles';
    style.textContent = `
        @media print {
            @page {
                margin: 0;
                size: auto;
            }
            
            body {
                margin: 0;
                padding: 0;
            }
            
            .control-panel,
            header,
            .map-info,
            .leaflet-control-zoom,
            .leaflet-control-attribution {
                display: none !important;
            }
            
            .main-content {
                display: block;
                margin: 0;
                padding: 0;
            }
            
            .map-wrapper {
                margin: 0;
                padding: 0;
                box-shadow: none;
            }
            
            .map-container {
                border-radius: 0;
                page-break-inside: avoid;
            }
        }
    `;
    document.head.appendChild(style);

    // Open print dialog
    setTimeout(() => {
        window.print();
        
        // Cleanup after print
        setTimeout(() => {
            document.title = originalTitle;
            const printStyle = document.getElementById('print-styles');
            if (printStyle) {
                printStyle.remove();
            }
        }, 500);
    }, 500);
}

function getCityNameFromSearch() {
    const searchInput = document.getElementById('city-search');
    const value = searchInput.value.trim();
    
    if (value) {
        // Clean up the city name for filename
        return value.replace(/[^a-zA-Z0-9-]/g, '-').toLowerCase();
    }
    
    return null;
}

// Helper function to download data as file
function downloadFile(data, filename, type) {
    const blob = new Blob([data], { type: type });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = filename;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    URL.revokeObjectURL(url);
}

// Export functions for external use
window.exportMapAsPNG = exportAsPNG;
window.exportMapAsSVG = exportAsSVG;
window.printMap = openPrintPreview;
