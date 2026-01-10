# ✅ Complete Setup Checklist

Use this checklist to ensure your development environment is properly configured for the StreamHub API Platform project.

## 🍎 macOS Environment Setup

### System Requirements
- [ ] macOS 10.15 (Catalina) or later
- [ ] At least 8GB RAM (16GB recommended)
- [ ] 20GB free disk space
- [ ] Stable internet connection

### Step 1: Essential Tools (15 minutes)

#### Homebrew
```bash
# Install Homebrew
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Verify
brew --version
```
- [ ] Homebrew installed and working

#### Go
```bash
# Install Go 1.21+
brew install go

# Verify
go version  # Should show 1.21 or higher

# Setup Go environment
echo 'export GOPATH=$HOME/go' >> ~/.zshrc
echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.zshrc
source ~/.zshrc
```
- [ ] Go installed (version 1.21+)
- [ ] GOPATH configured
- [ ] PATH includes $GOPATH/bin

#### Git
```bash
# Install Git
brew install git

# Configure Git
git config --global user.name "Your Name"
git config --global user.email "your.email@example.com"

# Verify
git --version
git config --list
```
- [ ] Git installed
- [ ] Git configured with name and email

#### Docker Desktop
```bash
# Install Docker Desktop
brew install --cask docker

# Start Docker Desktop
open -a Docker

# Wait for Docker to start (check menu bar icon)
# Verify
docker --version
docker-compose --version
docker ps
```
- [ ] Docker Desktop installed
- [ ] Docker Desktop running
- [ ] Can execute docker commands

#### VS Code
```bash
# Install VS Code
brew install --cask visual-studio-code

# Install 'code' command in PATH
# Open VS Code → Cmd+Shift+P → type "shell command" → "Install 'code' command in PATH"

# Verify
code --version
```
- [ ] VS Code installed
- [ ] 'code' command available in terminal

### Step 2: Development Tools (10 minutes)

```bash
# Install build tools
xcode-select --install

# Install additional CLI tools
brew install \
  make \
  postgresql@15 \
  redis \
  jq \
  httpie \
  tree

# Install Go development tools
go install github.com/99designs/gqlgen@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install golang.org/x/tools/gopls@latest
go install github.com/go-delve/delve/cmd/dlv@latest
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golang/mock/mockgen@latest

# Optional: Install k6 for load testing
brew install k6

# Optional: Install wscat for WebSocket testing
npm install -g wscat
```

- [ ] Xcode Command Line Tools installed
- [ ] PostgreSQL client installed
- [ ] Redis client installed
- [ ] Additional tools (jq, httpie, tree) installed
- [ ] gqlgen installed
- [ ] golangci-lint installed
- [ ] gopls (Go language server) installed
- [ ] delve (Go debugger) installed
- [ ] goimports installed
- [ ] mockgen installed
- [ ] k6 installed (optional)
- [ ] wscat installed (optional)

### Step 3: Project Setup (5 minutes)

```bash
# Create workspace directory
mkdir -p ~/workspace
cd ~/workspace

# Clone the project
git clone https://github.com/yourusername/streaming-platform-api.git
cd streaming-platform-api

# Download Go dependencies
go mod download
go mod verify

# Verify downloads
ls -la $GOPATH/pkg/mod
```

- [ ] Workspace directory created
- [ ] Project cloned
- [ ] Dependencies downloaded
- [ ] Dependencies verified

### Step 4: VS Code Configuration (5 minutes)

```bash
# Open project in VS Code
code .
```

When VS Code opens:
1. **Install Recommended Extensions**
   - [ ] Click "Install" when prompted for extensions
   - [ ] Wait for all extensions to install
   - [ ] Reload window if prompted

2. **Install Go Tools**
   - [ ] Press `Cmd+Shift+P`
   - [ ] Type "Go: Install/Update Tools"
   - [ ] Select all tools
   - [ ] Click OK
   - [ ] Wait for installation

3. **Verify Settings**
   - [ ] Check `.vscode/settings.json` exists
   - [ ] Check `.vscode/launch.json` exists
   - [ ] Check `.vscode/tasks.json` exists
   - [ ] Auto-format on save works (edit a .go file)

4. **Test Debugging**
   - [ ] Open `cmd/api-server/main.go`
   - [ ] Set a breakpoint (click left of line number)
   - [ ] Press `F5`
   - [ ] Select "Launch API Server"
   - [ ] Debugger should start (may fail if Docker not running)

### Step 5: Docker Services (5 minutes)

```bash
# Start all services
make docker-up

# Wait for services to be healthy (30-60 seconds)
# Watch logs
docker-compose -f deployments/docker/docker-compose.yml logs -f

# In another terminal, verify services
docker ps  # Should show 6-7 running containers

# Test connections
psql -h localhost -p 5432 -U streamhub -d streamhub
# Password: streamhub_password
# Type \q to exit

redis-cli -h localhost -p 6379
# Type PING, should respond PONG
# Type exit
```

- [ ] Docker Compose file exists
- [ ] Services started successfully
- [ ] PostgreSQL accessible (port 5432)
- [ ] Redis accessible (port 6379)
- [ ] RabbitMQ accessible (port 5672)
- [ ] Prometheus accessible (port 9090)
- [ ] Grafana accessible (port 3000)

### Step 6: Build and Test (5 minutes)

```bash
# Build the project
make build

# Verify binaries
ls -la bin/
# Should see: api-server, ws-server

# Run tests
make test

# Check test coverage
make test-coverage
open coverage.html
```

- [ ] Project builds without errors
- [ ] Binaries created in bin/
- [ ] All tests pass
- [ ] Coverage report generated
- [ ] Coverage > 80%

