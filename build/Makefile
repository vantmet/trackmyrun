build-linux:
	GOOS="linux" GOARCH="amd64" go build -C app/tmr -ldflags="-X 'main.Version=$(VERSION)'" .
	GOOS="linux" GOARCH="amd64" go build -C app/tmr-auth -ldflags="-X 'main.Version=$(VERSION)'" .
	GOOS="linux" GOARCH="amd64" go build -C app/tmr-strava -ldflags="-X 'main.Version=$(VERSION)'" .
