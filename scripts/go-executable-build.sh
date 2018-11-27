#!/usr/bin/env bash

set -eo pipefail

package=$1
if [[ -z "$package" ]]; then
  echo "usage: $0 <package-name> <[linux, windows, darwin, freebsd]/[amd64, 386, arm]>..."
  exit 1
fi
shift
package_split=(${package//\// })
package_name=${package_split[-1]}

platforms=("$@")

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name="$package_name-$GOOS-$GOARCH"
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi

    printf "Building $package $GOOS $GOARCH to $output_name\n"
    env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name $package
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
done

