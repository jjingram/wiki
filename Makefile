PORT=8080
DATABASE_URL='postgres://localhost:5432/wiki?user=wiki&password=wiki&sslmode=disable'

run:
	PORT=$(PORT) DATABASE_URL=$(DATABASE_URL) go run wiki.go
