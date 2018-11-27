#!/usr/bin/env bash

set -eo pipefail

package=$1
version=$2
if [[ -z "$package" ]] || [[ -z "$version" ]]; then
  echo "usage: $0 <package-name> <version> <<bin_name>/<arch_name>/<archive_type>>..."
  exit 1
fi
shift
shift
package_split=(${package//\// })
package_name=${package_split[-1]}

package_path="$GOPATH/src/$package"
readme="$package_path/README.md"
license="$package_path/LICENSE"
basename="${package_name}_${version}_"
insidebinname="${package_name}"
archives=("$@")

for archive in "${archives[@]}"
do
    archive_split=(${archive//\// })
    bin=${archive_split[0]}
    arch=${archive_split[1]}
    type=${archive_split[2]}

    archive_name="${basename}${arch}.${type}"
    echo "Creating $archive_name"

    if [ "$type" = "deb" ]; then
        echo "Archive deb not implemented"
    elif [ "$type" = "zip" ]; then
        cp "${bin}" "/tmp/$insidebinname"
        rm -f "$archive_name"
        zip "$archive_name" -j "$license"  -j "$readme" -j "/tmp/$insidebinname"
        echo zip "$archive_name" -j "$license"  -j "$readme" -j "/tmp/$insidebinname"
    elif [ "$type" = "tar.gz" ]; then
        cp "${bin}" "/tmp/$insidebinname"
        tar -czvf "$archive_name" --directory=$(dirname "$license") $(basename "$license") --directory=$(dirname "$readme") $(basename "$readme") --directory=$(dirname "/tmp/$insidebinname") "$insidebinname"
    else
        echo "Archive $type not implemented"
    fi
done
