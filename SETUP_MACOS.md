# macOS Development Setup Guide

Complete guide for setting up the StreamHub API Platform on macOS.

## Prerequisites

### 1. Install Homebrew (if not already installed)
```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

### 2. Install Go
```bash
# Install Go 1.21 or later
brew install go

# Verify installation
go version

# Set up Go environment (add to ~/.zshrc or ~/.bash_profile)
echo 'export GOPATH=$HOME/go' >> ~/.zshrc
echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.zshrc
source ~/.zshrc
```

### 3. Install Git
```bash
brew install git

# Configure Git
git config --global user.name "Your Name"
git config --global user.email "your.email@example.com"
```

### 4. Install Docker Desktop
```bash
# Install Docker Desktop
brew install --cask docker

# Start Docker Desktop from Applications
# Or launch from terminal:
open -a Docker

# Verify Docker is running
docker --version
docker-compose --version
```

### 5. Install VS Code
```bash
# Install Visual Studio Code
brew install --cask visual-studio-code

# Install code command in PATH
# Open VS Code, press Cmd+Shift+P, type "shell command"
# Select "Shell Command: Install 'code' command in PATH"
```

### 6. Install Development Tools
```bash
# Install Make
xcode-select --install

# Install additional tools
brew install \
  postgresql@15 \
  redis \
  golangci-lint \
  k6 \
  jq \
  httpie \
  tree

# Install Go development tools
go install github.com/99designs/gqlgen@latest
go install github.com/golang/mock/mockgen@latest
go install golang.org/x/tools/gopls@latest
go install github.com/go-delve/delve/cmd/dlv@latest
go install golang.org/x/tools/cmd/goimports@latest
```

## Project Setup

### 1. Clone the Repository
```bash
# Create a workspace directory
mkdir -p ~/workspace
cd ~/workspace

# Clone the project
git clone https://github.com/yourusername/streaming-platform-api.git
cd streaming-platform-api
```

### 2. Open in VS Code
```bash
# Open project in VS Code
code .
```

When VS Code opens, it will prompt you to:
- Install recommended extensions (click "Install All")
- Trust the workspace (click "Yes, I trust the authors")

### 3. Install Go Dependencies
```bash
# Download all Go modules
go mod download
go mod verify

# This may take a few minutes
```

### 4. Start Development Services
```bash
# Start PostgreSQL, Redis, RabbitMQ, etc.
make docker-up

# Wait for services to be healthy (about 30 seconds)
# Check logs:
docker-compose -f deployments/docker/docker-compose.yml logs -f
```

### 5. Verify Setup
```bash
# Run tests to verify everything works
make test

# Build the project
make build

# You should see binaries in bin/
ls -la bin/
```

## VS Code Setup

### Recommended Extensions (should auto-prompt to install)
- **Go** (golang.go) - Essential for Go development
- **GraphQL** (graphql.vscode-graphql) - GraphQL syntax highlighting
- **Docker** (ms-azuretools.vscode-docker) - Docker support
- **GitLens** (eamodio.gitlens) - Enhanced Git features
- **Error Lens** (usernamehw.errorlens) - Inline error display

### VS Code Settings
The project includes `.vscode/settings.json` with optimal settings for Go development:
- Auto-format on save
- Auto-import organization
- Linting with golangci-lint
- Integrated debugging

## Running the Application

### Method 1: Using Make (Recommended)
```bash
# Terminal 1: Start API server
make run-api

# Terminal 2: Start WebSocket server
make run-ws
```

### Method 2: Using VS Code Debugger
1. Press `F5` or go to Run → Start Debugging
2. Select "Launch All Servers" from the dropdown
3. Set breakpoints by clicking left of line numbers
4. Use Debug Console for interactive debugging

### Method 3: Direct Go Run
```bash
# Run API server
go run ./cmd/api-server/main.go

# Run WebSocket server (in another terminal)
go run ./cmd/ws-server/main.go
```

## Testing the Application

### 1. Access GraphQL Playground
```bash
# Open in browser
open http://localhost:8080/playground
```

Try this query:
```graphql
query {
  streams(limit: 5) {
    edges {
      node {
        id
        title
        viewerCount
      }
    }
  }
}
```

### 2. Test WebSocket Connection
```bash
# Install wscat
npm install -g wscat

# Connect to WebSocket server
wscat -c "ws://localhost:8081/ws?user_id=test_user"

# Send a subscribe message
{"type":"subscribe","data":{"room":"stream_123"}}
```

### 3. View Monitoring Dashboards
```bash
# Prometheus metrics
open http://localhost:9090

# Grafana dashboards
open http://localhost:3000
# Default login: admin/admin

# RabbitMQ management
open http://localhost:15672
# Login: streamhub/streamhub_password
```

## Common Tasks

### Run Tests
```bash
# All tests
make test

