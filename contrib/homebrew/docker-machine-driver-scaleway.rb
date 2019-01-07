require "language/go"

class DockerMachineDriverScaleway < Formula
  desc "Docker Machine driver for Scaleway"
  homepage "https://github.com/scaleway/docker-machine-driver-scaleway/"
  url "https://github.com/scaleway/docker-machine-driver-scaleway/archive/v1.6.tar.gz"
  sha256 "9d27fe10c2169ffa2a7de79d65ce5a682f36ca75e1a98681b93b84b34ed8a22d"

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
