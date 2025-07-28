class Gitsweeper < Formula
  desc "A command-line tool for cleaning up merged Git branches"
  homepage "https://github.com/petems/gitsweeper"
  version "v0.1.0"
  license "MIT"

  # Detect platform and prefer binary installation for speed
  if OS.mac? && Hardware::CPU.intel?
    url "https://github.com/petems/gitsweeper/releases/download/v0.1.0/gitsweeper-v0.1.0-darwin-amd64.tar.gz"
    sha256 "PLACEHOLDER_BINARY_SHA256"

    def install
      bin.install "gitsweeper"
    end
  elsif OS.mac? && Hardware::CPU.arm?
    url "https://github.com/petems/gitsweeper/releases/download/v0.1.0/gitsweeper-v0.1.0-darwin-arm64.tar.gz"
    sha256 "PLACEHOLDER_BINARY_ARM64_SHA256"

    def install
      bin.install "gitsweeper"
    end
  else
    # Fallback to source installation for other platforms or when binaries aren't available
    url "https://github.com/petems/gitsweeper/archive/refs/tags/v0.1.0.tar.gz"
    sha256 "PLACEHOLDER_SOURCE_SHA256"

    depends_on "go" => :build

    def install
      system "go", "build", *std_go_args(ldflags: "-s -w"), "./..."
    end
  end

  test do
    # Test that the binary was installed and can show help
    assert_match "gitsweeper", shell_output("#{bin}/gitsweeper --help")
    
    # Test version output
    assert_match version.to_s, shell_output("#{bin}/gitsweeper --version 2>&1", 1)
  end
end