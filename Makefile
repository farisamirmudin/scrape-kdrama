dev:
	@templ generate
	@go run main.go

build:
	@templ generate
	@go build -o build/app .
	@./build/app
