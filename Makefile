dev:
	@templ generate
	@go run main.go

start:
	@templ generate
	@go build -o build/app .
	@./build/app