# With coverage
make test-coverage
open coverage.html

# Integration tests
make test-integration

# Load tests
make test-load
```

### Code Quality
```bash
# Run linter
make lint

# Format code
make fmt

# Generate GraphQL code
make generate
```

### Database Operations
```bash
# Run migrations
make migrate

# Rollback migrations
make migrate-down

# Create new migration
make migrate-create NAME=add_users_table

# Connect to PostgreSQL
psql -h localhost -p 5432 -U streamhub -d streamhub
# Password: streamhub_password
```

### Docker Management
```bash
# View logs
make docker-logs

# Restart services
make docker-down
make docker-up

# Clean everything
make docker-down
docker system prune -a
```

## Debugging Tips

### 1. VS Code Debugging
- Set breakpoints by clicking left of line numbers
- Press F5 to start debugging
- Use Debug Console for REPL
- Step through code with F10 (step over) and F11 (step into)

### 2. Delve (Go debugger)
```bash
# Debug specific file
dlv debug ./cmd/api-server/main.go

# Attach to running process
dlv attach $(pgrep api-server)
```

### 3. Logging
```bash
# Follow application logs
tail -f logs/api.log

# Watch Docker logs
docker-compose -f deployments/docker/docker-compose.yml logs -f api
```

## Performance Profiling

### CPU Profiling
```bash
# Start profiling
go test -cpuprofile=cpu.prof -bench=.

# Analyze profile
go tool pprof cpu.prof
# Commands: top, list, web
```

### Memory Profiling
```bash
# Generate memory profile
go test -memprofile=mem.prof -bench=.

# Analyze profile
go tool pprof mem.prof
```

### Live Profiling
```bash
# API server exposes pprof endpoints
open http://localhost:8080/debug/pprof/

# Generate heap profile
curl http://localhost:8080/debug/pprof/heap > heap.prof
go tool pprof heap.prof
```

## GitHub Workflow

### 1. Create a New Branch
```bash
# Create feature branch
git checkout -b feature/add-user-authentication

# Make changes, then commit
git add .
git commit -m "feat: add JWT authentication"
```

### 2. Push to GitHub
```bash
# First time pushing this project
git remote add origin https://github.com/yourusername/streaming-platform-api.git
git branch -M main
git push -u origin main

# For feature branches
git push -u origin feature/add-user-authentication
```

### 3. Create Pull Request
1. Go to GitHub repository
2. Click "Pull requests" → "New pull request"
3. Select your branch
4. Fill in PR template
5. Wait for CI checks to pass
6. Request review

## Troubleshooting

### Docker Issues
```bash
# Reset Docker
docker system prune -a --volumes

# Check Docker memory
# Go to Docker Desktop → Settings → Resources
# Increase memory to at least 4GB

# Restart Docker Desktop
killall Docker && open -a Docker
```

### Port Already in Use
```bash
# Find process using port 8080
lsof -ti:8080

# Kill the process
kill -9 $(lsof -ti:8080)
```

### Go Module Issues
```bash
# Clean module cache
go clean -modcache

# Re-download modules
go mod download
```

### VS Code Go Extension Issues
```bash
# Reinstall Go tools
Go: Install/Update Tools (Cmd+Shift+P)
# Select "Install All"
```

## macOS-Specific Tips

### 1. Use Multiple Terminal Tabs
- Cmd+T: New tab
- Cmd+W: Close tab
- Cmd+Shift+[: Previous tab
- Cmd+Shift+]: Next tab

### 2. Spotlight Search
- Cmd+Space, type "Docker" to launch Docker Desktop
- Cmd+Space, type "code ." to open current directory in VS Code

### 3. iTerm2 (Alternative Terminal)
```bash
# Install iTerm2 for better terminal experience
brew install --cask iterm2
```

### 4. Use Fish or Zsh
```bash
# Install Fish shell (optional)
brew install fish
chsh -s /usr/local/bin/fish

# Install Oh My Zsh (for Zsh)
sh -c "$(curl -fsSL https://raw.github.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"
```

## Next Steps

1. ✅ Complete this setup guide
2. 📖 Read the [Architecture Documentation](docs/architecture.md)
3. 🧪 Review the [Testing Guide](docs/testing.md)
4. 🚀 Try building a new feature
5. 📝 Update documentation as you learn

## Quick Reference

```bash
# Daily workflow commands
make docker-up      # Start services
make run           # Run both servers
make test          # Run tests
make lint          # Check code quality
code .             # Open in VS Code
make docker-logs   # View logs
make docker-down   # Stop services
```

## Getting Help

- 📚 Documentation: `docs/` directory
- 🐛 Issues: Create a GitHub issue
- 💬 Questions: Use GitHub Discussions
- 📧 Email: your.email@example.com

---

**Happy Coding! 🚀**
