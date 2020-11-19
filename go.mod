module github.com/scaleway/docker-machine-driver-scaleway

replace (
	github.com/Sirupsen/logrus => github.com/sirupsen/logrus v1.4.2
	github.com/docker/docker => github.com/docker/engine v1.4.2-0.20190822205725-ed20165a37b4
	github.com/renstrom/fuzzysearch v1.1.0 => github.com/lithammer/fuzzysearch v1.1.0
)

go 1.15

require (
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/creack/goselect v0.1.1 // indirect
	github.com/docker/docker v1.13.2-0.20170601211448-f5ec1e2936dc // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/docker/machine v0.16.2
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/google/go-cmp v0.3.1 // indirect
	github.com/gorilla/websocket v1.4.1 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.3 // indirect
	github.com/mattn/go-isatty v0.0.11 // indirect
	github.com/moul/gotty-client v1.7.0 // indirect
	github.com/moul/http2curl v1.0.0 // indirect
	github.com/pkg/errors v0.8.1 // indirect
	github.com/renstrom/fuzzysearch v1.1.0 // indirect
	github.com/samalba/dockerclient v0.0.0-20160531175551-a30362618471 // indirect
	github.com/scaleway/scaleway-cli v1.11.1
	github.com/scaleway/scaleway-sdk-go v1.0.0-beta.7
	github.com/sirupsen/logrus v1.7.0
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/stretchr/objx v0.1.1 // indirect
	golang.org/x/crypto v0.0.0-20201117144127-c1f2f97bffc9 // indirect
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e // indirect
	golang.org/x/sys v0.0.0-20201119102817-f84b799fce68 // indirect
	gopkg.in/yaml.v2 v2.3.0
	gotest.tools v2.2.0+incompatible // indirect
	moul.io/anonuuid v1.2.1
)
