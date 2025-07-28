#!/bin/bash

set -e

echo "🧪 Testing goreleaser configuration..."

# Check if goreleaser is installed
if ! command -v goreleaser &> /dev/null; then
    echo "❌ goreleaser not found!"
    echo "Install goreleaser: https://goreleaser.com/install/"
    echo "Or run: go install github.com/goreleaser/goreleaser@latest"
    exit 1
fi

echo "📋 Checking goreleaser version..."
goreleaser --version

echo ""
echo "🔧 Validating configuration..."
goreleaser check

echo ""
echo "🏗️  Running goreleaser build test..."
goreleaser build --snapshot --clean --single-target

echo ""
echo "📦 Running goreleaser release test (dry run)..."
goreleaser release --snapshot --clean --skip publish

echo ""
echo "✅ goreleaser configuration test passed!"
echo ""
echo "🎯 Configuration includes:"
echo "  • Multi-platform builds (linux, windows, darwin with amd64/arm64)"
echo "  • Single optimized build with stripped symbols"
echo "  • Automatic checksum generation"
echo "  • Homebrew Casks integration"
echo "  • Automatic changelog generation"
echo ""
echo "🚀 To create a real release:"
echo "1. Tag your release: git tag v1.0.0"
echo "2. Push the tag: git push origin v1.0.0"
echo "3. The GitHub workflow will automatically run goreleaser"
echo ""
echo "🧪 To test locally:"
echo "  make release-snapshot  # Creates a snapshot release"
echo "  make build-all         # Builds for all platforms"
echo ""
echo "🍺 Homebrew integration:"
echo "  • Automatic updates to petems/homebrew-gitsweeper"
echo "  • Users can install with: brew tap petems/gitsweeper && brew install gitsweeper" 