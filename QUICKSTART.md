# Quick Start Guide - GitHub + VS Code + macOS

Get up and running in **10 minutes**! ⚡

## Step 1: Install Prerequisites (5 minutes)

```bash
# Open Terminal (Cmd+Space, type "Terminal")

# Install Homebrew
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Install everything at once
brew install go git && brew install --cask docker visual-studio-code

# Start Docker Desktop
open -a Docker
```

## Step 2: Setup Project (2 minutes)

```bash
# Create workspace
mkdir -p ~/workspace && cd ~/workspace

# Clone the project
git clone <YOUR_GITHUB_URL>
cd streaming-platform-api

# Open in VS Code
code .
```

When VS Code opens:
1. Click **"Install"** when prompted for recommended extensions
2. Click **"Yes, I trust"** when prompted

## Step 3: Install Dependencies (2 minutes)

```bash
# In VS Code terminal (Ctrl+`):

# Install Go tools
go install github.com/99designs/gqlgen@latest
go install golang.org/x/tools/gopls@latest

# Download project dependencies
go mod download
```

## Step 4: Run the Application (1 minute)

### Option A: VS Code Debugger (Recommended)
1. Press `F5`
2. Select "Launch All Servers"
3. ✅ Done! Servers are running with debugging enabled

### Option B: Terminal
```bash
# Terminal 1
make docker-up    # Start databases

# Terminal 2  
make run-api      # Start API server

# Terminal 3
make run-ws       # Start WebSocket server
```

## Step 5: Verify It's Working

```bash
# Open GraphQL Playground
open http://localhost:8080/playground

# Test the API with this query:
```

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

## GitHub Setup

### Push to Your GitHub

```bash
# Run the setup script
./setup-github.sh

# Follow the prompts, then:
git push -u origin main
```

### Create Repository on GitHub
1. Go to https://github.com/new
2. Repository name: `streaming-platform-api`
3. Description: "Production-ready API platform for real-time streaming"
4. Make it **Public** (for portfolio visibility)
5. **Don't** initialize with README
6. Click "Create repository"

## Essential VS Code Shortcuts

| Action | Shortcut |
|--------|----------|
| Open terminal | `Ctrl + `` |
| Command palette | `Cmd + Shift + P` |
| Quick open file | `Cmd + P` |
| Start debugging | `F5` |
| Toggle sidebar | `Cmd + B` |
| Find in files | `Cmd + Shift + F` |
| Format document | `Shift + Alt + F` |

## Common Tasks

```bash
# Run tests
make test

# Check code quality
make lint

# View logs
make docker-logs

# Stop everything
make docker-down
```

## Troubleshooting

### Port already in use?
```bash
# Kill process on port 8080
kill -9 $(lsof -ti:8080)
```

### Docker not starting?
```bash
# Restart Docker
killall Docker && open -a Docker
```

### VS Code not recognizing Go?
1. Press `Cmd + Shift + P`
2. Type "Go: Install/Update Tools"
3. Select all and click OK

## Next Steps

1. ✅ Complete this quick start
2. 📖 Read [SETUP_MACOS.md](SETUP_MACOS.md) for detailed setup
3. 🏗️ Read [docs/architecture.md](docs/architecture.md) to understand the system
4. 🧪 Read [docs/testing.md](docs/testing.md) for testing guide
5. 🚀 Start building!

## Daily Development Workflow

```bash
# Morning: Start services
make docker-up
code .              # Open VS Code

# Development
F5                  # Start debugging
Cmd+S              # Auto-format on save
make test          # Run tests

# Evening: Stop services
make docker-down
```

## Help & Resources

- 📚 **Full Documentation**: See `docs/` folder
- 🐛 **Issues**: [Report bugs on GitHub](https://github.com/yourusername/streaming-platform-api/issues)
- 💬 **Questions**: Use GitHub Discussions
- 📝 **Contributing**: See `CONTRIBUTING.md`

## Project Structure Overview

```
streaming-platform-api/
├── cmd/              ← Entry points (main.go files)
├── internal/         ← Private application code
│   ├── graphql/      ← GraphQL schema & resolvers
│   ├── websocket/    ← WebSocket hub & clients
│   └── events/       ← Event publishing
├── docs/             ← Documentation
├── .vscode/          ← VS Code configuration
└── Makefile          ← Build commands
```

---

**You're all set! 🎉**

Press `F5` in VS Code and start exploring the code with debugging enabled!
