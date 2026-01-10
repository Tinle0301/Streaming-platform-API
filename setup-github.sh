#!/bin/bash

# GitHub Repository Setup Script for macOS
# This script helps initialize the repository on GitHub

set -e

echo "🚀 StreamHub API Platform - GitHub Setup"
echo "========================================"
echo ""

# Color codes
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Check if git is installed
if ! command -v git &> /dev/null; then
    echo -e "${RED}❌ Git is not installed. Please install git first.${NC}"
    exit 1
fi

# Check if we're in a git repository
if [ ! -d .git ]; then
    echo -e "${BLUE}📦 Initializing git repository...${NC}"
    git init
    echo -e "${GREEN}✓ Git repository initialized${NC}"
else
    echo -e "${GREEN}✓ Already in a git repository${NC}"
fi

# Ask for GitHub username
echo ""
echo -e "${YELLOW}Enter your GitHub username:${NC}"
read -r github_username

if [ -z "$github_username" ]; then
    echo -e "${RED}❌ GitHub username is required${NC}"
    exit 1
fi

# Ask for repository name
echo ""
echo -e "${YELLOW}Enter repository name (default: streaming-platform-api):${NC}"
read -r repo_name
repo_name=${repo_name:-streaming-platform-api}

# Create .gitignore if it doesn't exist
if [ ! -f .gitignore ]; then
    echo -e "${BLUE}📝 Creating .gitignore...${NC}"
    cat > .gitignore << 'EOF'
# See the full .gitignore in the project
bin/
*.out
.env
.DS_Store
vendor/
EOF
    echo -e "${GREEN}✓ .gitignore created${NC}"
fi

# Check for untracked files
echo ""
echo -e "${BLUE}📋 Checking git status...${NC}"
git status

# Add all files
echo ""
echo -e "${YELLOW}Do you want to add all files to git? (y/n)${NC}"
read -r add_files

if [ "$add_files" = "y" ] || [ "$add_files" = "Y" ]; then
    git add .
    echo -e "${GREEN}✓ Files staged${NC}"
fi

# Initial commit
if ! git rev-parse HEAD &> /dev/null; then
    echo ""
    echo -e "${BLUE}💾 Creating initial commit...${NC}"
    git commit -m "feat: initial commit - StreamHub API Platform

- Complete GraphQL API implementation
- WebSocket real-time messaging system
- Event-driven architecture with Redis/RabbitMQ
- Docker and Kubernetes deployment configs
- Comprehensive testing suite
- Complete documentation

This project demonstrates production-ready API platform engineering
skills for roles like Twitch API Platform Engineer."
    echo -e "${GREEN}✓ Initial commit created${NC}"
else
    echo -e "${GREEN}✓ Repository already has commits${NC}"
fi

# Set main branch
current_branch=$(git branch --show-current)
if [ "$current_branch" != "main" ]; then
    echo ""
    echo -e "${BLUE}🔀 Renaming branch to 'main'...${NC}"
    git branch -M main
    echo -e "${GREEN}✓ Branch renamed to main${NC}"
fi

# Set up remote
remote_url="https://github.com/$github_username/$repo_name.git"
echo ""
echo -e "${BLUE}🔗 Setting up remote repository...${NC}"
echo -e "   URL: ${remote_url}"

if git remote get-url origin &> /dev/null; then
    echo -e "${YELLOW}⚠️  Remote 'origin' already exists. Do you want to update it? (y/n)${NC}"
    read -r update_remote
    if [ "$update_remote" = "y" ] || [ "$update_remote" = "Y" ]; then
        git remote set-url origin "$remote_url"
        echo -e "${GREEN}✓ Remote updated${NC}"
    fi
else
    git remote add origin "$remote_url"
    echo -e "${GREEN}✓ Remote added${NC}"
fi

# Instructions for GitHub
echo ""
echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"
echo -e "${YELLOW}📝 Next Steps:${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"
echo ""
echo "1. Go to GitHub and create a new repository:"
echo -e "   ${BLUE}https://github.com/new${NC}"
echo ""
echo "   Repository name: ${GREEN}$repo_name${NC}"
echo "   Description: Production-ready API platform for real-time streaming"
echo "   Visibility: Public (recommended for portfolio)"
echo "   DON'T initialize with README, .gitignore, or license"
echo ""
echo "2. After creating the repository, run:"
echo -e "   ${GREEN}git push -u origin main${NC}"
echo ""
echo "3. (Optional) Set up GitHub Pages for documentation:"
echo "   - Go to Settings → Pages"
echo "   - Source: Deploy from a branch"
echo "   - Branch: main, /docs"
echo ""
echo "4. (Optional) Add repository topics:"
echo "   golang, graphql, websocket, api, real-time, docker, kubernetes"
echo ""
echo "5. (Optional) Add secrets for CI/CD:"
echo "   - Go to Settings → Secrets and variables → Actions"
echo "   - Add: DOCKER_USERNAME and DOCKER_PASSWORD"
echo ""
echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"
echo ""
echo -e "${GREEN}✅ Local repository setup complete!${NC}"
echo ""
echo -e "${YELLOW}Pro Tips:${NC}"
echo "• Add a GitHub profile README showcasing this project"
echo "• Pin this repository to your GitHub profile"
echo "• Add a detailed description and topics to the repository"
echo "• Enable GitHub Actions for automatic CI/CD"
echo "• Consider adding a project website using GitHub Pages"
echo ""
echo -e "${BLUE}Repository URL:${NC}"
echo -e "${GREEN}$remote_url${NC}"
echo ""
echo "Happy coding! 🚀"