### Step 7: Run Application (5 minutes)

#### Method A: Using Make (2 terminals)

Terminal 1:
```bash
make run-api
# Should start on port 8080
```

Terminal 2:
```bash
make run-ws
# Should start on port 8081
```

- [ ] API server starts without errors
- [ ] WebSocket server starts without errors
- [ ] No port conflicts

#### Method B: Using VS Code Debugger

```bash
# In VS Code:
# Press F5
# Select "Launch All Servers"
```

- [ ] Both servers start in debug mode
- [ ] Can set breakpoints
- [ ] Can step through code

### Step 8: Verify Everything Works (5 minutes)

```bash
# Test GraphQL API
open http://localhost:8080/playground

# Try this query in playground:
```
```graphql
query {
  __schema {
    types {
      name
    }
  }
}
```

- [ ] GraphQL Playground opens
- [ ] Schema query works
- [ ] No errors in server logs

```bash
# Test WebSocket
wscat -c "ws://localhost:8081/ws?user_id=test_user"

# Send message:
{"type":"ping"}

# Should receive pong response
```

- [ ] WebSocket connection successful
- [ ] Can send/receive messages

```bash
# Test monitoring
open http://localhost:9090  # Prometheus
open http://localhost:3000  # Grafana (admin/admin)
open http://localhost:15672 # RabbitMQ (streamhub/streamhub_password)
```

- [ ] Prometheus accessible
- [ ] Grafana accessible
- [ ] RabbitMQ management accessible

### Step 9: GitHub Setup (10 minutes)

```bash
# Run setup script
./setup-github.sh

# Follow prompts:
# - Enter your GitHub username
# - Enter repository name (or use default)
# - Confirm file additions
```

Then on GitHub:
1. [ ] Go to https://github.com/new
2. [ ] Create repository: `streaming-platform-api`
3. [ ] Make it Public
4. [ ] Don't initialize with README
5. [ ] Click "Create repository"

Back in terminal:
```bash
git push -u origin main
```

- [ ] Repository created on GitHub
- [ ] Code pushed to GitHub
- [ ] Can view repository online
- [ ] README displays correctly

### Step 10: Optional Enhancements (10 minutes)

#### Add Repository Description
- [ ] Add description: "Production-ready API platform for real-time streaming"
- [ ] Add topics: `golang`, `graphql`, `websocket`, `api`, `real-time`, `docker`, `kubernetes`
- [ ] Add a website URL (if you have one)

#### Enable GitHub Actions
- [ ] Go to repository → Actions tab
- [ ] Enable workflows
- [ ] Verify CI/CD runs on next push

#### Setup Branch Protection
- [ ] Go to Settings → Branches
- [ ] Add rule for `main` branch
- [ ] Require status checks to pass
- [ ] Require pull request reviews

#### Add Secrets for CI/CD (Optional)
- [ ] Go to Settings → Secrets → Actions
- [ ] Add `DOCKER_USERNAME`
- [ ] Add `DOCKER_PASSWORD`

#### Pin Repository
- [ ] Go to your GitHub profile
- [ ] Click "Customize your pins"
- [ ] Select this repository
- [ ] Add to pinned repositories

---

## 🎉 Final Verification

Run this comprehensive check:

```bash
# Full system check
echo "=== System Check ==="
echo "Go version: $(go version)"
echo "Docker version: $(docker --version)"
echo "Git version: $(git --version)"
echo "VS Code: $(code --version | head -1)"
echo ""
echo "=== Go Tools ==="
which gqlgen
which golangci-lint
which gopls
echo ""
echo "=== Docker Services ==="
docker ps --format "table {{.Names}}\t{{.Status}}"
echo ""
echo "=== Project Status ==="
echo "Git branch: $(git branch --show-current)"
echo "Git remote: $(git remote get-url origin)"
echo "Go modules: $(go list -m all | wc -l) dependencies"
echo ""
echo "=== Verification Complete ==="
```

### Everything Working Checklist
- [ ] All system requirements met
- [ ] All tools installed and working
- [ ] Project cloned and dependencies downloaded
- [ ] VS Code configured with extensions
- [ ] Docker services running
- [ ] Project builds successfully
- [ ] All tests pass
- [ ] Application runs (both servers)
- [ ] GraphQL Playground accessible
- [ ] WebSocket connections work
- [ ] Monitoring dashboards accessible
- [ ] GitHub repository created and pushed
- [ ] CI/CD pipeline running

---

## 🆘 Troubleshooting

If any step fails, see:
- [SETUP_MACOS.md](SETUP_MACOS.md) - Detailed setup guide
- [QUICKSTART.md](QUICKSTART.md) - Fast setup alternative
- [docs/architecture.md](docs/architecture.md) - System overview

Common issues:
- **Port in use**: `kill -9 $(lsof -ti:8080)`
- **Docker not starting**: Restart Docker Desktop
- **Go tools missing**: Run `Go: Install/Update Tools` in VS Code
- **Module errors**: Run `go mod tidy`

---

## 📝 Next Steps

After completing this checklist:

1. ✅ **Review Documentation**
   - Read [docs/architecture.md](docs/architecture.md)
   - Review [docs/testing.md](docs/testing.md)
   - Check [docs/deployment.md](docs/deployment.md) (if exists)

2. 🎯 **Start Development**
   - Try modifying a handler
   - Add a new GraphQL query
   - Write a test

3. 🚀 **Deploy** (Future)
   - Set up AWS credentials
   - Review Terraform configs
   - Deploy to staging

---

**Estimated Total Time**: 60-90 minutes

**Questions?** Create an issue on GitHub!

**Ready to code?** Press `F5` in VS Code and start debugging! 🚀
