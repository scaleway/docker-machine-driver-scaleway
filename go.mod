module github.com/scaleway/docker-machine-driver-scaleway

replace (
	github.com/Sirupsen/logrus => github.com/sirupsen/logrus v1.4.2
	github.com/docker/docker => github.com/docker/engine v1.4.2-0.20190822205725-ed20165a37b4
	github.com/docker/machine => /Users/qbrosse/go/src/github.com/docker/machine
	github.com/renstrom/fuzzysearch v1.1.0 => github.com/lithammer/fuzzysearch v1.1.0
// github.com/docker/machine => github.com/quentinbrosse/machine v0.7.1-0.20191224095347-2e18b6f9f2db
)

go 1.11

require (
	github.com/creack/goselect v0.1.1 // indirect
	github.com/docker/docker v1.13.2-0.20170601211448-f5ec1e2936dc // indirect
	github.com/docker/machine v0.16.2
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/google/go-cmp v0.3.1 // indirect
	github.com/gorilla/websocket v1.4.1 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/mattn/go-isatty v0.0.11 // indirect
	github.com/moul/gotty-client v1.7.0 // indirect
	github.com/moul/http2curl v1.0.0 // indirect
	github.com/pkg/errors v0.8.1 // indirect
	github.com/renstrom/fuzzysearch v1.1.0 // indirect
	github.com/scaleway/scaleway-cli v1.10.2-0.20190910170054-f0ff7c3359bb
	github.com/sirupsen/logrus v1.4.2
	github.com/smartystreets/goconvey v1.6.4 // indirect
	golang.org/x/crypto v0.0.0-20191219195013-becbf705a915 // indirect
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e // indirect
	golang.org/x/sys v0.0.0-20191224085550-c709ea063b76 // indirect
	gotest.tools v2.2.0+incompatible // indirect
	moul.io/anonuuid v1.2.1
)
