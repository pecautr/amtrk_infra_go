# Contributing to Amtrak NEC Optimization

Thank you for your interest in contributing! This project aims to improve train dispatching and reduce delays on Amtrak's Northeast Corridor.

## How to Contribute

### Reporting Issues

- Use GitHub Issues to report bugs or suggest features
- Include detailed information: steps to reproduce, expected vs actual behavior
- For performance issues, include timing information and system specs

### Code Contributions

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/your-feature-name`
3. **Make your changes**: Follow the coding standards below
4. **Write tests**: Ensure your code is well-tested
5. **Commit**: Use clear, descriptive commit messages
6. **Push**: `git push origin feature/your-feature-name`
7. **Submit a Pull Request**: Describe your changes and their motivation

### Coding Standards

- **Go Code**: Follow standard Go conventions
  - Run `go fmt ./...` before committing
  - Use meaningful variable and function names
  - Add comments for exported functions and types
  - Keep functions focused and testable

- **Testing**: 
  - Write unit tests for new functionality
  - Aim for >80% code coverage
  - Use table-driven tests where appropriate

- **Documentation**:
  - Update README.md for user-facing changes
  - Update relevant docs/ files for architectural changes
  - Add inline comments for complex logic

### Areas for Contribution

1. **Data Integration**
   - Contribute to [Catenary Transit Amtrak GTFS-RT](https://github.com/CatenaryTransit/amtrak-gtfs-rt)
   - Improve GTFS-RT parsing and validation
   - Add support for additional data sources

2. **Machine Learning**
   - Improve delay prediction models
   - Add new ML models (passenger demand, weather impact)
   - Optimize model training and inference

3. **Optimization**
   - Enhance MIP solver performance
   - Add new constraint types
   - Implement alternative optimization algorithms

4. **Vehicle Performance**
   - Add new trainset types and profiles
   - Improve physics calculations
   - Add energy consumption modeling

5. **Documentation**
   - Improve setup guides
   - Add tutorials and examples
   - Translate documentation

6. **Testing**
   - Increase test coverage
   - Add integration tests
   - Performance benchmarking

### Development Setup

See [docs/development.md](docs/development.md) for detailed setup instructions.

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/amtrk_infra_go.git
cd amtrk_infra_go

# Add upstream remote
git remote add upstream https://github.com/ryanpecora/amtrk_infra_go.git

# Install dependencies
go mod download

# Run tests
go test ./...

# Build
go build -o bin/orchestrator cmd/orchestrator/main.go
```

### Commit Message Guidelines

Use clear, descriptive commit messages:

- `feat: add conflict prediction for curved track sections`
- `fix: correct delay calculation for express trains`
- `docs: update ML model training guide`
- `test: add integration tests for optimization engine`
- `refactor: simplify vehicle performance calculator`

### Pull Request Process

1. Update documentation for any changed functionality
2. Ensure all tests pass: `go test ./...`
3. Update the README.md if needed
4. Your PR will be reviewed by maintainers
5. Address any feedback or requested changes
6. Once approved, your PR will be merged

### Code Review

All submissions require review. We aim to:

- Review PRs within 1 week
- Provide constructive feedback
- Help you improve your contribution

### Community Guidelines

- Be respectful and constructive
- Focus on the problem, not the person
- Welcome newcomers and help them learn
- Share knowledge and best practices

## Data Source Attribution

This project uses real-time data from:

- **[Catenary Transit Amtrak GTFS-RT](https://github.com/CatenaryTransit/amtrak-gtfs-rt)** - Please consider contributing to this project as well!

## Questions?

- Open a GitHub Issue for technical questions
- Check existing issues before creating new ones
- Join discussions in pull requests

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

Thank you for helping improve transit operations! 🚆
