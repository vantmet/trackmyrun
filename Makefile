build-linux:
	cd app; $ENV:GOOS="linux";$ENV:GOARCH="amd64"; go build -o ../docker-images/app
	docker build -t vantmet/tmr .