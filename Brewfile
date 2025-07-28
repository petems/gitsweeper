# Brewfile for gitsweeper development environment
# Run `brew bundle` to install all dependencies

# Core development tools
brew "go"
brew "git"
brew "make"

# Testing tools
brew "golangci-lint"

# Ruby for Cucumber/Aruba tests
brew "ruby"

# Docker for integration tests
cask "docker" unless File.exist?("/Applications/Docker.app")

# Optional: Useful development tools
brew "gh"               # GitHub CLI
brew "jq"               # JSON processing
brew "tree"             # Directory structure visualization
brew "htop"             # Process monitoring

# Optional: Git tools
brew "git-delta"        # Better git diff
brew "lazygit"          # Terminal UI for git commands