run:
	@go build -ldflags "-X main.Version=dev"
	@ENV=development AWS_REGION=eu-west-1 ./shouts-api