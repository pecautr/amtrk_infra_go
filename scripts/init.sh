#!/bin/bash
# Initialization script for Amtrak NEC Optimization system

set -e

echo "Initializing Amtrak NEC Optimization System..."

# Create necessary directories
mkdir -p data/raw
mkdir -p data/processed
mkdir -p models
mkdir -p logs

# Create placeholder files
touch data/raw/.gitkeep
touch data/processed/.gitkeep
touch models/.gitkeep

echo "✓ Directories created"

# Initialize Go modules
echo "Downloading Go dependencies..."
go mod download
echo "✓ Go modules downloaded"

# Create default model metadata
cat > models/metadata.json << EOF
{
  "version": 0,
  "trained_at": "",
  "data_points": 0,
  "models": {
    "delay_predictor": "Not trained",
    "conflict_predictor": "Not trained",
    "dispatch_advisor": "Not trained",
    "performance_analyzer": "Not trained"
  }
}
EOF
echo "✓ Model metadata initialized"

# Create database schema (if PostgreSQL is available)
if command -v psql &> /dev/null; then
    echo "PostgreSQL detected. Would you like to initialize the database? (y/n)"
    read -r response
    if [[ "$response" == "y" ]]; then
        echo "Creating database schema..."
        psql -U amtrak -d amtrak_nec -f scripts/schema.sql
        echo "✓ Database initialized"
    fi
else
    echo "⚠ PostgreSQL not found. Skipping database initialization."
fi

echo ""
echo "================================================"
echo "Amtrak NEC Optimization System initialized!"
echo "================================================"
echo ""
echo "Next steps:"
echo "1. Set environment variables:"
echo "   export AMTRAK_API_KEY=your_api_key"
echo "   export DB_PASSWORD=your_db_password"
echo ""
echo "2. Build the application:"
echo "   go build -o bin/orchestrator cmd/orchestrator/main.go"
echo ""
echo "3. Run the orchestrator:"
echo "   ./bin/orchestrator"
echo ""
echo "For more information, see docs/development.md"
