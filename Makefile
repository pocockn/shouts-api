run:
	@go build -ldflags "-X main.Version=dev"
	@ENV=development ./shouts-api