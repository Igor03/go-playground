swagger:
	swag init -o swagger
run:
	go run main.go
dev:
	swag init -o swagger && go run main.go