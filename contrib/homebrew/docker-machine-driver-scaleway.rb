require "language/go"

class DockerMachineDriverScaleway < Formula
  desc "Docker Machine driver for Scaleway"
  homepage "https://github.com/scaleway/docker-machine-driver-scaleway/"
  url "https://github.com/scaleway/docker-machine-driver-scaleway/archive/v1.2.1.tar.gz"
  sha256 "dcfb35c4dd0a5ae7836d7275598a5de6b123b8f983ae3faaba42c51c28212ec6"

  head "https://github.com/scaleway/docker-machine-driver-scaleway.git"

  depends_on "go" => :build
  depends_on "docker-machine" => :recommended

  def install
    ENV["GOPATH"] = buildpath
    path = buildpath/"src/github.com/scaleway/docker-machine-driver-scaleway"
    path.install Dir["{*,.git,.gitignore}"]

    cd path do
      system "go", "build", "-o", "#{bin}/docker-machine-driver-scaleway", "./main.go"
    end
  end

  test do
    output = shell_output("#{Formula["docker-machine"].bin}/docker-machine create --driver scaleway -h")
    assert_match "scaleway-name", output
  end
end
