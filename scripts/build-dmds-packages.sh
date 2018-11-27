#!/usr/bin/env bash

set -eo pipefail

version=$1
if [[ -z "$version" ]]; then
  echo "usage: $0 <version>"
  exit 1
fi

package="github.com/scaleway/docker-machine-driver-scaleway"

archives=("docker-machine-driver-scaleway-linux-arm/linux_arm/tar.gz" "docker-machine-driver-scaleway-linux-amd64/linux_amd64/tar.gz" "docker-machine-driver-scaleway-linux-386/linux_386/tar.gz" "docker-machine-driver-scaleway-freebsd-arm/freebsd_arm/zip" "docker-machine-driver-scaleway-freebsd-amd64/freebsd_amd64/zip" "docker-machine-driver-scaleway-freebsd-386/freebsd_386/zip" "docker-machine-driver-scaleway-darwin-386/darwin_386/zip" "docker-machine-driver-scaleway-darwin-amd64/darwin_amd64/zip" "docker-machine-driver-scaleway-linux-amd64/amd64/deb" "docker-machine-driver-scaleway-linux-arm/armhf/deb" "docker-machine-driver-scaleway-linux-386/i386/deb")

mkdir -p "./release"
cd "./release"

for archive in "${archives[@]}"
do
  archive_split=(${archive//\// })
  bin=${archive_split[0]}
  bin_split=(${bin//-/ })
  ../go-executable-build.sh "$package" "${bin_split[-2]}/${bin_split[-1]}" \
    && ../packages-build.sh "$package" "$version" "$archive" \
    && rm -f "$bin"
done
