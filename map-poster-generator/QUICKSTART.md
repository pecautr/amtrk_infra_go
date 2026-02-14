# Quick Start Guide - City Map Poster Generator

## 🚀 Get Started in 30 Seconds

### For Web Deployment

1. **Upload to your web server:**
   ```bash
   cp -r map-poster-generator /your/web/directory/
   ```

2. **Open in browser:**
   ```
   https://yoursite.com/map-poster-generator/index-full.html
   ```

3. **Done!** The application loads Leaflet from CDN automatically.

### For Local Testing

1. **Start server:**
   ```bash
   cd map-poster-generator
   python3 -m http.server 8080
   ```

2. **Open browser:**
   ```
   http://localhost:8080/index-full.html
   ```

3. **Create your map!**

## 📖 Basic Usage

1. **Search for a city** - Type city name and click Search
2. **Adjust the view** - Pan (drag) and zoom (scroll) to frame your area
3. **Choose a style** - Select from Black & White, Grayscale, Inverted, or Sepia
4. **Customize colors** - Pick custom street and background colors
5. **Set poster size** - Choose preset or enter custom dimensions
6. **Adjust border** - Use slider to set border width (0-50mm)
7. **Export** - Click "Export as PNG" or use Print Preview

## 🎨 Popular Presets

### Classic Black & White
- Style: Black & White (Classic)
- Street Color: #000000
- Background: #ffffff
- Border: 20mm

### Modern Minimalist
- Style: Grayscale
- Street Color: #2c2c2c
- Background: #f5f5f5
- Border: 10mm

### Dark Mode
- Style: Inverted
- Street Color: #ffffff
- Background: #1a1a1a
- Border: 15mm

### Vintage
- Style: Sepia
- Street Color: #5d4e37
- Background: #f4e4c1
- Border: 25mm

## 🖨️ Printing Tips

1. **Maximize window** before capturing/exporting
2. **Use landscape orientation** for wider maps
3. **Set zoom level 13-15** for neighborhood detail
4. **Zoom 11-12** for city overview
5. **Zoom 16-18** for street-level detail

## 💡 Pro Tips

- **Perfect framing**: Zoom out slightly more than needed, borders will crop it nicely
- **High quality**: Use Print Preview → Save as PDF for best quality
- **Custom sizes**: Works great for non-standard frames
- **Multiple cities**: Create a series with consistent styling
- **Experiment**: Try different combinations - there's no wrong answer!

## 🛠️ Troubleshooting

**Map not loading?**
- Check internet connection
- Verify you're using `index-full.html` not `index.html`
- Check browser console for errors

**Search not working?**
- Requires internet for geocoding
- Try different search terms (e.g., "Paris, France" instead of just "Paris")

**Export issues?**
- Use Print Preview → Save as PDF as alternative
- Browser screenshot tools work great too

**Colors not applying?**
- CSS filters work on the whole map layer
- Custom colors are experimental (best results with grayscale style)

## 📚 More Information

- **Full Documentation**: [README.md](README.md)
- **Deployment Guide**: [DEPLOYMENT.md](DEPLOYMENT.md)
- **Advanced Configuration**: See JavaScript files in `js/` directory

## 🎯 Example Workflow

**Create a Cincinnati poster in 2 minutes:**

1. Open `index-full.html`
2. Type "Cincinnati, OH" → Click Search
3. Zoom to desired level (try 13)
4. Select "Black & White (Classic)"
5. Choose poster size "18 × 24 in"
6. Adjust border to 20mm
7. Click "Print Preview"
8. Print or Save as PDF
9. Take to print shop or print at home!

## 🌟 Featured Cities to Try

- Cincinnati, OH (default)
- New York City, NY
- Chicago, IL
- San Francisco, CA
- Boston, MA
- Washington, DC
- Portland, OR
- Philadelphia, PA
- Seattle, WA
- Austin, TX

## ❓ Need Help?

- Check the [README](README.md) for detailed documentation
- See [DEPLOYMENT.md](DEPLOYMENT.md) for deployment options
- Open an issue in the repository for bugs or feature requests

---

**Happy Map Making! 🗺️**
