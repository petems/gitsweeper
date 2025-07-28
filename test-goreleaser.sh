#!/bin/bash

set -e

echo "Testing goreleaser configuration..."

# Check if goreleaser is installed
if ! command -v goreleaser &> /dev/null; then
    echo "Install goreleaser: https://goreleaser.com/install/"
    exit 1
fi

echo "Running goreleaser build test..."
goreleaser build --snapshot --clean --single-target

echo "Running goreleaser release test (dry run)..."
goreleaser release --snapshot --clean --skip-validate

echo "✅ goreleaser configuration test passed!"
echo ""
echo "To create a real release:"
echo "1. Tag your release: git tag v1.0.0"
echo "2. Push the tag: git push origin v1.0.0"
echo "3. The GitHub workflow will automatically run goreleaser"
echo ""
echo "To test locally:"
echo "  make release-snapshot  # Creates a snapshot release"
echo "  make build-all         # Builds for all platforms" 